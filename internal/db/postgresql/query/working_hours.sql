-- name: GetWorkingHours :many
SELECT * FROM working_hours
WHERE service_provider_id = $1;

-- name: CreateWorkingHours :one
INSERT INTO working_hours (
  service_provider_id, day_of_week, start_time, end_time
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateWorkingHours :exec
UPDATE working_hours
SET start_time = $3,
    end_time = $4
WHERE service_provider_id = $1 AND day_of_week = $2;

-- name: DeleteWorkingHours :exec
DELETE FROM working_hours
WHERE service_provider_id = $1 AND day_of_week = $2;