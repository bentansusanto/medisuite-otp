-- name: GetPermissionByModuleAction :one
SELECT * FROM permissions
WHERE module = $1 AND action = $2 AND is_active = true
LIMIT 1;


