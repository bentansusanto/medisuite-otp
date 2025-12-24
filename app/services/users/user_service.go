package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	userDTO "medisuite-api/app/dto/users"
	"medisuite-api/app/repo"
	"medisuite-api/common/emails"
	errWrap "medisuite-api/common/errors"
	"medisuite-api/config"
	errConsts "medisuite-api/constants/errors"
	"medisuite-api/constants/roles"
	"medisuite-api/constants/success"
	sessiondb "medisuite-api/pkg/db/user_sessions"
	userdb "medisuite-api/pkg/db/users"
	"medisuite-api/pkg/jwt"

	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(ctx context.Context, userDTO userDTO.RegisterDTO) (*userDTO.AuthResponse, error)
	VerifyAccount(ctx context.Context, token string) (*userDTO.AuthResponse, error)
	ResendVerifyAccount(ctx context.Context, email string) error
	LoginUser(ctx context.Context, req *userDTO.LoginDTO, clientIP string) (*userDTO.AuthResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	GetUser(ctx context.Context, userID uuid.UUID) (*userDTO.AuthResponse, error)
	RefreshToken(ctx context.Context, token string, clientIP string) (*userDTO.AuthResponse, error)
	ForgotPassword(ctx context.Context, req *userDTO.EmailRequest) error
	ResetPassword(ctx context.Context, req *userDTO.ResetPasswordDTO) error
}

type UserService struct {
	r repo.IRepo
}

func NewUserService(r repo.IRepo) IUserService {
	return &UserService{r: r}
}

// Service method for creating a new user.
func (s *UserService) CreateUser(ctx context.Context, req userDTO.RegisterDTO) (*userDTO.AuthResponse, error) {
	// find user by email
	findUser, err := s.r.UserRepo().FindUserByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("Error finding user by email", "error", err, "email", req.Email)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	// Check if user already exists
	if findUser != nil {
		slog.Error("User already exists with this email", "email", req.Email)
		return nil, errWrap.WrapError(errConsts.ErrUserEmailAlreadyExists)
	}

	// get role
	role, err := s.r.RoleRepo().FindRoleById(ctx, req.RoleID)
	if err != nil {
		slog.Error("Error finding role by id", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	// check if role exists
	if role == nil {
		slog.Error("Role not found", "role_id", req.RoleID)
		return nil, errWrap.WrapError(errConsts.ErrRoleNotFound)
	}

	// check if role can self register
	if !role.CanSelfRegister {
		slog.Error("Role cannot self register", "role_id", req.RoleID, "role_code", role.Code)
		return nil, errWrap.WrapError(errConsts.ErrRoleSelfRegister)
	}

	// check owner already exist (only if registering as OWNER)
	if role.Code == roles.OWNER {
		users, err := s.r.UserRepo().FindAllUsers(ctx)
		if err != nil {
			slog.Error("Error fetching users to check owner count", "error", err)
			return nil, errWrap.WrapError(errConsts.ErrSQLError)
		}

		// Check if any user already has OWNER role
		ownerExists := false
		for _, user := range users {
			// You need to check if user has OWNER role
			// Since FindAllUsers returns role_code, we can check
			userWithRole, err := s.r.UserRepo().FindUserById(ctx, user.ID)
			if err == nil && userWithRole != nil && userWithRole.RoleCode == roles.OWNER {
				ownerExists = true
				break
			}
		}

		if ownerExists {
			slog.Error("Owner already exists in the system")
			return nil, errWrap.WrapError(errConsts.ErrOwnerAlreadyExists)
		}
	}

	// hashpassword
	hashedPassword, err := config.HashPassword(req.Password)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrUserPassword)
	}

	verifyCode := config.GenerateRandomToken(32)
	verifyExpiresAt := time.Now().Add(24 * time.Hour)

	// payload
	payload := userdb.CreateUserParams{
		Name:            req.Name,
		Email:           req.Email,
		PhoneNumber:     req.PhoneNumber,
		Password:        hashedPassword,
		RoleID:          req.RoleID,
		VerifyCode:      verifyCode,
		VerifyExpiresAt: verifyExpiresAt,
	}

	newUser, err := s.r.UserRepo().Create(ctx, payload)
	if err != nil {
		slog.Error("Error creating user in repository", "error", err, "email", req.Email)
		return nil, err // Return the error as is (it's already wrapped)
	}

	// send email verification
	site := "http://localhost:3002"
	verificationLink := fmt.Sprintf(site+"/verify-account?verify_token=%s", verifyCode)
	emailBody := fmt.Sprintf("Thank you for registering with Bizpos. Please verify your account by clicking the link below:\n\n%s\n\nThis link will expire in 24 hours.", verificationLink)

	// send email (async or in background goroutine to not block response)
	go func() {
		errMail := emails.SendEmail([]string{newUser.Email}, nil,
			"Verify Your Account",
			emailBody)
		if errMail != nil {
			slog.Error("Error sending verification email", "error", errMail, "email", newUser.Email)
		} else {
			slog.Debug("Verification email sent successfully", "email", newUser.Email)
		}
	}()

	slog.Debug(success.SuccessCreateUser, "user_id", newUser.ID, "email", newUser.Email)

	// response body
	response := &userDTO.AuthResponse{
		ID:          newUser.ID,
		Name:        newUser.Name,
		Email:       newUser.Email,
		PhoneNumber: newUser.PhoneNumber,
		IsVerified:  newUser.IsVerified,
		RoleID:      newUser.RoleID,
		RoleDetail: &userDTO.RoleResponse{
			Name:            role.Name,
			Code:            role.Code,
			Level:           role.Level,
			Description:     role.Description,
			CanSelfRegister: role.CanSelfRegister,
		},
	}

	return response, nil
}

// helper verify account
func (s *UserService) validateVerificationToken(token string) error {
	// Check if token is empty
	if token == "" {
		slog.Error("verification token is empty")
		return errWrap.WrapError(errConsts.ErrTokenInvalid)
	}

	// Check token length (assuming tokens are 32 chars)
	token = strings.TrimSpace(token)
	if len(token) < 16 || len(token) > 128 {
		slog.Error("verification token is invalid")
		return errWrap.WrapError(errConsts.ErrTokenInvalid)
	}

	// Check for invalid characters (basic validation)
	for _, char := range token {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-' || char == '_') {
			slog.Error("verification token is invalid")
			return errWrap.WrapError(errConsts.ErrTokenInvalid)
		}
	}

	return nil
}

// Service method for verify account
func (s *UserService) VerifyAccount(ctx context.Context, verify_token string) (*userDTO.AuthResponse, error) {
	// Validate token
	if err := s.validateVerificationToken(verify_token); err != nil {
		return nil, err
	}

	// Find user by token
	user, err := s.r.UserRepo().FindUserByVerify(ctx, verify_token)
	if err != nil {
		switch {
		case errors.Is(err, errConsts.ErrUserNotFound):
			return nil, errWrap.WrapError(errConsts.ErrTokenNotFound)
		case errors.Is(err, errConsts.ErrUserAlreadyVerified):
			return nil, errWrap.WrapError(errConsts.ErrUserAlreadyVerified)
		case errors.Is(err, errConsts.ErrVerifyCodeExpired):
			return nil, errWrap.WrapError(errConsts.ErrVerifyCodeExpired)
		default:
			return nil, errWrap.WrapError(errConsts.ErrSQLError)
		}
	}
	if user.VerifyCode == "" {
		slog.Error("Token not found")
		return nil, errWrap.WrapError(errConsts.ErrTokenNotFound)
	}

	if user.IsVerified == true {
		slog.Error("user already verified")
		return nil, errWrap.WrapError(errConsts.ErrUserAlreadyVerified)
	}

	if time.Now().After(user.VerifyExpiresAt) {
		slog.Error("verify token expired")
		return nil, errWrap.WrapError(errConsts.ErrVerifyCodeExpired)
	}

	// Update user verification status
	updateUser, err := s.r.UserRepo().UpdateUser(ctx, userdb.UpdateUserParams{
		ID:              user.ID,
		Name:            user.Name,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Password:        user.Password,
		IsVerified:      true,
		VerifyCode:      nil,
		VerifyExpiresAt: nil,
	})
	if err != nil {
		slog.Error("User update failed", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrUpdatedUser)
	}

	slog.Debug(success.SuccessCreateUser)

	// response body
	response := &userDTO.AuthResponse{
		ID:          updateUser.ID,
		Name:        updateUser.Name,
		Email:       updateUser.Email,
		PhoneNumber: updateUser.PhoneNumber,
		IsVerified:  updateUser.IsVerified,
		RoleID:      updateUser.RoleID,
		UpdatedAt:   updateUser.UpdatedAt,
	}

	return response, nil
}

// Service method for resend verify account
func (s *UserService) ResendVerifyAccount(ctx context.Context, email string) error {
	// Validation email
	if email == "" {
		slog.Error("Email is empty")
		return errWrap.WrapError(errConsts.ErrUserEmailInvalid)
	}

	// find user by email
	findUser, err := s.r.UserRepo().FindUserByEmail(ctx, email)
	if err != nil {
		slog.Error("Error finding user by email", "error", err, "email", email)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}

	if findUser == nil {
		slog.Error("user email not found", "email", email)
		return errWrap.WrapError(errConsts.ErrUserEmailNotFound)
	}

	// check user already verified
	if findUser.IsVerified == true {
		slog.Error("user already verified", "email", email)
		return errWrap.WrapError(errConsts.ErrUserAlreadyVerified)
	}

	verifyToken := config.GenerateRandomToken(32)
	expiredAt := time.Now().Add(time.Hour * 24)

	// payload
	payload := userdb.UpdateUserParams{
		ID:              findUser.ID,
		Name:            findUser.Name,
		PhoneNumber:     findUser.PhoneNumber,
		Email:           findUser.Email,
		Password:        findUser.Password,
		IsVerified:      false,
		VerifyCode:      &verifyToken,
		VerifyExpiresAt: &expiredAt,
	}

	// execute update user
	updatedUser, err := s.r.UserRepo().UpdateUser(ctx, payload)
	if err != nil {
		slog.Error("error updating user", "error", err)
		return errWrap.WrapError(errConsts.ErrUpdatedUser)
	}

	slog.Debug(success.SuccessResendVerifyAccount)

	// send email verification and email body
site := "http://localhost:3002"
	verificationLink := fmt.Sprintf(site+"/verify-account?verify_token=%s", verifyToken)
	emailBody := fmt.Sprintf("Thank you for registering with Bizpos. Please verify your account by clicking the link below:\n\n%s\n\nThis link will expire in 24 hours.", verificationLink)

	// create send email
	errMail := emails.SendEmail([]string{updatedUser.Email}, nil,
		"Verify Your Account",
		emailBody)
	if errMail != nil {
		slog.Error("error sending email", "error", errMail)
	}

	return nil
}

// Service method for login
func (s *UserService) LoginUser(ctx context.Context, req *userDTO.LoginDTO, clientIP string) (*userDTO.AuthResponse, error) {
	// check email if already exists
	// Validation email
	if req.Email == "" {
		slog.Error("Email is empty")
		return nil, errWrap.WrapError(errConsts.ErrUserEmailInvalid)
	}

	// find user by email
	findUser, err := s.r.UserRepo().FindUserByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("Error finding user by email", "error", err, "email", req.Email)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	if findUser == nil {
		slog.Error("user not found", "email", req.Email)
		return nil, errWrap.WrapError(errConsts.ErrUserNotFound)
	}

	// check user already verified
	if findUser.IsVerified == false {
		slog.Error("user not verified", "email", req.Email)
		return nil, errWrap.WrapError(errConsts.ErrUserNotVerified)
	}

	// verify password
	_, err = config.VerifyPassword(req.Password, findUser.Password)
	if err != nil {
		slog.Error("password is invalid", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrInvalidCredentials)
	}

	// generate expired access token and refresh token
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	// generate access token
	accessToken, err := jwt.GenerateAccessToken(findUser.ID, findUser.RoleName, accessTTL)
	if err != nil {
		slog.Error("failed to generate access token", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	// generate refresh token
	refreshToken, err := jwt.GenerateRefreshToken(findUser.ID, findUser.RoleName, refreshTTL)
	if err != nil {
		slog.Error("failed to generate refresh token", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	tokens := sessiondb.CreateSessionParams{
		UserID:    findUser.ID,
		RefToken:  refreshToken,
		IsBlocked: false,
		ClientIp:  clientIP,
		ExpiresAt: time.Now().Add(refreshTTL),
	}

	// save new refresh token
	_, err = s.r.UserRepo().SaveSession(ctx, tokens)
	if err != nil {
		slog.Error("Failed to save token", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	// add Bearer prefix to access token for client usage
	tokenWithPrefix := fmt.Sprintf("Bearer %s", accessToken)

	slog.Debug(success.SuccessLoginUser)

	// create auth response
	response := &userDTO.AuthResponse{
		ID:           findUser.ID,
		Name:         findUser.Name,
		Email:        findUser.Email,
		IsVerified:   findUser.IsVerified,
		RoleID:       findUser.RoleID,
		Token:        tokenWithPrefix,
		RefreshToken: refreshToken,
	}

	return response, nil
}

// Service method for logout
func (s *UserService) Logout(ctx context.Context, userID uuid.UUID) error {
	// find token by user id
	findToken, err := s.r.UserRepo().FindSessionByUserId(ctx, userID)
	if err != nil {
		slog.Error("Failed to find token", "error", err)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}

	// check if token is nil
	if findToken == nil {
		slog.Error("token not found", "error", err)
		return errWrap.WrapError(errConsts.ErrTokenNotFound)
	}

	// delete token by user id
	err = s.r.UserRepo().DeleteSession(ctx, findToken.RefToken)
	if err != nil {
		slog.Error("Failed to delete token", "error", err)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}

	slog.Debug(success.SuccessLogoutUser)

	return nil
}

// Service method for Get user
func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*userDTO.AuthResponse, error) {
	// find user if exists
	findUser, err := s.r.UserRepo().FindUserById(ctx, userID)
	if err != nil {
		slog.Error("Failed to find user", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrUserNotFound)
	}
	if findUser == nil {
		slog.Error("user not found", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrUserNotFound)
	}

	if findUser.IsVerified == false {
		slog.Error("user not verified", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrUserNotVerified)
	}

	slog.Debug(success.SuccessFindUser)
	// response body
	response := &userDTO.AuthResponse{
		ID:         findUser.ID,
		Name:       findUser.Name,
		Email:      findUser.Email,
		IsVerified: findUser.IsVerified,
		RoleID:     findUser.RoleID,
		RoleDetail: &userDTO.RoleResponse{
			Name:            findUser.RoleName,
			Code:            findUser.RoleCode,
			Level:           findUser.RoleLevel,
			Description:     findUser.RoleDescription,
			CanSelfRegister: findUser.RoleCanSelfRegister,
		},
	}
	return response, nil
}

// Service method for refresh token
func (s *UserService) RefreshToken(ctx context.Context, token string, clientIP string) (*userDTO.AuthResponse, error) {
	// check if refresh token is valid
	findToken, err := s.r.UserRepo().FindSessions(ctx, token)
	if err != nil {
		slog.Error("Failed to get token", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	if findToken == nil {
		slog.Error("token not found", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrTokenNotFound)
	}

	// check if token expired
	if findToken.ExpiresAt.Before(time.Now()) {
		slog.Error("token expired", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrTokenExpired)
	}

	// get user to access role
	findUser, err := s.r.UserRepo().FindUserById(ctx, findToken.UserID)
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	if findUser == nil {
		slog.Error("user not found", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrUserNotFound)
	}

	// generate new access token
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	accessToken, err := jwt.GenerateAccessToken(findUser.ID, findUser.RoleName, accessTTL)
	if err != nil {
		slog.Error("failed to generate access token", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	// generate new refresh token
	newRefreshToken, err := jwt.GenerateRefreshToken(findUser.ID, findUser.RoleName, refreshTTL)
	if err != nil {
		slog.Error("failed to generate refresh token", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	// delete old refresh token
	err = s.r.UserRepo().DeleteSession(ctx, findToken.RefToken)
	if err != nil {
		slog.Error("failed to delete token", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	tokens := sessiondb.CreateSessionParams{
		UserID:    findUser.ID,
		RefToken:  newRefreshToken,
		ClientIp:  clientIP,
		ExpiresAt: time.Now().Add(refreshTTL),
	}

	// save new refresh token
	_, err = s.r.UserRepo().SaveSession(ctx, tokens)
	if err != nil {
		slog.Error("Failed to save token", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	// add Bearer prefix to access token for client usage
	tokenWithPrefix := fmt.Sprintf("Bearer %s", accessToken)

	slog.Debug(success.SuccessLoginUser)

	// create auth response
	response := &userDTO.AuthResponse{
		ID:           findUser.ID,
		Name:         findUser.Name,
		Email:        findUser.Email,
		IsVerified:   findUser.IsVerified,
		RoleID:       findUser.RoleID,
		Token:        tokenWithPrefix,
		RefreshToken: newRefreshToken,
	}

	return response, nil
}

// Service method for forgot password
func (s *UserService) ForgotPassword(ctx context.Context, req *userDTO.EmailRequest) error {
	// check email if already exists
	findUser, err := s.r.UserRepo().FindUserByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("error finding user by email", "error", err)
		return errWrap.WrapError(errConsts.ErrUserEmailNotFound)
	}
	if findUser == nil {
		slog.Error("user already exists", "email", req.Email)
		return errWrap.WrapError(errConsts.ErrUserAlreadyExists)
	}

	if findUser.IsVerified == false {
		slog.Error("user not verified", "email", req.Email)
		return errWrap.WrapError(errConsts.ErrUserNotVerified)
	}

	// generate forgot password token
	forgotToken := config.GenerateRandomToken(32)
	expiredAt := time.Now().Add(time.Hour * 24)

	// update user for reset password
	payload := userdb.UpdateUserParams{
		ID:              findUser.ID,
		Name:            findUser.Name,
		Email:           findUser.Email,
		PhoneNumber:     findUser.PhoneNumber,
		IsVerified:      findUser.IsVerified,
		VerifyCode:      &forgotToken,
		VerifyExpiresAt: &expiredAt,
	}

	_, err = s.r.UserRepo().UpdateUser(ctx, payload)
	if err != nil {
		slog.Error("error updating user", "error", err)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}

	// Send forgot password email
	site := "http://localhost:3002"
	verificationLink := fmt.Sprintf(site+"/reset-password?verify_token=%s", forgotToken)
	emailBody := fmt.Sprintf("You have requested to reset your password. Please click the link below to reset your password:\n\n%s\n\nThis link will expire in 24 hours.", verificationLink)

	// Send email
	errMail := emails.SendEmail([]string{findUser.Email}, nil,
		"Reset Password",
		emailBody)
	if errMail != nil {
		slog.Error("error sending forgot password email", "error", errMail)
		return errWrap.WrapError(errConsts.ErrSQLError)
	} else {
		slog.Debug(success.SuccessEmailSent)
	}

	return nil
}

// Service method for reset password
func (s *UserService) ResetPassword(ctx context.Context, req *userDTO.ResetPasswordDTO) error {
	// check verify code exists
	findVerifyCode, err := s.r.UserRepo().FindUserByVerify(ctx, req.VerifyCode)
	if err != nil {
		if errors.Is(err, errConsts.ErrUserNotFound) {
			return errWrap.WrapError(errConsts.ErrVerifyCodeNotFound)
		}
		slog.Error("error finding verify code", "error", err)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}
	slog.Debug("verify code found", "verify code", findVerifyCode)

	// check verify code expired
	if findVerifyCode.VerifyExpiresAt.Before(time.Now()) {
		slog.Error("verify code expired", "verify code", req.VerifyCode)
		return errWrap.WrapError(errConsts.ErrVerifyCodeExpired)
	}

	// check if password and retry password not match
	if req.Password != req.RetryPassword {
		slog.Error("password and retry password not match", "password", req.Password)
		return errWrap.WrapError(errConsts.ErrUserPasswordNotMatch)
	}

	// hash password
	hashPassword, err := config.HashPassword(req.Password)
	if err != nil {
		slog.Error("error hashing password", "error", err)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}

	// execute update user
	_, err = s.r.UserRepo().UpdateUser(ctx, userdb.UpdateUserParams{
		ID:              findVerifyCode.ID,
		Name:            findVerifyCode.Name,
		Email:           findVerifyCode.Email,
		PhoneNumber:     findVerifyCode.PhoneNumber,
		Password:        hashPassword,
		IsVerified:      true,
		VerifyCode:      nil,
		VerifyExpiresAt: nil,
	})
	if err != nil {
		slog.Error("User update failed", "error", err)
		return errWrap.WrapError(errConsts.ErrUpdatedUser)
	}

	slog.Debug(success.SuccessResetPassword)

	// Send reset success email
	emailBody := "You have successfully reset your password."

	// Send email
	errMail := emails.SendEmail([]string{findVerifyCode.Email}, nil,
		"Reset Password Success",
		emailBody)
	if errMail != nil {
		slog.Error("error sending password reset email", "error", errMail)
		// We don't return error here as password is already reset
	}
	slog.Debug(success.SuccessEmailSent)

	return nil
}
