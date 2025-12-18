package categories

import (
	"context"
	"log/slog"

	categoryDTO "medisuite-api/app/dto/treatments"
	"medisuite-api/app/repo"
	errWrap "medisuite-api/common/errors"
	errConsts "medisuite-api/constants/errors"
	"medisuite-api/constants/success"

	"github.com/google/uuid"
)

type ICategoryService interface {
	Create(ctx context.Context, name_category string) (*categoryDTO.CategoryResponse, error)
	Update(ctx context.Context, categoryID uuid.UUID, name_category string) (*categoryDTO.CategoryResponse, error)
	Delete(ctx context.Context, categoryID uuid.UUID) error
	FindAll(ctx context.Context) ([]categoryDTO.CategoryResponse, error)
	FindById(ctx context.Context, categoryID uuid.UUID) (*categoryDTO.CategoryResponse, error)
}

type CategoryService struct {
	r repo.IRepo
}

func NewCategoryService(r repo.IRepo) ICategoryService {
	return &CategoryService{r: r}
}

// Services method for creating a new category.
func (s *CategoryService) Create(ctx context.Context, name_category string) (*categoryDTO.CategoryResponse, error) {
	// find category if exist
	findCategory, err := s.r.CategoryRepo().FindCategoryByName(ctx, name_category)
	if err != nil {
		slog.Error("Error finding name_category", "error", err, "name_category", name_category)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	if findCategory != nil {
		slog.Error("Category already exist", "error", err, "name_category", name_category)
		return nil, errWrap.WrapError(errConsts.ErrCategoryExist)
	}

	// create category
	category, err := s.r.CategoryRepo().Create(ctx, name_category)
	if err != nil {
		slog.Error("Error creating category", "error", err, "name_category", name_category)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	slog.Info(success.SuccessCreateCategory)

	response := &categoryDTO.CategoryResponse{
		ID:           category.ID,
		NameCategory: category.NameCategory,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}

	return response, nil
}

// Services method for updating a category.
func (s *CategoryService) Update(ctx context.Context, categoryID uuid.UUID, name_category string) (*categoryDTO.CategoryResponse, error) {
	// find category if exist
	findCategory, err := s.r.CategoryRepo().FindById(ctx, categoryID)
	if err != nil {
		slog.Error("Error finding name_category", "error", err, "name_category", name_category)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	if findCategory == nil {
		slog.Error("Category not found", "error", err, "name_category", name_category)
		return nil, errWrap.WrapError(errConsts.ErrFindCategoryId)
	}

	// update category
	category, err := s.r.CategoryRepo().Update(ctx, categoryID, name_category)
	if err != nil {
		slog.Error("Error updating category", "error", err, "name_category", name_category)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	slog.Info(success.SuccessUpdateCategory)

	response := &categoryDTO.CategoryResponse{
		ID:           category.ID,
		NameCategory: category.NameCategory,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}

	return response, nil
}

// Services method for deleting a category.
func (s *CategoryService) Delete(ctx context.Context, categoryID uuid.UUID) error {
	// find category if exist
	findCategory, err := s.r.CategoryRepo().FindById(ctx, categoryID)
	if err != nil {
		slog.Error("Error finding name_category", "error", err)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}
	if findCategory == nil {
		slog.Error("Category not found", "error", err)
		return errWrap.WrapError(errConsts.ErrFindCategoryId)
	}

	// delete category
	err = s.r.CategoryRepo().Delete(ctx, categoryID)
	if err != nil {
		return errWrap.WrapError(errConsts.ErrSQLError)
	}

	slog.Info(success.SuccessDeleteCategory)

	return nil
}

// Services method for finding all categories.
func (s *CategoryService) FindAll(ctx context.Context) ([]categoryDTO.CategoryResponse, error) {
	// find category if exist
	categories, err := s.r.CategoryRepo().FindAll(ctx)
	if err != nil {
		slog.Error("Error finding name_category", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	response := make([]categoryDTO.CategoryResponse, len(categories))
	for i, c := range categories {
		response[i] = categoryDTO.CategoryResponse{
			ID:           c.ID,
			NameCategory: c.NameCategory,
			CreatedAt:    c.CreatedAt,
			UpdatedAt:    c.UpdatedAt,
		}
	}
	return response, nil
}

// Services method for finding a category by ID.
func (s *CategoryService) FindById(ctx context.Context, categoryID uuid.UUID) (*categoryDTO.CategoryResponse, error) {
	// find category if exist
	findCategory, err := s.r.CategoryRepo().FindById(ctx, categoryID)
	if err != nil {
		slog.Error("Error finding name_category", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	if findCategory == nil {
		slog.Error("Category not found", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrFindCategoryId)
	}

	response := &categoryDTO.CategoryResponse{
		ID:           findCategory.ID,
		NameCategory: findCategory.NameCategory,
		CreatedAt:    findCategory.CreatedAt,
		UpdatedAt:    findCategory.UpdatedAt,
	}
	return response, nil
}
