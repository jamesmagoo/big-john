-- Create service_providers table
CREATE TABLE service_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    specialty VARCHAR(255)
);

-- Create working_hours table
CREATE TABLE working_hours (
    id SERIAL PRIMARY KEY,
    service_provider_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    CONSTRAINT fk_service_provider
        FOREIGN KEY(service_provider_id) 
        REFERENCES service_providers(id)
        ON DELETE CASCADE,
    CONSTRAINT valid_day_of_week 
        CHECK (day_of_week >= 0 AND day_of_week <= 6)
);

-- Create appointments table
CREATE TABLE appointments (
    id SERIAL PRIMARY KEY,
    service_provider_id INTEGER NOT NULL,
    client_name VARCHAR(255) NOT NULL,
    client_email VARCHAR(255) NOT NULL,
    appointment_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_service_provider
        FOREIGN KEY(service_provider_id) 
        REFERENCES service_providers(id)
        ON DELETE CASCADE
);

-- Create index on appointments for faster lookups
CREATE INDEX idx_appointments_service_provider_date ON appointments(service_provider_id, appointment_date);

-- Create index on working_hours for faster lookups
CREATE INDEX idx_working_hours_service_provider_day ON working_hours(service_provider_id, day_of_week);