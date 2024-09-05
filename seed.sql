-- Seed data for service_providers
INSERT INTO service_providers (name, specialty) VALUES
('Dr. Jane Smith', 'General Practitioner'),
('Dr. John Doe', 'Dentist'),
('Dr. Emily Brown', 'Psychologist');

-- Seed data for working_hours
-- Assuming IDs 1, 2, 3 for the service providers inserted above
INSERT INTO working_hours (service_provider_id, day_of_week, start_time, end_time) VALUES
(1, 1, '09:00', '17:00'), -- Monday
(1, 2, '09:00', '17:00'), -- Tuesday
(1, 3, '09:00', '17:00'), -- Wednesday
(1, 4, '09:00', '17:00'), -- Thursday
(1, 5, '09:00', '15:00'), -- Friday
(2, 1, '08:00', '16:00'), -- Monday
(2, 2, '08:00', '16:00'), -- Tuesday
(2, 3, '08:00', '16:00'), -- Wednesday
(2, 4, '08:00', '16:00'), -- Thursday
(2, 5, '08:00', '14:00'), -- Friday
(3, 1, '10:00', '18:00'), -- Monday
(3, 2, '10:00', '18:00'), -- Tuesday
(3, 3, '10:00', '18:00'), -- Wednesday
(3, 4, '10:00', '18:00'), -- Thursday
(3, 5, '10:00', '16:00'); -- Friday

-- Seed data for appointments
-- Let's add some appointments for the next 7 days
INSERT INTO appointments (service_provider_id, client_name, client_email, appointment_date, start_time, end_time, status)
SELECT
  (ARRAY[1, 2, 3])[floor(random() * 3 + 1)::int] as service_provider_id,
  'Client ' || generate_series as client_name,
  'client' || generate_series || '@example.com' as client_email,
  current_date + (generate_series % 7 || ' days')::interval as appointment_date,
  ('09:00'::time + (floor(random() * 16) * '30 minutes'::interval))::time as start_time,
  ('09:00'::time + (floor(random() * 16) * '30 minutes'::interval) + '30 minutes'::interval)::time as end_time,
  (ARRAY['booked', 'completed', 'canceled'])[floor(random() * 3 + 1)::int] as status
FROM generate_series(1, 50);
