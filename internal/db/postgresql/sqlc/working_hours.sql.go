// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: working_hours.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createWorkingHours = `-- name: CreateWorkingHours :one
INSERT INTO working_hours (
  service_provider_id, day_of_week, start_time, end_time
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, service_provider_id, day_of_week, start_time, end_time
`

type CreateWorkingHoursParams struct {
	ServiceProviderID int32       `json:"service_provider_id"`
	DayOfWeek         int32       `json:"day_of_week"`
	StartTime         pgtype.Time `json:"start_time"`
	EndTime           pgtype.Time `json:"end_time"`
}

func (q *Queries) CreateWorkingHours(ctx context.Context, arg CreateWorkingHoursParams) (WorkingHour, error) {
	row := q.db.QueryRow(ctx, createWorkingHours,
		arg.ServiceProviderID,
		arg.DayOfWeek,
		arg.StartTime,
		arg.EndTime,
	)
	var i WorkingHour
	err := row.Scan(
		&i.ID,
		&i.ServiceProviderID,
		&i.DayOfWeek,
		&i.StartTime,
		&i.EndTime,
	)
	return i, err
}

const deleteWorkingHours = `-- name: DeleteWorkingHours :exec
DELETE FROM working_hours
WHERE service_provider_id = $1 AND day_of_week = $2
`

type DeleteWorkingHoursParams struct {
	ServiceProviderID int32 `json:"service_provider_id"`
	DayOfWeek         int32 `json:"day_of_week"`
}

func (q *Queries) DeleteWorkingHours(ctx context.Context, arg DeleteWorkingHoursParams) error {
	_, err := q.db.Exec(ctx, deleteWorkingHours, arg.ServiceProviderID, arg.DayOfWeek)
	return err
}

const getWorkingHours = `-- name: GetWorkingHours :many
SELECT id, service_provider_id, day_of_week, start_time, end_time FROM working_hours
WHERE service_provider_id = $1
`

func (q *Queries) GetWorkingHours(ctx context.Context, serviceProviderID int32) ([]WorkingHour, error) {
	rows, err := q.db.Query(ctx, getWorkingHours, serviceProviderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WorkingHour
	for rows.Next() {
		var i WorkingHour
		if err := rows.Scan(
			&i.ID,
			&i.ServiceProviderID,
			&i.DayOfWeek,
			&i.StartTime,
			&i.EndTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateWorkingHours = `-- name: UpdateWorkingHours :exec
UPDATE working_hours
SET start_time = $3,
    end_time = $4
WHERE service_provider_id = $1 AND day_of_week = $2
`

type UpdateWorkingHoursParams struct {
	ServiceProviderID int32       `json:"service_provider_id"`
	DayOfWeek         int32       `json:"day_of_week"`
	StartTime         pgtype.Time `json:"start_time"`
	EndTime           pgtype.Time `json:"end_time"`
}

func (q *Queries) UpdateWorkingHours(ctx context.Context, arg UpdateWorkingHoursParams) error {
	_, err := q.db.Exec(ctx, updateWorkingHours,
		arg.ServiceProviderID,
		arg.DayOfWeek,
		arg.StartTime,
		arg.EndTime,
	)
	return err
}
