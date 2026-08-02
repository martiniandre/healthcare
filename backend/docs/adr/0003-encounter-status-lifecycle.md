# Encounter status lifecycle — born in-progress, minimal transitions

Encounter creation produced divergent statuses per transport (gRPC defaulted `in-progress`, HTTP defaulted `finished`) and no transition rules existed. We decided the Encounter is born `in-progress` and only `in-progress → finished` and `in-progress → cancelled` are valid transitions for now; the scheduling module (planned) will extend the lifecycle with `planned → arrived → in-progress`.

Considered: full FHIR R4 state machine (planned/arrived/triaged/in-progress/onleave/finished/cancelled) enforced immediately — rejected until scheduling exists; no transition validation — rejected because it allows status regressions; defaulting `finished` — rejected as it misrepresents a live consultation.

Consequences: directly-recorded consultations always start `in-progress`; the frontend gains explicit Finalizar/Cancelar actions; `planned` arrives with the Appointment module.
