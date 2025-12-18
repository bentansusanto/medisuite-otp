package categories

import (
	"context"
	"log/slog"

	errWrap "medisuite-api/common/errors"
	errConsts "medisuite-api/constants/errors"
	categorydb "medisuite-api/pkg/db/categories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ICategoryRepo interface {
	Create(ctx context.Context, name_category string) (*categorydb.Category, error)
	Update(ctx context.Context, id uuid.UUID, name_category string) (*categorydb.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context) ([]categorydb.Category, error)
	FindById(ctx context.Context, id uuid.UUID) (*categorydb.Category, error)
	FindCategoryByName(ctx context.Context, name_category string) (*categorydb.Category, error)
}

type CategoryRepo struct {
	cq *categorydb.Queries
}

func NewCategoryRepo(cq *categorydb.Queries) ICategoryRepo {
	return &CategoryRepo{cq: cq}
}

// Repository method for creating a new category.
func (r *CategoryRepo) Create(ctx context.Context, nameCategory string) (*categorydb.Category, error) {
	category, err := r.cq.CreateCategory(ctx, nameCategory)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Error finding name_category", "error", err, "name_category", nameCategory)
			return nil, nil
		}
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &category, nil
}

// Repository method for updating a category.
func (r *CategoryRepo) Update(ctx context.Context, id uuid.UUID, nameCategory string) (*categorydb.Category, error) {
	category, err := r.cq.UpdateCategory(ctx, id, nameCategory)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Error finding name_category", "error", err, "name_category", nameCategory)
			return nil, nil
		}
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &category, nil
}

// Repository method for deleting a category.
func (r *CategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.cq.DeleteCategory(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Error finding name_category", "error", err, "name_category", id)
			return errWrap.WrapError(errConsts.ErrSQLError)
		}
		return errWrap.WrapError(errConsts.ErrSQLError)
	}
	return nil
}

// Repository method for finding all categories.
func (r *CategoryRepo) FindAll(ctx context.Context) ([]categorydb.Category, error) {
	categories, err := r.cq.FindCategories(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Error finding name_category", "error", err)
			return nil, nil
		}
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return categories, nil
}

// Repository method for finding a category by ID.
func (r *CategoryRepo) FindById(ctx context.Context, id uuid.UUID) (*categorydb.Category, error) {
	category, err := r.cq.FindCategoryById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Error finding name_category", "error", err)
			return nil, nil
		}
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &category, nil
}

// Repository method for finding a category by name.
func (r *CategoryRepo) FindCategoryByName(ctx context.Context, nameCategory string) (*categorydb.Category, error) {
	category, err := r.cq.FindCategoryByName(ctx, nameCategory)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("Error finding name_category", "error", err)
			return nil, nil
		}
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &category, nil
}
