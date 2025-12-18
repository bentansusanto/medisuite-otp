-- +goose Up

-- Role table
CREATE TABLE IF NOT EXISTS roles(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	code VARCHAR(50) NOT NULL UNIQUE,
	level INTEGER NOT NULL DEFAULT 0,
	description TEXT NULL,
	can_self_register BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Permissions table
CREATE TABLE permissions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	module VARCHAR(50) NOT NULL,
	action VARCHAR(50) NOT NULL,
	name VARCHAR(100) NOT NULL,
	description TEXT,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	UNIQUE(module, action),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Role permissions junction table
CREATE TABLE role_permissions (
	role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
	permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
	PRIMARY KEY (role_id, permission_id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User table
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  password VARCHAR NOT NULL,
  phone_number VARCHAR(50) NOT NULL,
  role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  is_verified BOOLEAN NOT NULL DEFAULT false,
  verify_code TEXT NULL,
  verify_expires_at TIMESTAMPTZ NULL,
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Partial unique index for email (only for non-deleted users)
-- This allows the same email to be reused after soft delete
CREATE UNIQUE INDEX idx_users_email_active
ON users(email)
WHERE deleted_at IS NULL;

-- User sessions for refresh tokens
CREATE TABLE user_sessions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	ref_token VARCHAR(512) NOT NULL,
	client_ip VARCHAR(45) NOT NULL,
	is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
	expires_at TIMESTAMPTZ NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO roles (name, code, level, description, can_self_register) VALUES
	('Owner', 'owner', 100, 'Owner system and clinic', true),
	('Doctor', 'doctor', 80, 'Doctor professional', false),
	('Admin', 'admin', 90, 'Admin professional', false),
	('Apoteker', 'apoteker', 70, 'Apoteker professional', false),
	('Cashier', 'cashier', 60, 'Cashier professional', false),
	('Patient', 'patient', 10, 'Patient', true);

INSERT INTO permissions (module, action, name, description) VALUES
    -- User management
    ('user', 'create', 'Create User', 'Create new users'),
    ('user', 'read', 'View Users', 'View user list'),
    ('user', 'update', 'Update User', 'Update user information'),
    ('user', 'delete', 'Delete User', 'Delete users'),
    -- Patient management
    ('patient', 'create', 'Create Patient', 'Create patient records'),
    ('patient', 'read', 'View Patients', 'View patient list'),
    ('patient', 'update', 'Update Patient', 'Update patient information'),
    ('patient', 'delete', 'Delete Patient', 'Delete patient records'),
        -- Appointment management
    ('appointment', 'create', 'Create Appointment', 'Create appointments'),
    ('appointment', 'read', 'View Appointments', 'View appointment list'),
    ('appointment', 'update', 'Update Appointment', 'Update appointments'),
    ('appointment', 'delete', 'Delete Appointment', 'Delete appointments');


-- Assign permissions to roles
-- Owner gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE code = 'owner'),
    id
FROM permissions;

-- Admin gets most permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE code = 'admin'),
    id
FROM permissions
WHERE module NOT IN ('user', 'role', 'permission')
   OR (module = 'user' AND action IN ('read', 'update'));

-- Apoteker permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE code = 'apoteker'),
    id
FROM permissions
WHERE module IN ('medicine', 'patient')
   AND action IN ('create', 'read', 'update');

-- Doctor permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE code = 'doctor'),
    id
FROM permissions
WHERE module IN ('patient', 'appointment', 'medicine')
   AND action IN ('read', 'update', 'create');

-- Patient permissions (minimal)
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE code = 'patient'),
    id
FROM permissions
WHERE module = 'appointment'
   AND action IN ('create', 'read', 'update', 'delete');


-- +goose Down
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
