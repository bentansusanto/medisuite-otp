package users

import (
	"errors"
	"net/http"

	userDTO "medisuite-api/app/dto/users"
	"medisuite-api/app/services"
	errValidation "medisuite-api/common/errors"
	"medisuite-api/common/response"
	"medisuite-api/constants/success"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IUserHandler interface {
	Register(c *gin.Context)
	VerifyAccount(c *gin.Context)
	ResendVerify(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetUser(c *gin.Context)
	RefreshToken(c *gin.Context)
	ForgotPassword(c *gin.Context)
	ResetPassword(c *gin.Context)
}

type UserHandler struct {
	s services.IService
}

func NewUserHandler(s services.IService) IUserHandler {
	return &UserHandler{s: s}
}

// Handler method for user registration
func (h *UserHandler) Register(c *gin.Context) {
	reqDTO := userDTO.RegisterDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	// validate request
	validate := validator.New()
	if err := validate.Struct(reqDTO); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	// execute register service
	result, err := h.s.UserService().CreateUser(c, reqDTO)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// return success response (without token for registration)
	resMessage := success.SuccessCreateUser
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusCreated,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}

// Handler method for verify account
func (h *UserHandler) VerifyAccount(c *gin.Context) {
	token := c.Query("verify_token")

	// check if token is empty
	if token == "" {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: errors.New("token is required"),
			Gin:   c,
		})
		return
	}

	// check if token is valid
	result, err := h.s.UserService().VerifyAccount(c, token)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// return success response (without token for verification)
	resMessage := success.SuccessVerifyUser
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}

// Handler method for resend verify account
func (h *UserHandler) ResendVerify(c *gin.Context) {
	reqDTO := &userDTO.EmailRequest{}
	if err := c.ShouldBindJSON(reqDTO); err != nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	// validate request
	validate := validator.New()
	if err := validate.Struct(reqDTO); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	// execute resend verify account service
	err := h.s.UserService().ResendVerifyAccount(c, reqDTO.Email)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// return success response
	resMessage := success.SuccessResendVerifyAccount
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Gin:     c,
	})
}

// Handler method for login
func (h *UserHandler) Login(c *gin.Context) {
	reqDTO := &userDTO.LoginDTO{}
	if err := c.ShouldBindJSON(reqDTO); err != nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	// validate request
	validate := validator.New()
	if err := validate.Struct(reqDTO); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}
	// execute login service (pass client IP for session tracking)
	clientIP := c.ClientIP()
	result, err := h.s.UserService().LoginUser(c, reqDTO, clientIP)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// set refresh token to cookie
	cookieSecure := false
	cookieDomain := ""

	// in production, use secure cookies and proper domain
	if gin.Mode() == gin.ReleaseMode {
		cookieSecure = true
	}

	c.SetCookie(
		"refresh_token",     // name
		result.RefreshToken, // value
		7*24*60*60,          // maxAge (7 days in seconds)
		"/",                 // path
		cookieDomain,        // domain
		cookieSecure,        // secure (true in production)
		true,                // httpOnly
	)

	// return success response (only access token in body)
	dataResult := map[string]any{
		"id":       result.ID,
		"name":     result.Name,
		"email":    result.Email,
		"roleid":   result.RoleID,
		"verified": result.IsVerified,
		"token":    result.Token,
	}
	resMessage := "User logged in successfully"
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    dataResult,
		Gin:     c,
	})
}

// Handler method for logout
func (h *UserHandler) Logout(c *gin.Context) {
	// execute logout service
	userIDs, exists := c.Get("userID")
	if !exists {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: errors.New("user ID not found"),
			Gin:   c,
		})
		return
	}

	// type asset userIDs to string
	userID, ok := userIDs.(uuid.UUID)
	if !ok {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: errors.New("invalid user ID format"),
			Gin:   c,
		})
		return
	}

	// validate user ID is not empty
	if userID == uuid.Nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: errors.New("user ID cannot be empty"),
			Gin:   c,
		})
		return
	}

	// execute logout service
	err := h.s.UserService().Logout(c, userID)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// clear refresh token cookie
	cookieSecure := false
	cookieDomain := ""

	// in production, use secure cookies and proper domain
	if gin.Mode() == gin.ReleaseMode {
		cookieSecure = true
		// cookieDomain can be set from config if needed
	}

	c.SetCookie(
		"refresh_token", // name
		"",              // value (empty to clear)
		-1,              // maxAge (negative to delete)
		"/",             // path
		cookieDomain,    // domain
		cookieSecure,    // secure (true in production)
		true,            // httpOnly
	)
	// clear refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	// return success response
	resMessage := success.SuccessLogoutUser
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    nil,
		Gin:     c,
	})
}

// Handler method for get user
func (h *UserHandler) GetUser(c *gin.Context) {
	// execute get user service
	userIDs, exists := c.Get("userID")
	if !exists {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: errors.New("user ID not found"),
			Gin:   c,
		})
		return
	}

	// type asset userIDs to string
	userID, ok := userIDs.(uuid.UUID)
	if !ok {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: errors.New("invalid user ID format"),
			Gin:   c,
		})
		return
	}

	// validate user ID is not empty
	if userID == uuid.Nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: errors.New("user ID cannot be empty"),
			Gin:   c,
		})
		return
	}

	// execute get user service
	user, err := h.s.UserService().GetUser(c, userID)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// return success response (without token)
	resMessage := success.SuccessFindUser
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data: map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.RoleDetail,
			"verified": user.IsVerified,
		},
		Gin: c,
	})
}

// Handler method for refresh token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	// get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		errMessage := "Refresh token not found in cookie"
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusUnauthorized,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// execute refresh token service
	clientIP := c.ClientIP()
	result, err := h.s.UserService().RefreshToken(c, refreshToken, clientIP)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusUnauthorized,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// set new refresh token to cookie
	cookieSecure := false
	cookieDomain := ""

	// in production, use secure cookies and proper domain
	if gin.Mode() == gin.ReleaseMode {
		cookieSecure = true
	}

	c.SetCookie(
		"refresh_token",     // name
		result.RefreshToken, // value
		7*24*60*60,          // maxAge (7 days in seconds)
		"/",                 // path
		cookieDomain,        // domain
		cookieSecure,        // secure (true in production)
		true,                // httpOnly
	)

	// return success response (only access token in body)
	dataResult := map[string]any{
		"id":       result.ID,
		"name":     result.Name,
		"email":    result.Email,
		"role_id":  result.RoleID,
		"verified": result.IsVerified,
		"token":    result.Token,
	}
	resMessage := success.SuccessRefreshToken
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    dataResult,
		Gin:     c,
	})
}

// Handler method for forgot password
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	reqDTO := &userDTO.EmailRequest{}
	if err := c.ShouldBindJSON(reqDTO); err != nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}
	// validate request
	validate := validator.New()
	if err := validate.Struct(reqDTO); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}
	// execute forgot password service
	err := h.s.UserService().ForgotPassword(c, reqDTO)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}
	// return success response
	resMessage := success.SuccessSendEmailForgotPassword
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Gin:     c,
	})
}

// Handler method for reset password
func (h *UserHandler) ResetPassword(c *gin.Context) {
	token := c.Query("verify_token")

	// check if token is empty
	if token == "" {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: errors.New("token is required"),
			Gin:   c,
		})
		return
	}

	reqDTO := &userDTO.ResetPasswordDTO{}
	if err := c.ShouldBindJSON(reqDTO); err != nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}
	// validate request
	validate := validator.New()
	if err := validate.Struct(reqDTO); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	// execute reset password service
	err := h.s.UserService().ResetPassword(c, reqDTO)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}
	// return success response
	resMessage := success.SendEmailSuccessResetPassword
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Gin:     c,
	})
}
