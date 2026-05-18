CREATE TABLE telemetry_rooms (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    passcode VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE telemetry_beds (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES telemetry_rooms(id) ON DELETE CASCADE,
    bed_number VARCHAR(50) NOT NULL,
    patient_name VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    gender VARCHAR(50) NOT NULL,
    bpm INT NOT NULL,
    spo2 INT NOT NULL,
    temperature NUMERIC(4, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    condition VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO telemetry_rooms (id, name, passcode, description) VALUES
('b1a4a4b4-d599-4daa-b7a6-2921de4e52d1', 'Sala Verde', '1234', 'Quarto VIP 101 - Estável'),
('b1a4a4b4-d599-4daa-b7a6-2921de4e52d2', 'Sala Vermelha', '4321', 'UTI Adulto - Crítico'),
('b1a4a4b4-d599-4daa-b7a6-2921de4e52d3', 'Sala Amarela', '9999', 'Semi-Intensivo 103 - Moderado');

INSERT INTO telemetry_beds (id, room_id, bed_number, patient_name, age, gender, bpm, spo2, temperature, status, condition) VALUES
('c1a4a4b4-d599-4daa-b7a6-2921de4e52b1', 'b1a4a4b4-d599-4daa-b7a6-2921de4e52d1', 'Leito 01', 'Ana Silva', 34, 'Feminino', 78, 98, 36.5, 'normal', 'Normal'),
('c1a4a4b4-d599-4daa-b7a6-2921de4e52b2', 'b1a4a4b4-d599-4daa-b7a6-2921de4e52d2', 'Leito 02', 'Bruno Costa', 45, 'Masculino', 52, 95, 37.1, 'warning', 'Bradicardia'),
('c1a4a4b4-d599-4daa-b7a6-2921de4e52b3', 'b1a4a4b4-d599-4daa-b7a6-2921de4e52d2', 'Leito 03', 'Carlos Oliveira', 62, 'Masculino', 115, 92, 38.4, 'danger', 'Taquicardia'),
('c1a4a4b4-d599-4daa-b7a6-2921de4e52b4', 'b1a4a4b4-d599-4daa-b7a6-2921de4e52d3', 'Leito 04', 'Danielle Souza', 28, 'Feminino', 82, 99, 36.7, 'normal', 'Normal');
