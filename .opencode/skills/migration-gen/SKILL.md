---
name: migration-gen
description: Creates a new numbered SQL migration pair (up/down) in backend/migrations following the golang-migrate naming convention {NNN}_{snake_case}.up.sql/.down.sql, picking the next sequence number after the highest existing one. Use when the user says "crie a migration", "nova migration", "migration para [tabela]", "migration gen", or when go-module-gen indicates an operational module needs a migration.
---

# Migration Generator

Create `{NNN}_{snake_case_name}.up.sql` and `{NNN}_{snake_case_name}.down.sql` under `backend/migrations/`.

## Step 0 — Clarification (MANDATORY)

Before writing SQL, confirm with the user:

- Table name and columns (with types and constraints).
- Indexes, foreign keys, and unique constraints.
- Whether data backfill or defaults are required.
- The exact reversible operation for the `.down` migration.

Never assume the schema silently.

## Step 1 — Determine the next sequence number

List the existing migrations and compute the next number:

```bash
Get-ChildItem backend/migrations -Filter "*.up.sql" | ForEach-Object { $_.Name.Split('_')[0] } | Sort-Object
```

`NNN` = highest existing number + 1 (zero-padded to 3 digits).

## Step 2 — Name the migration

Use `snake_case` describing the change:

- `013_create_appointments.up.sql` / `.down.sql`
- `014_add_fhir_resource_id_to_employees.up.sql` / `.down.sql`

The same base name must be used for both `.up` and `.down`.

## Step 3 — Write the `.up` migration

Idempotent DDL for PostgreSQL (used by `golang-migrate`):

```sql
CREATE TABLE IF NOT EXISTS module_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ...
);
```

Operational modules (`auth`, `staff`, `telemetry`, `audit_logs`, `notifications`, `schedule`) write to local PostgreSQL and require migrations. Clinical modules (FHIR) never create migrations.

## Step 4 — Write the `.down` migration

Reversible DDL, ordered to drop child tables before parents:

```sql
DROP TABLE IF EXISTS module_table;
```

## Step 5 — Validate

- Confirm exactly one `.up.sql` and one `.down.sql` exist for the new number.
- Confirm the number is unique (no duplicates in `backend/migrations/`).
- If the migration runner is exercised locally, run it; otherwise report the migration as pending.

## Definition of Done

- [ ] `NNN` follows the highest existing number.
- [ ] Both `.up.sql` and `.down.sql` exist with the same base name.
- [ ] DDL is PostgreSQL-compatible and reversible.
- [ ] No comments were added to the generated SQL.
