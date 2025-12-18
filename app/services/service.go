package services

import (
	"medisuite-api/app/repo"
	categoryService "medisuite-api/app/services/categories"
	treatmentService "medisuite-api/app/services/treatments"
	userService "medisuite-api/app/services/users"
)

type IService interface {
	UserService() userService.IUserService
	CategoryService() categoryService.ICategoryService
	TreatmentService() treatmentService.ITreatmentService
}

type Service struct {
	r repo.IRepo
}

func NewService(r repo.IRepo) IService {
	return &Service{r: r}
}

func (s *Service) UserService() userService.IUserService {
	return userService.NewUserService(s.r)
}

func (s *Service) CategoryService() categoryService.ICategoryService {
	return categoryService.NewCategoryService(s.r)
}

func (s *Service) TreatmentService() treatmentService.ITreatmentService {
	return treatmentService.NewTreatmentService(s.r)
}
