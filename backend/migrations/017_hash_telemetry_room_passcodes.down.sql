UPDATE telemetry_rooms SET passcode = '1234' WHERE name = 'Sala Verde';
UPDATE telemetry_rooms SET passcode = '4321' WHERE name = 'Sala Vermelha';
UPDATE telemetry_rooms SET passcode = '9999' WHERE name = 'Sala Amarela';

ALTER TABLE telemetry_rooms ALTER COLUMN passcode TYPE VARCHAR(50);
