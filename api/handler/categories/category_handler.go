package categories

import (
	"net/http"

	categoryDTO "medisuite-api/app/dto/treatments"
	"medisuite-api/app/services"
	errValidation "medisuite-api/common/errors"
	"medisuite-api/common/response"
	"medisuite-api/constants/success"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ICategoryHandler interface {
	CreateCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	FindAllCategory(c *gin.Context)
	FindByIdCategory(c *gin.Context)
}

type CategoryHandler struct {
	s services.IService
}

func NewCategoryHandler(s services.IService) ICategoryHandler {
	return &CategoryHandler{s: s}
}

// Handler method for creating a new category.
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	reqDTO := categoryDTO.CategoryDTO{}
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

	// execute create category service
	result, err := h.s.CategoryService().Create(c, reqDTO.NameCategory)
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
	resMessage := success.SuccessCreateCategory
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusCreated,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}

// Handler method for updating a category.
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	reqDTO := categoryDTO.CategoryDTO{}
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
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
	categoryID, err := uuid.Parse(id)
	if err != nil {
		errMessage := "Invalid category ID format"
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
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

	// execute update category service
	result, err := h.s.CategoryService().Update(c, categoryID, reqDTO.NameCategory)
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
	resMessage := success.SuccessUpdateCategory
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}

// Handler method for deleting a category.
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code: http.StatusBadRequest,
			Gin:  c,
		})
		return
	}

	// parse id to uuid
	categoryID, err := uuid.Parse(id)
	if err != nil {
		errMessage := "Invalid category ID format"
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// execute delete category service
	err = h.s.CategoryService().Delete(c, categoryID)
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
	resMessage := success.SuccessDeleteCategory
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Gin:     c,
	})
}

// Handler method for finding all categories.
func (h *CategoryHandler) FindAllCategory(c *gin.Context) {
	// execute find all category service
	result, err := h.s.CategoryService().FindAll(c)
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
	resMessage := success.SuccessFindAllCategory
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}

// Handler method for finding a category by ID.
func (h *CategoryHandler) FindByIdCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.HttpResponse(response.ParamHttpResp[any]{
			Code: http.StatusBadRequest,
			Gin:  c,
		})
		return
	}

	// parse id to uuid
	categoryID, err := uuid.Parse(id)
	if err != nil {
		errMessage := "Invalid category ID format"
		response.HttpResponse(response.ParamHttpResp[any]{
			Code:    http.StatusBadRequest,
			Error:   err,
			Message: &errMessage,
			Gin:     c,
		})
		return
	}

	// execute find by id category service
	result, err := h.s.CategoryService().FindById(c, categoryID)
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
	resMessage := success.SuccessFindCategoryById
	response.HttpResponse(response.ParamHttpResp[any]{
		Code:    http.StatusOK,
		Message: &resMessage,
		Data:    result,
		Gin:     c,
	})
}
