-- name: GetRoleByID :one
SELECT id, name, code, level, description, can_self_register, created_at, updated_at FROM roles r WHERE r.id = $1;

-- name: GetRoleByCode :one
SELECT id, name, code, level, description, can_self_register, created_at, updated_at FROM roles WHERE code = $1 LIMIT 1;

-- name: GetAllRoles :many
SELECT id, name, code, level, description, can_self_register, created_at, updated_at FROM roles;
