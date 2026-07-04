package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/healthcare/backend/internal/shared/eventbus"
)

type Simulator struct {
	repo     Repository
	eventBus eventbus.Bus
	interval time.Duration
	stopChan chan struct{}
}

func NewSimulator(repo Repository, eventBus eventbus.Bus) *Simulator {
	return &Simulator{
		repo:     repo,
		eventBus: eventBus,
		interval: 4 * time.Second,
		stopChan: make(chan struct{}),
	}
}

func (simulator *Simulator) Start(ctx context.Context) {
	go simulator.run(ctx)
}

func (simulator *Simulator) Stop() {
	close(simulator.stopChan)
}

func (simulator *Simulator) run(ctx context.Context) {
	ticker := time.NewTicker(simulator.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			simulator.tick(ctx)
		case <-simulator.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (simulator *Simulator) tick(ctx context.Context) {
	rooms, err := simulator.repo.GetRooms(ctx)
	if err != nil {
		slog.Warn("telemetry simulator: failed to get rooms", "error", err)
		return
	}

	for _, room := range rooms {
		beds, bedsErr := simulator.repo.GetBedsByRoomID(ctx, room.ID)
		if bedsErr != nil {
			slog.Warn("telemetry simulator: failed to get beds", "room_id", room.ID, "error", bedsErr)
			continue
		}

		for _, bed := range beds {
			previousStatus := bed.Status
			previousCondition := bed.Condition

			simulator.fluctuateVitals(bed)

			if updateErr := simulator.repo.UpdateBedCondition(ctx, bed); updateErr != nil {
				slog.Warn("telemetry simulator: failed to update bed", "bed_id", bed.ID, "error", updateErr)
				continue
			}

			if bed.Status == "danger" && previousStatus != "danger" && simulator.eventBus != nil {
				simulator.eventBus.Publish(ctx, eventbus.Event{
					Name: "telemetry.alert",
					Data: map[string]any{
						"title":         "Alerta Clínico - Leito " + bed.BedNumber,
						"body":          fmt.Sprintf("Paciente %s apresenta condição %s (BPM: %d, SpO2: %d%%).", bed.PatientName, bed.Condition, bed.Bpm, bed.Spo2),
						"resource_type": "bed",
						"resource_id":   bed.ID.String(),
					},
				})

				slog.Warn("telemetry alert published",
					"bed_id", bed.ID,
					"patient", bed.PatientName,
					"condition", bed.Condition,
					"previous_condition", previousCondition,
				)
			}
		}
	}
}

func (simulator *Simulator) fluctuateVitals(bed *Bed) {
	switch bed.Condition {
	case "Bradicardia":
		bed.Bpm = randomInt32(48, 58)
		bed.Spo2 = randomInt32(94, 97)
		bed.Temperature = randomFloat64(36.4, 37.2)
		bed.Status = "warning"

	case "Taquicardia":
		bed.Bpm = randomInt32(105, 130)
		bed.Spo2 = randomInt32(88, 94)
		bed.Temperature = randomFloat64(37.8, 39.1)
		bed.Status = "danger"

	case "Parada Cardíaca":
		bed.Bpm = 0
		bed.Spo2 = 0
		bed.Temperature = randomFloat64(34.0, 35.5)
		bed.Status = "danger"

	default:
		bed.Bpm = randomInt32(65, 95)
		bed.Spo2 = randomInt32(96, 100)
		bed.Temperature = randomFloat64(36.2, 37.3)
		bed.Status = "normal"
	}
}

func randomInt32(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

func randomFloat64(min, max float64) float64 {
	value := min + rand.Float64()*(max-min)
	return float64(int(value*10)) / 10
}
