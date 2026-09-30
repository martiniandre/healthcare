UPDATE telemetry_beds SET condition = CASE condition
    WHEN 'normal' THEN 'Normal'
    WHEN 'bradycardia' THEN 'Bradicardia'
    WHEN 'tachycardia' THEN 'Taquicardia'
    WHEN 'cardiac-arrest' THEN 'Parada Cardíaca'
    ELSE 'Normal'
END;

ALTER TABLE telemetry_beds ALTER COLUMN condition DROP DEFAULT;

DROP INDEX IF EXISTS idx_telemetry_beds_condition;
