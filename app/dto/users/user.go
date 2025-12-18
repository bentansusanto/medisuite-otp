package users

import (
	"time"

	"github.com/google/uuid"
)

type RegisterDTO struct {
	Name        string    `json:"name" validation:"required"`
	Email       string    `json:"email" validation:"required, email"`
	Password    string    `json:"password" validation:"required min=8 max=20 alphaNum cap special"`
	PhoneNumber string    `json:"phone_number" validation:"required"`
	RoleID      uuid.UUID `json:"role_id" validation:"required"`
}

type AuthResponse struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	PhoneNumber  string        `json:"phone_number"`
	IsVerified   bool          `json:"is_verified"`
	RoleID       uuid.UUID     `json:"role_id"`
	RoleDetail   *RoleResponse `json:"role_detail,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Token        string        `json:"token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
}

type RoleResponse struct {
	Name            string `json:"name"`
	Code            string `json:"code"`
	Level           int32  `json:"level"`
	Description     string `json:"description"`
	CanSelfRegister bool   `json:"can_self_register"`
}

type LoginDTO struct {
	Email    string `json:"email" validation:"required"`
	Password string `json:"password" validation:"required min=8 max=20 alphaNum cap special"`
}

type EmailRequest struct {
	Email string `json:"email" validation:"required"`
}

type ResetPasswordDTO struct {
	Password      string `json:"password" validation:"required min=8 max=20 alphaNum cap special"`
	RetryPassword string `json:"retry_password" validation:"required min=8 max=20 alphaNum cap special"`
	VerifyCode    string `json:"verify_code" validation:"required"`
}
