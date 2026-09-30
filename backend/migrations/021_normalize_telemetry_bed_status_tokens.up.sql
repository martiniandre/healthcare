UPDATE telemetry_beds SET condition = CASE condition
    WHEN 'Normal' THEN 'normal'
    WHEN 'Bradicardia' THEN 'bradycardia'
    WHEN 'Taquicardia' THEN 'tachycardia'
    WHEN 'Parada Cardíaca' THEN 'cardiac-arrest'
    WHEN 'Parada Cardiaca' THEN 'cardiac-arrest'
    ELSE 'normal'
END;

ALTER TABLE telemetry_beds ALTER COLUMN condition SET DEFAULT 'normal';

CREATE INDEX idx_telemetry_beds_condition ON telemetry_beds (condition);
