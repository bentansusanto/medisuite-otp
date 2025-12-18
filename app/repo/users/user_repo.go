package users

import (
	"context"
	"log/slog"
	"time"

	errWrap "medisuite-api/common/errors"
	errConsts "medisuite-api/constants/errors"
	sessiondb "medisuite-api/pkg/db/user_sessions"
	userdb "medisuite-api/pkg/db/users"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type IUserRepo interface {
	Create(ctx context.Context, req userdb.CreateUserParams) (*userdb.CreateUserRow, error)
	FindUserByEmail(ctx context.Context, email string) (*userdb.FindUserByEmailRow, error)
	FindUserById(ctx context.Context, id uuid.UUID) (*userdb.FindUserByIdRow, error)
	FindAllUsers(ctx context.Context) ([]userdb.FindAllUsersRow, error)
	UpdateUser(ctx context.Context, req userdb.UpdateUserParams) (*userdb.UpdateUserRow, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (*userdb.User, error)
	FindUserByVerify(ctx context.Context, token string) (*userdb.FindUserByVerifyCodeRow, error)
	FindSessionByUserId(ctx context.Context, userID uuid.UUID) (*sessiondb.UserSession, error)
	SaveSession(ctx context.Context, req sessiondb.CreateSessionParams) (*sessiondb.UserSession, error)
	DeleteSession(ctx context.Context, token string) error
	FindSessions(ctx context.Context, token string) (*sessiondb.UserSession, error)
}

type UserRepo struct {
	uq *userdb.Queries
	sq *sessiondb.Queries
}

func NewUserRepo(uq *userdb.Queries, sq *sessiondb.Queries) IUserRepo {
	return &UserRepo{uq: uq, sq: sq}
}

// Repository method for creating a new user.
func (r *UserRepo) Create(ctx context.Context, req userdb.CreateUserParams) (*userdb.CreateUserRow, error) {
	// Set verification code and expiry (24 hours from now)

	user, err := r.uq.CreateUser(ctx, userdb.CreateUserParams{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		PhoneNumber:     req.PhoneNumber,
		RoleID:          req.RoleID,
		IsVerified:      false,
		VerifyCode:      req.VerifyCode,
		VerifyExpiresAt: req.VerifyExpiresAt,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// foreign key violation
			if pgErr.Code == "23503" { // foreign_key_violation
				slog.Error("FK violation on users.role_id", "role_id", req.RoleID, "error", pgErr)
				return nil, errWrap.WrapError(errConsts.ErrRoleNotFound)
			}
			// unique violation email
			if pgErr.Code == "23505" {
				slog.Error("Unique violation on users.email", "email", req.Email, "error", pgErr)
				return nil, errWrap.WrapError(errConsts.ErrUserEmailAlreadyExists)
			}
		}
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	return &user, nil
}

// Repository method for finding a user by email.
func (r *UserRepo) FindUserByEmail(ctx context.Context, email string) (*userdb.FindUserByEmailRow, error) {
	// find user by email in database
	user, err := r.uq.FindUserByEmail(ctx, email)
	if err != nil {
		// no rows means user does not exist
		if err == pgx.ErrNoRows {
			slog.Debug("User not found by email", "email", email)
			return nil, nil
		}
		slog.Error("Error finding user by email", "error", err, "email", email)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &user, nil
}

// Repository method for finding a user by ID.
func (r *UserRepo) FindUserById(ctx context.Context, id uuid.UUID) (*userdb.FindUserByIdRow, error) {
	// find user by id in database
	user, err := r.uq.FindUserById(ctx, id)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &user, nil
}

// Repository method for find all user
func (r *UserRepo) FindAllUsers(ctx context.Context) ([]userdb.FindAllUsersRow, error) {
	// Pass 1000 as limit (sensible default) and nil for isVerified to fetch all
	users, err := r.uq.FindAllUsers(ctx, 1000)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return users, nil
}

// Repository method for updating a user.
func (r *UserRepo) UpdateUser(ctx context.Context, req userdb.UpdateUserParams) (*userdb.UpdateUserRow, error) {
	// update user in database
	user, err := r.uq.UpdateUser(ctx, req)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &user, nil
}

// Repository method for deleting a user.
func (r *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) (*userdb.User, error) {
	// delete user in database
	user, err := r.uq.DeleteUser(ctx, id)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrUserDeleted)
	}
	return &user, nil
}

// Repository method for validating a user by verification code.
func (r *UserRepo) FindUserByVerify(ctx context.Context, verify_token string) (*userdb.FindUserByVerifyCodeRow, error) {
	user, err := r.uq.FindUserByVerifyCode(ctx, verify_token)
	if err != nil {
		// token tidak ditemukan di DB
		if err == pgx.ErrNoRows {
			return nil, errWrap.WrapError(errConsts.ErrUserNotFound)
		}
		// error DB lain
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	// check if user already has a verification code (sudah diverifikasi)
	if user.VerifyCode == "" {
		return nil, errWrap.WrapError(errConsts.ErrUserAlreadyVerified)
	}

	// check if the verification code has expired
	if time.Now().After(user.VerifyExpiresAt) {
		return nil, errWrap.WrapError(errConsts.ErrVerifyCodeExpired)
	}

	// verification code is still valid
	return &user, nil
}

// Repository method for retrieving sessions by user ID.
func (r *UserRepo) FindSessionByUserId(ctx context.Context, userID uuid.UUID) (*sessiondb.UserSession, error) {
	session, err := r.sq.GetSessionsByUserId(ctx, userID)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &session, nil
}

// Repository method for retrieving sessions by ref_token
func (r *UserRepo) FindSessions(ctx context.Context, token string) (*sessiondb.UserSession, error) {
	session, err := r.sq.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &session, nil
}

// Repository method for saving a new session.
func (r *UserRepo) SaveSession(ctx context.Context, req sessiondb.CreateSessionParams) (*sessiondb.UserSession, error) {
	// Create session in database
	session, err := r.sq.CreateSession(ctx, req)
	if err != nil {
		slog.Error("Failed to create user session", "error", err, "user_id", req.UserID, "ref_token", req.RefToken)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &session, nil
}

// Repository method for deleting a session.
func (r *UserRepo) DeleteSession(ctx context.Context, token string) error {
	err := r.sq.DeleteSession(ctx, token)
	if err != nil {
		return errWrap.WrapError(errConsts.ErrSQLError)
	}
	return nil
}
