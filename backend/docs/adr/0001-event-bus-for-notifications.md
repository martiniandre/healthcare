# Event Bus for Notifications

Context: Notifications can be triggered by gRPC handler completions (e.g., Encounter created), async processes (e.g., Vertex AI exam analysis done), and real-time server-side events (e.g., Telemetry threshold breach). An interceptor-based approach would miss the latter two categories.

Decision: Use an in-process event bus (`internal/shared/eventbus/`). Each module publishes typed domain events; the notifications module subscribes and persists to PostgreSQL. WebSocket streaming delivers to the frontend. The `Register()` wiring in `main.go` connects publishers to subscribers explicitly.

Rejected: gRPC interceptor (cannot capture async/real-time events), pull-based polling (adds latency and load), external message broker (unnecessary operational complexity for in-process delivery).
