-- Drop indexes
DROP INDEX IF EXISTS idx_appointments_service_provider_date;
DROP INDEX IF EXISTS idx_working_hours_service_provider_day;

-- Drop appointments table
DROP TABLE IF EXISTS appointments;

-- Drop working_hours table
DROP TABLE IF EXISTS working_hours;

-- Drop service_providers table
DROP TABLE IF EXISTS service_providers;