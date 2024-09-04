-- name: ListServiceProviders :many
SELECT * FROM service_providers
ORDER BY name;

-- name: CreateServiceProvider :one
INSERT INTO service_providers (
  name, specialty
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateServiceProvider :exec
UPDATE service_providers
SET name = $2,
    specialty = $3
WHERE id = $1;

-- name: DeleteServiceProvider :exec
DELETE FROM service_providers
WHERE id = $1;