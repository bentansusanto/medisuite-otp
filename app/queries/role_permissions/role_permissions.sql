-- name: GetRolePermissions :many
SELECT
    p.module, p.action, p.name
FROM role_permissions rp
JOIN permissions p ON rp.permission_id = p.id
WHERE rp.role_id = $1::uuid AND p.is_active = true;
