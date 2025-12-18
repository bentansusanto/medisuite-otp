package treatments

import (
	"net/http"

	treatmentDTO "medisuite-api/app/dto/treatments"
	"medisuite-api/app/services"
	errValidation "medisuite-api/common/errors"
	"medisuite-api/common/response"
	"medisuite-api/constants/success"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ITreatmentHandler interface {
	CreateTreatment(c *gin.Context)
	UpdateTreatment(c *gin.Context)
	DeleteTreatment(c *gin.Context)
	FindAllTreatment(c *gin.Context)
	FindByIdTreatment(c *gin.Context)
}

type TreatmentHandler struct {
	s services.IService
}

func NewTreatmentHandler(s services.IService) ITreatmentHandler {
	return &TreatmentHandler{s: s}
}

// Handler method for creating a new treatment.
func (h *TreatmentHandler) CreateTreatment(c *gin.Context) {
	reqDTO := treatmentDTO.TreatmentDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	// validation request
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

	// execute create treatment service
	result, err := h.s.TreatmentService().Create(c, reqDTO)
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
	resMessage := success.SuccessCreateTreatment
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusCreated,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}

// Handler method for updating a treatment.
func (h *TreatmentHandler) UpdateTreatment(c *gin.Context) {
	reqDTO := treatmentDTO.TreatmentDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	// validation request
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

	id := c.Param("id")
	if id == "" {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code: http.StatusBadRequest,
			Gin:  c,
		})
		return
	}

	// parse id to uuid
	treatmentID, err := uuid.Parse(id)
	if err != nil {
		errMessage := "Invalid treatment ID format"
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// execute update treatment service
	result, err := h.s.TreatmentService().Update(c, treatmentID, reqDTO)
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
	resMessage := success.SuccessUpdateTreatment
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}

// Handler method for deleting a treatment.
func (h *TreatmentHandler) DeleteTreatment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code: http.StatusBadRequest,
			Gin:  c,
		})
		return
	}

	// parse id to uuid
	treatmentID, err := uuid.Parse(id)
	if err != nil {
		errMessage := "Invalid treatment ID format"
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// execute delete treatment service
	err = h.s.TreatmentService().Delete(c, treatmentID)
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
	resMessage := success.SuccessDeleteTreatment
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Gin:     c,
	})
}

// Handler method for finding all treatments.
func (h *TreatmentHandler) FindAllTreatment(c *gin.Context) {
	// execute find all treatment service
	result, err := h.s.TreatmentService().FindAll(c)
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
	resMessage := success.SuccessFindAllTreatment
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}

// Handler method for finding a treatment by ID.
func (h *TreatmentHandler) FindByIdTreatment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code: http.StatusBadRequest,
			Gin:  c,
		})
		return
	}

	// parse id to uuid
	treatmentID, err := uuid.Parse(id)
	if err != nil {
		errMessage := "Invalid treatment ID format"
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// execute find by id treatment service
	result, err := h.s.TreatmentService().FindById(c, treatmentID)
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
	resMessage := success.SuccessFindTreatmentById
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}
