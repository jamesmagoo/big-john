-- name: GetAppointment :one
SELECT * FROM appointments
WHERE id = $1 LIMIT 1;

-- name: ListAppointments :many
SELECT * FROM appointments
WHERE service_provider_id = $1 AND appointment_date = $2
ORDER BY start_time;

-- name: CreateAppointment :one
INSERT INTO appointments (
  service_provider_id, client_name, client_email, appointment_date, start_time, end_time, status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateAppointmentStatus :exec
UPDATE appointments
SET status = $2
WHERE id = $1;

-- name: DeleteAppointment :exec
DELETE FROM appointments
WHERE id = $1;

-- name: CheckAvailability :many
SELECT a.id
FROM appointments a
WHERE a.service_provider_id = $1
  AND a.appointment_date = $2
  AND (
    (a.start_time <= $3 AND a.end_time > $3)
    OR (a.start_time < $4 AND a.end_time >= $4)
    OR (a.start_time >= $3 AND a.end_time <= $4)
  );