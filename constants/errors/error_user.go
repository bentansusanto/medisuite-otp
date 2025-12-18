package errors

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrUserInvalid             = errors.New("user invalid")
	ErrUserPassword            = errors.New("user password")
	ErrUserPasswordNotMatch    = errors.New("user password not match")
	ErrUserEmailNotFound       = errors.New("user email not found")
	ErrUserEmailAlreadyExists  = errors.New("user email already exists")
	ErrUserEmailInvalid        = errors.New("user email invalid")
	ErrUserEmailPassword       = errors.New("user email password")
	ErrRoleAttempted           = errors.New("invalid role attempted")
	ErrOwnerAlreadyExists      = errors.New("owner already exists")
	ErrTokenInvalid            = errors.New("verification token is invalid")
	ErrTokenExpired            = errors.New("verification token is expired")
	ErrVerifyCodeExpired       = errors.New("verification code is expired")
	ErrUserAlreadyVerified     = errors.New("user is already verified")
	ErrUserNotVerified         = errors.New("user is not verified")
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrTokenNotFound           = errors.New("token not found")
	ErrVerifyCodeNotFound      = errors.New("verify code not found")
	ErrUserFailedToVerify      = errors.New("user failed to verify")
	ErrCreateUser              = errors.New("failed to create new user")
	ErrUpdatedUser             = errors.New("Failed to update user")
	ErrUserNotOwner            = errors.New("User not owner")
	ErrUserDeleted             = errors.New("User deleted")
	ErrInsufficientPermissions = errors.New("insufficient permissions for this action")
	ErrInvalidRole             = errors.New("invalid role")
)

var FindUserErr = []error{
	ErrUserNotFound,
	ErrUserInvalid,
	ErrUserPassword,
	ErrUserEmailNotFound,
	ErrUserEmailAlreadyExists,
	ErrUserEmailInvalid,
	ErrUserEmailPassword,
	ErrRoleAttempted,
	ErrOwnerAlreadyExists,
	ErrUserAlreadyVerified,
	ErrTokenInvalid,
	ErrTokenExpired,
	ErrUserNotVerified,
	ErrInvalidCredentials,
	ErrUserPasswordNotMatch,
	ErrTokenNotFound,
	ErrVerifyCodeNotFound,
	ErrUserFailedToVerify,
	ErrCreateUser,
	ErrUpdatedUser,
	ErrUserNotOwner,
	ErrUserDeleted,
	ErrInsufficientPermissions,
	ErrInvalidRole,
}
