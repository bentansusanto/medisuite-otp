-- name: CreateTreatment :one
INSERT INTO treatments
(category_id, name_treatment, description, thumbnail, price, duration, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: FindTreatments :many
SELECT id, category_id, name_treatment, description, thumbnail, price, duration, is_active, created_at, updated_at FROM treatments;

-- name: FindTreatmentById :one
SELECT id, category_id, name_treatment, description, thumbnail, price, duration, is_active, created_at, updated_at FROM treatments WHERE id = $1;

-- name: FindTreatmentByName :one
SELECT id, category_id, name_treatment, description, thumbnail, price, duration, is_active, created_at, updated_at FROM treatments WHERE name_treatment = $1;

-- name: UpdateTreatment :one
UPDATE treatments SET
  category_id = $2,
  name_treatment = $3,
  description = $4,
  thumbnail = $5,
  price = $6,
  duration = $7,
  is_active = $8
WHERE id = $1 RETURNING *;

-- name: DeleteTreatment :exec
DELETE FROM treatments WHERE id = $1;
