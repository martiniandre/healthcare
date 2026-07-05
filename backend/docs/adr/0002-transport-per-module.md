# Transport per Module

**Status:** Accepted

**Context:** Each module needs a transport (gRPC, HTTP, or both). Browser-frontend modules cannot call gRPC without a proxy (gRPC-web, Envoy). Internal server-to-server communication benefits from gRPC's typed contracts and streaming.

**Decision:**
- **gRPC + HTTP (dual):** Modules that serve internal consumers AND the browser frontend. gRPC handler registered on the shared `grpc.Server` in `Register()`, HTTP handler passed to the router. Used by: allergy, audit_logs, auth, condition, diagnostic_report, encounter, exam_analyzer, health, imaging, medication, observation, patients, staff, telemetry.
- **HTTP-only:** Purely browser-facing modules with no internal consumers and no streaming data. Used by: analytics (read-only dashboard, no internal consumer), notifications (SSE real-time delivery is HTTP-native), portal (patient-facing, no internal consumer).

**Rejected options:**
- gRPC-only for browser modules — would require gRPC-web proxy for zero benefit
- Dual transport for analytics/portal — adds complexity with no internal consumers
- Removing HTTP from telemetry — browser dashboard needs HTTP reads
