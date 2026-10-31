package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxStoredErrorLength = 500

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
}

type Record struct {
	ID        string
	Name      string
	Payload   []byte
	CreatedAt time.Time
	Attempts  int
}

type Failure struct {
	ID      string
	Cause   error
	RetryAt time.Time
	Dead    bool
}

type Stats struct {
	Pending       int64
	Dead          int64
	OldestPending time.Duration
}

type Store interface {
	Enqueue(ctx context.Context, executor Executor, event eventbus.Event) (string, error)
	ClaimBatch(ctx context.Context, batchSize int, lease time.Duration) ([]Record, error)
	MarkProcessed(ctx context.Context, eventID string) error
	MarkFailed(ctx context.Context, failure Failure) error
	PurgeProcessed(ctx context.Context, olderThan time.Duration) (int64, error)
	Stats(ctx context.Context) (Stats, error)
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (store *PostgresStore) Enqueue(ctx context.Context, executor Executor, event eventbus.Event) (string, error) {
	if executor == nil {
		executor = store.pool
	}
	payload, marshalError := json.Marshal(event.Data)
	if marshalError != nil {
		return "", fmt.Errorf("outbox: encode payload: %w", marshalError)
	}
	if event.Data == nil {
		payload = []byte("{}")
	}
	eventID := event.ID
	if eventID == "" {
		eventID = uuid.NewString()
	}
	_, execError := executor.Exec(ctx,
		`INSERT INTO outbox_events (id, name, payload) VALUES ($1, $2, $3)`,
		eventID, event.Name, payload,
	)
	if execError != nil {
		return "", fmt.Errorf("outbox: enqueue %s: %w", event.Name, execError)
	}
	return eventID, nil
}

func (store *PostgresStore) ClaimBatch(ctx context.Context, batchSize int, lease time.Duration) ([]Record, error) {
	rows, queryError := store.pool.Query(ctx, `
		UPDATE outbox_events
		SET attempts = attempts + 1,
		    available_at = NOW() + make_interval(secs => $2::double precision)
		WHERE id IN (
			SELECT id FROM outbox_events
			WHERE processed_at IS NULL AND dead_at IS NULL AND available_at <= NOW()
			ORDER BY available_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, name, payload, created_at, attempts`,
		batchSize, lease.Seconds(),
	)
	if queryError != nil {
		return nil, fmt.Errorf("outbox: claim batch: %w", queryError)
	}
	defer rows.Close()

	claimed := make([]Record, 0, batchSize)
	for rows.Next() {
		var record Record
		if scanError := rows.Scan(&record.ID, &record.Name, &record.Payload, &record.CreatedAt, &record.Attempts); scanError != nil {
			return nil, fmt.Errorf("outbox: scan claimed event: %w", scanError)
		}
		claimed = append(claimed, record)
	}
	if rowsError := rows.Err(); rowsError != nil {
		return nil, rowsError
	}
	sort.SliceStable(claimed, func(left, right int) bool {
		return claimed[left].CreatedAt.Before(claimed[right].CreatedAt)
	})
	return claimed, nil
}

func (store *PostgresStore) MarkProcessed(ctx context.Context, eventID string) error {
	_, execError := store.pool.Exec(ctx,
		`UPDATE outbox_events SET processed_at = NOW(), last_error = NULL WHERE id = $1`, eventID)
	return execError
}

func (store *PostgresStore) MarkFailed(ctx context.Context, failure Failure) error {
	lastError := truncateError(failure.Cause)
	if failure.Dead {
		_, execError := store.pool.Exec(ctx,
			`UPDATE outbox_events SET dead_at = NOW(), last_error = $2 WHERE id = $1`,
			failure.ID, lastError)
		return execError
	}
	_, execError := store.pool.Exec(ctx,
		`UPDATE outbox_events SET available_at = $3, last_error = $2 WHERE id = $1`,
		failure.ID, lastError, failure.RetryAt)
	return execError
}

func (store *PostgresStore) PurgeProcessed(ctx context.Context, olderThan time.Duration) (int64, error) {
	commandTag, execError := store.pool.Exec(ctx,
		`DELETE FROM outbox_events WHERE processed_at IS NOT NULL AND processed_at < NOW() - make_interval(secs => $1::double precision)`,
		olderThan.Seconds())
	if execError != nil {
		return 0, execError
	}
	return commandTag.RowsAffected(), nil
}

func (store *PostgresStore) Stats(ctx context.Context) (Stats, error) {
	var stats Stats
	var oldestSeconds float64
	queryError := store.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE processed_at IS NULL AND dead_at IS NULL),
			COUNT(*) FILTER (WHERE dead_at IS NOT NULL),
			COALESCE(EXTRACT(EPOCH FROM NOW() - MIN(created_at) FILTER (WHERE processed_at IS NULL AND dead_at IS NULL)), 0)
		FROM outbox_events`).Scan(&stats.Pending, &stats.Dead, &oldestSeconds)
	if queryError != nil {
		return Stats{}, queryError
	}
	stats.OldestPending = time.Duration(oldestSeconds * float64(time.Second))
	return stats, nil
}

func truncateError(cause error) string {
	if cause == nil {
		return ""
	}
	message := cause.Error()
	if len(message) > maxStoredErrorLength {
		return message[:maxStoredErrorLength]
	}
	return message
}
