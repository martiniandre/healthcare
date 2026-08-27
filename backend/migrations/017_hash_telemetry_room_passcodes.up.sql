ALTER TABLE telemetry_rooms ALTER COLUMN passcode TYPE VARCHAR(72);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

UPDATE telemetry_rooms SET passcode = crypt(passcode, gen_salt('bf'));
