package errors

import "errors"

var (
	ErrFindCategories     = errors.New("Error find categories")
	ErrCreateCategory     = errors.New("Error create category")
	ErrUpdateCategory     = errors.New("Error update category")
	ErrDeleteCategory     = errors.New("Error delete category")
	ErrFindCategoryId     = errors.New("Error find category id")
	ErrCategoryExist      = errors.New("Category exist")
	ErrFindCategoryByName = errors.New("Error find category by name")
	ErrFindTreatment      = errors.New("Error find treatment")
	ErrCreateTreatment    = errors.New("Error create treatment")
	ErrUpdateTreatment    = errors.New("Error update treatment")
	ErrDeleteTreatment    = errors.New("Error delete treatment")
	ErrFindTreatmentId    = errors.New("Error find treatment id")
	ErrFindTreatmentByName = errors.New("Error find treatment by name")
	ErrTreatmentExist     = errors.New("Treatment exist")
)

var ServiceErrorMessage = []error{
	ErrFindCategories,
	ErrCreateCategory,
	ErrUpdateCategory,
	ErrDeleteCategory,
	ErrFindCategoryId,
	ErrCategoryExist,
	ErrFindCategoryByName,
	ErrFindTreatment,
	ErrCreateTreatment,
	ErrUpdateTreatment,
	ErrDeleteTreatment,
	ErrFindTreatmentId,
	ErrFindTreatmentByName,
	ErrTreatmentExist,
}
