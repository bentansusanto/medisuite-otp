package response

import (
	"net/http"

	httpstatus "medisuite-api/constants/http_status"

	errConsts "medisuite-api/constants/errors"

	"github.com/gin-gonic/gin"
)

// response data
type Response[T any] struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Data    T       `json:"data"`
	Token   *string `json:"token,omitempty"`
}

// detail response field data
type ParamHttpResp[T any] struct {
	Code    int
	Error   error
	Data    T
	Message *string
	Gin     *gin.Context
	// Token   string
}

/*
	HttpResponse is a generic function that handles HTTP responses in a standardized format
	It takes a single parameter of type ParamHttpResp[T] which contains all necessary response parameters
**/

func HttpResponse[T any](param ParamHttpResp[T]) {
	// Check if there's no error in the response
	if param.Error == nil {
		/*
			If no error, send a successful response
			The response includes:
			- Status: "success" (from constants.Success)
			- Message: Custom message if provided, otherwise default HTTP status text
			- Data: The actual data payload of type T
			- Token: Optional JWT token (if provided)
		**/
		message := http.StatusText(param.Code)
		if param.Message != nil {
			message = *param.Message
		}
		param.Gin.JSON(param.Code, Response[T]{
			Status:  httpstatus.Success,
			Message: message,
			Data:    param.Data,
			// Token:   &param.Token,
		})
		return
	}

	/*
		If we reach here, there was an error
		Set default error message
	**/
	message := errConsts.ErrInternalServerError.Error()
	// Prioritize explicit message first
	if param.Message != nil {
		message = *param.Message
	} else if param.Error != nil {
		// Check if error is in our mapped errors
		if errConsts.MappingError(param.Error) {
			message = param.Error.Error()
		}
	}

	param.Gin.JSON(param.Code, Response[T]{
		Status:  httpstatus.Error,
		Message: message,
		Data:    param.Data,
	})
}
