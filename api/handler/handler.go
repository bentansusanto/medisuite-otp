package handler

import (
	categoryHandler "medisuite-api/api/handler/categories"
	treatmentHandler "medisuite-api/api/handler/treatments"
	userHandler "medisuite-api/api/handler/users"
	"medisuite-api/app/services"
)

type IHandler interface {
	UserHandler() userHandler.IUserHandler
	CategoryHandler() categoryHandler.ICategoryHandler
	TreatmentHandler() treatmentHandler.ITreatmentHandler
}

type Handler struct {
	s services.IService
}

func NewHandler(s services.IService) IHandler {
	return &Handler{s: s}
}

func (h *Handler) UserHandler() userHandler.IUserHandler {
	return userHandler.NewUserHandler(h.s)
}

func (h *Handler) CategoryHandler() categoryHandler.ICategoryHandler {
	return categoryHandler.NewCategoryHandler(h.s)
}

func (h *Handler) TreatmentHandler() treatmentHandler.ITreatmentHandler {
	return treatmentHandler.NewTreatmentHandler(h.s)
}
