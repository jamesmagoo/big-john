// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: appointment.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkAvailability = `-- name: CheckAvailability :many
SELECT a.id
FROM appointments a
WHERE a.service_provider_id = $1
  AND a.appointment_date = $2
  AND (
    (a.start_time <= $3 AND a.end_time > $3)
    OR (a.start_time < $4 AND a.end_time >= $4)
    OR (a.start_time >= $3 AND a.end_time <= $4)
  )
`

type CheckAvailabilityParams struct {
	ServiceProviderID int32       `json:"service_provider_id"`
	AppointmentDate   pgtype.Date `json:"appointment_date"`
	StartTime         pgtype.Time `json:"start_time"`
	StartTime_2       pgtype.Time `json:"start_time_2"`
}

func (q *Queries) CheckAvailability(ctx context.Context, arg CheckAvailabilityParams) ([]int32, error) {
	rows, err := q.db.Query(ctx, checkAvailability,
		arg.ServiceProviderID,
		arg.AppointmentDate,
		arg.StartTime,
		arg.StartTime_2,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createAppointment = `-- name: CreateAppointment :one
INSERT INTO appointments (
  service_provider_id, client_name, client_email, appointment_date, start_time, end_time, status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, service_provider_id, client_name, client_email, appointment_date, start_time, end_time, status, created_at
`

type CreateAppointmentParams struct {
	ServiceProviderID int32       `json:"service_provider_id"`
	ClientName        string      `json:"client_name"`
	ClientEmail       string      `json:"client_email"`
	AppointmentDate   pgtype.Date `json:"appointment_date"`
	StartTime         pgtype.Time `json:"start_time"`
	EndTime           pgtype.Time `json:"end_time"`
	Status            string      `json:"status"`
}

func (q *Queries) CreateAppointment(ctx context.Context, arg CreateAppointmentParams) (Appointment, error) {
	row := q.db.QueryRow(ctx, createAppointment,
		arg.ServiceProviderID,
		arg.ClientName,
		arg.ClientEmail,
		arg.AppointmentDate,
		arg.StartTime,
		arg.EndTime,
		arg.Status,
	)
	var i Appointment
	err := row.Scan(
		&i.ID,
		&i.ServiceProviderID,
		&i.ClientName,
		&i.ClientEmail,
		&i.AppointmentDate,
		&i.StartTime,
		&i.EndTime,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAppointment = `-- name: DeleteAppointment :exec
DELETE FROM appointments
WHERE id = $1
`

func (q *Queries) DeleteAppointment(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteAppointment, id)
	return err
}

const getAppointment = `-- name: GetAppointment :one
SELECT id, service_provider_id, client_name, client_email, appointment_date, start_time, end_time, status, created_at FROM appointments
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAppointment(ctx context.Context, id int32) (Appointment, error) {
	row := q.db.QueryRow(ctx, getAppointment, id)
	var i Appointment
	err := row.Scan(
		&i.ID,
		&i.ServiceProviderID,
		&i.ClientName,
		&i.ClientEmail,
		&i.AppointmentDate,
		&i.StartTime,
		&i.EndTime,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const listAppointments = `-- name: ListAppointments :many
SELECT id, service_provider_id, client_name, client_email, appointment_date, start_time, end_time, status, created_at FROM appointments
WHERE service_provider_id = $1 AND appointment_date = $2
ORDER BY start_time
`

type ListAppointmentsParams struct {
	ServiceProviderID int32       `json:"service_provider_id"`
	AppointmentDate   pgtype.Date `json:"appointment_date"`
}

func (q *Queries) ListAppointments(ctx context.Context, arg ListAppointmentsParams) ([]Appointment, error) {
	rows, err := q.db.Query(ctx, listAppointments, arg.ServiceProviderID, arg.AppointmentDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Appointment
	for rows.Next() {
		var i Appointment
		if err := rows.Scan(
			&i.ID,
			&i.ServiceProviderID,
			&i.ClientName,
			&i.ClientEmail,
			&i.AppointmentDate,
			&i.StartTime,
			&i.EndTime,
			&i.Status,
			&i.CreatedAt,
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

const updateAppointmentStatus = `-- name: UpdateAppointmentStatus :exec
UPDATE appointments
SET status = $2
WHERE id = $1
`

type UpdateAppointmentStatusParams struct {
	ID     int32  `json:"id"`
	Status string `json:"status"`
}

func (q *Queries) UpdateAppointmentStatus(ctx context.Context, arg UpdateAppointmentStatusParams) error {
	_, err := q.db.Exec(ctx, updateAppointmentStatus, arg.ID, arg.Status)
	return err
}
