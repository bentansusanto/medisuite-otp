-- name: CreateUser :one
INSERT INTO users(
  name,
  email,
  password,
  phone_number,
  role_id,
  is_verified,
  verify_code,
  verify_expires_at
)VALUES(
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING
  id,
  name,
  email,
  phone_number,
  role_id,
  is_verified,
  verify_code,
  verify_expires_at,
  created_at,
  updated_at
;

-- name: FindUserByEmail :one
SELECT
    u.id, u.email, u.password, u.name, u.phone_number,
    u.role_id, u.is_verified, COALESCE(u.verify_code, '') AS verify_code,COALESCE(u.verify_expires_at, NOW()) AS verify_expires_at,u.created_at, u.updated_at,
    r.id as role_id, r.name as role_name, r.code as role_code,
    r.level as role_level, r.description as role_description, r.can_self_register as role_can_self_register
FROM users u
JOIN roles r ON u.role_id = r.id
WHERE u.email = $1;

-- name: FindUserById :one
SELECT
	u.id, u.email, u.name, u.phone_number,
	u.role_id, u.is_verified,
	COALESCE(u.verify_code, '') AS verify_code,
	COALESCE(u.verify_expires_at, NOW()) AS verify_expires_at,
	u.created_at, u.updated_at,
	r.id as role_id, r.name as role_name, r.code as role_code,
	r.level as role_level, r.description as role_description, r.can_self_register as role_can_self_register
FROM users u
JOIN roles r ON u.role_id = r.id
WHERE u.id = $1;

-- name: FindUserByVerifyCode :one
SELECT
    u.id, u.email, u.password, u.name, u.phone_number,
    u.role_id, u.is_verified, u.verify_code, u.verify_expires_at,
    u.created_at, u.updated_at,
    r.id as role_id, r.name as role_name, r.code as role_code,
    r.level as role_level, r.description as role_description, r.can_self_register as role_can_self_register
FROM users u
JOIN roles r ON u.role_id = r.id
WHERE u.verify_code = $1;


-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE($2, name),
    phone_number = COALESCE($3, phone_number),
    email = COALESCE($4, email),
    password = COALESCE($5, password),
    is_verified = COALESCE($6, is_verified),
    verify_code = $7,
    verify_expires_at = $8,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    name,
    email,
    phone_number,
    role_id,
    is_verified,
    created_at,
    updated_at;

-- name: DeleteUser :one
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: FindAllUsers :many
SELECT
    u.id, u.email, u.name, u.phone_number,
    u.role_id, u.is_verified,
    COALESCE(u.verify_code, '') as verify_code,
    COALESCE(u.verify_expires_at, NOW()) as verify_expires_at,
    u.created_at, u.updated_at,
    r.id as role_id_ref,
    r.name as role_name,
    r.code as role_code,
    r.level as role_level,
    COALESCE(r.description, '') as role_description,
    r.can_self_register as role_can_self_register
FROM users u
JOIN roles r ON u.role_id = r.id
ORDER BY u.created_at DESC
LIMIT $1;
