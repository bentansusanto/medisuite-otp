package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"medisuite-api/app/repo"
	errWrap "medisuite-api/common/errors"
	"medisuite-api/common/response"
	errConstants "medisuite-api/constants/errors"
	constants "medisuite-api/constants/http_status"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Recovered from panic", r)
				c.JSON(http.StatusInternalServerError, response.Response[any]{
					Status:  constants.Error,
					Message: errConstants.ErrInternalServerError.Error(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:  http.StatusUnauthorized,
				Error: errWrap.WrapError(errConstants.ErrUnauthorized),
				Gin:   c,
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No Bearer prefix found
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:  http.StatusUnauthorized,
				Error: errWrap.WrapError(errConstants.ErrUnauthorized),
				Gin:   c,
			})
			c.Abort()
			return
		}

		// Verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:  http.StatusUnauthorized,
				Error: errWrap.WrapError(errConstants.ErrUnauthorized),
				Gin:   c,
			})
			c.Abort()
			return
		}

		// Set user info in context if needed
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if rawID, ok := claims["user_id"]; ok {
				// JWT library typically unmarshals UUIDs as strings
				if idStr, ok := rawID.(string); ok {
					if userID, err := uuid.Parse(idStr); err == nil {
						c.Set("userID", userID)
					} else {
						slog.Error("invalid user_id in token claims", "value", idStr, "error", err)
					}
				}
			}
			// Extract role from JWT claims
			if rawRole, ok := claims["role"]; ok {
				if roleStr, ok := rawRole.(string); ok {
					c.Set("roleCode", roleStr)
				}
			}
		}

		c.Next()
	}
}

// Rate limiter middleware
func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	// return rate limiter middleware and check if request over limit
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:  http.StatusTooManyRequests,
				Error: errWrap.WrapError(errConstants.ErrTooManyRequests),
				Gin:   c,
			})
			c.Abort()
			return // Important: return early when rate limited
		}
		c.Next()
	}
}

// RequirePermission checks if the authenticated user has the required permission
// Usage: router.POST("/users", AuthMiddleware(), RequirePermission(repo, "user", "create"), handler)
func RequirePermission(repository repo.IRepo, module, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get userID from context (set by AuthMiddleware)
		userIDVal, exists := c.Get("userID")
		if !exists {
			errMessage := "User not authenticated"
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusUnauthorized,
				Error:   errWrap.WrapError(errConstants.ErrUnauthorized),
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			errMessage := "Invalid user ID"
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusUnauthorized,
				Error:   errWrap.WrapError(errConstants.ErrUnauthorized),
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		// Get user to retrieve role_id
		user, err := repository.UserRepo().FindUserById(c.Request.Context(), userID)
		if err != nil {
			errMessage := "Failed to get user information"
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusInternalServerError,
				Error:   err,
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		// Get role permissions
		permissions, err := repository.RolePermissionRepo().GetRolePermissions(c.Request.Context(), user.RoleID)
		if err != nil {
			errMessage := "Failed to check permissions"
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusInternalServerError,
				Error:   err,
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		// Check if user has the required permission
		hasPermission := false
		for _, perm := range permissions {
			if perm.Module == module && perm.Action == action {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			errMessage := "You do not have permission to perform this action"
			slog.Warn("Permission denied", "user_id", userID, "module", module, "action", action)
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusForbidden,
				Error:   errWrap.WrapError(errConstants.ErrInsufficientPermissions),
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole checks if the authenticated user has one of the allowed roles
// Usage: router.GET("/admin/dashboard", AuthMiddleware(), RequireRole(repo, "owner", "admin"), handler)
func RequireRole(repository repo.IRepo, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get userID from context (set by AuthMiddleware)
		userIDVal, exists := c.Get("userID")
		if !exists {
			errMessage := "User not authenticated"
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusUnauthorized,
				Error:   errWrap.WrapError(errConstants.ErrUnauthorized),
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			errMessage := "Invalid user ID"
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusUnauthorized,
				Error:   errWrap.WrapError(errConstants.ErrUnauthorized),
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		// Get user to retrieve role_id
		user, err := repository.UserRepo().FindUserById(c.Request.Context(), userID)
		if err != nil {
			errMessage := "Failed to get user information"
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusInternalServerError,
				Error:   err,
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		// Get role information
		role, err := repository.RoleRepo().FindRoleById(context.Background(), user.RoleID)
		if err != nil {
			errMessage := "Failed to get role information"
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusInternalServerError,
				Error:   err,
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		// Check if user's role is in the allowed roles
		roleAllowed := false
		roleCode := role.Code
		for _, allowedRole := range allowedRoles {
			if roleCode == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			errMessage := "Your role does not have access to this resource"
			slog.Warn("Role access denied", "user_id", userID, "role", roleCode, "allowed_roles", allowedRoles)
			response.HttpResponse(response.ParamHttpResp[any]{
				Code:    http.StatusForbidden,
				Error:   errWrap.WrapError(errConstants.ErrInvalidRole),
				Message: &errMessage,
				Gin:     c,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
