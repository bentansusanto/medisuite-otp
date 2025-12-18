package errors

import "errors"

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleInvalid       = errors.New("role invalid")
	ErrRoleAlreadyExists = errors.New("role already exists")
	ErrRoleDeleted       = errors.New("role deleted")
	ErrRoleSelfRegister  = errors.New("role self register")
)

var RoleErrorMessage = []error{
	ErrRoleNotFound,
	ErrRoleInvalid,
	ErrRoleAlreadyExists,
	ErrRoleDeleted,
	ErrRoleSelfRegister,
}
