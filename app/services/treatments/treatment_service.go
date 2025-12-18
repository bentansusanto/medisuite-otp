package treatments

import (
	"context"
	"log/slog"

	treatmentDTO "medisuite-api/app/dto/treatments"
	"medisuite-api/app/repo"
	errWrap "medisuite-api/common/errors"
	errConsts "medisuite-api/constants/errors"
	"medisuite-api/constants/success"
	treatmentdb "medisuite-api/pkg/db/treatments"
	"github.com/google/uuid"
)

type ITreatmentService interface {
	Create(ctx context.Context, req treatmentDTO.TreatmentDTO) (*treatmentDTO.TreatmentResponse, error)
	Update(ctx context.Context, treatmentID uuid.UUID, req treatmentDTO.TreatmentDTO) (*treatmentDTO.TreatmentResponse, error)
	Delete(ctx context.Context, treatmentID uuid.UUID) error
	FindAll(ctx context.Context) ([]treatmentDTO.TreatmentResponse, error)
	FindById(ctx context.Context, treatmentID uuid.UUID) (*treatmentDTO.TreatmentResponse, error)
}

type TreatmentService struct {
	r repo.IRepo
}

func NewTreatmentService(r repo.IRepo) ITreatmentService {
	return &TreatmentService{r: r}
}

// Services method for creating a new treatment.
func (s *TreatmentService) Create(ctx context.Context, req treatmentDTO.TreatmentDTO) (*treatmentDTO.TreatmentResponse, error) {
	// find treatment if exist
	treatment, err := s.r.TreatmentRepo().FindByName(ctx, req.NameTreatment)
	if err != nil {
		slog.Error("Error finding name_treatment", "error", err, "name_treatment", req.NameTreatment)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	if treatment != nil {
		slog.Error("Treatment already exist", "error", err, "name_treatment", req.NameTreatment)
		return nil, errWrap.WrapError(errConsts.ErrTreatmentExist)
	}

	category, err := s.r.CategoryRepo().FindById(ctx, req.CategoryID)
	if err != nil {
		slog.Error("Error finding category", "error", err, "category_id", req.CategoryID)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	if category == nil {
		slog.Error("Category not found", "error", err, "category_id", req.CategoryID)
		return nil, errWrap.WrapError(errConsts.ErrFindCategoryId)
	}

	payload := treatmentdb.CreateTreatmentParams{
		CategoryID:    req.CategoryID,
		NameTreatment: req.NameTreatment,
		Description:   req.Description,
		Thumbnail:     req.Thumbnail,
		Price:         req.Price,
		Duration:      req.Duration,
		IsActive:      req.IsActive,
	}

	// create treatment
	treatment, err = s.r.TreatmentRepo().Create(ctx, payload)
	if err != nil {
		slog.Error("Error creating treatment", "error", err, "name_treatment", req.NameTreatment)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	slog.Info(success.SuccessCreateTreatment)

	response := &treatmentDTO.TreatmentResponse{
		ID:            treatment.ID,
		CategoryID:    treatment.CategoryID,
		NameTreatment: treatment.NameTreatment,
		Description:   treatment.Description,
		Thumbnail:     treatment.Thumbnail,
		Price:         treatment.Price,
		Duration:      treatment.Duration,
		IsActive:      treatment.IsActive,
		CreatedAt:     treatment.CreatedAt,
		UpdatedAt:     treatment.UpdatedAt,
	}

	return response, nil
}

// Services method for updating a treatment.
func (s *TreatmentService) Update(ctx context.Context, treatmentID uuid.UUID, req treatmentDTO.TreatmentDTO) (*treatmentDTO.TreatmentResponse, error) {
	// find treatment if exist
	findTreatment, err := s.r.TreatmentRepo().FindById(ctx, treatmentID)
	if err != nil {
		slog.Error("Error finding treatment", "error", err, "treatment_id", treatmentID)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	if findTreatment == nil {
		slog.Error("Treatment not found", "error", err, "treatment_id", treatmentID)
		return nil, errWrap.WrapError(errConsts.ErrFindTreatmentId)
	}

	// find category if exist
	category, err := s.r.CategoryRepo().FindById(ctx, req.CategoryID)
	if err != nil {
		slog.Error("Error finding category", "error", err, "category_id", req.CategoryID)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	if category == nil {
		slog.Error("Category not found", "error", err, "category_id", req.CategoryID)
		return nil, errWrap.WrapError(errConsts.ErrFindCategoryId)
	}

	payload := treatmentdb.UpdateTreatmentParams{
		ID:            treatmentID,
		CategoryID:    req.CategoryID,
		NameTreatment: req.NameTreatment,
		Description:   req.Description,
		Thumbnail:     req.Thumbnail,
		Price:         req.Price,
		Duration:      req.Duration,
		IsActive:      req.IsActive,
	}

	// update treatment
	updateTreatment, err := s.r.TreatmentRepo().Update(ctx, payload)
	if err != nil {
		slog.Error("Error updating treatment", "error", err, "treatment_id", treatmentID)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	response := &treatmentDTO.TreatmentResponse{
		ID:            updateTreatment.ID,
		CategoryID:    updateTreatment.CategoryID,
		NameTreatment: updateTreatment.NameTreatment,
		Description:   updateTreatment.Description,
		Thumbnail:     updateTreatment.Thumbnail,
		Price:         updateTreatment.Price,
		Duration:      updateTreatment.Duration,
		IsActive:      updateTreatment.IsActive,
		UpdatedAt:     updateTreatment.UpdatedAt,
	}

	return response, nil
}

// Services method for deleting a treatment.
func (s *TreatmentService) Delete(ctx context.Context, treatmentID uuid.UUID) error {
	// find treatment if exist
	findTreatment, err := s.r.TreatmentRepo().FindById(ctx, treatmentID)
	if err != nil {
		slog.Error("Error finding treatment", "error", err, "treatment_id", treatmentID)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}

	if findTreatment == nil {
		slog.Error("Treatment not found", "error", err, "treatment_id", treatmentID)
		return errWrap.WrapError(errConsts.ErrFindTreatmentId)
	}

	// delete treatment
	err = s.r.TreatmentRepo().Delete(ctx, treatmentID)
	if err != nil {
		slog.Error("Error deleting treatment", "error", err, "treatment_id", treatmentID)
		return errWrap.WrapError(errConsts.ErrSQLError)
	}

	slog.Info(success.SuccessDeleteTreatment)
	return nil
}

// Services method for finding all treatments.
func (s *TreatmentService) FindAll(ctx context.Context) ([]treatmentDTO.TreatmentResponse, error) {
	// find treatment if exist
	treatments, err := s.r.TreatmentRepo().FindAll(ctx)
	if err != nil {
		slog.Error("Error finding treatment", "error", err)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	response := make([]treatmentDTO.TreatmentResponse, len(treatments))
	for i, t := range treatments {
		response[i] = treatmentDTO.TreatmentResponse{
			ID:            t.ID,
			CategoryID:    t.CategoryID,
			NameTreatment: t.NameTreatment,
			Description:   t.Description,
			Thumbnail:     t.Thumbnail,
			Price:         t.Price,
			Duration:      t.Duration,
			IsActive:      t.IsActive,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
		}
	}

	return response, nil
}

// Services method for finding a treatment by ID.
func (s *TreatmentService) FindById(ctx context.Context, treatmentID uuid.UUID) (*treatmentDTO.TreatmentResponse, error) {
	// find treatment if exist
	findTreatment, err := s.r.TreatmentRepo().FindById(ctx, treatmentID)
	if err != nil {
		slog.Error("Error finding treatment", "error", err, "treatment_id", treatmentID)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}

	if findTreatment == nil {
		slog.Error("Treatment not found", "error", err, "treatment_id", treatmentID)
		return nil, errWrap.WrapError(errConsts.ErrFindTreatmentId)
	}

	response := treatmentDTO.TreatmentResponse{
		ID:            findTreatment.ID,
		CategoryID:    findTreatment.CategoryID,
		NameTreatment: findTreatment.NameTreatment,
		Description:   findTreatment.Description,
		Thumbnail:     findTreatment.Thumbnail,
		Price:         findTreatment.Price,
		Duration:      findTreatment.Duration,
		IsActive:      findTreatment.IsActive,
		CreatedAt:     findTreatment.CreatedAt,
		UpdatedAt:     findTreatment.UpdatedAt,
	}

	return &response, nil
}
