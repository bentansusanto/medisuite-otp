package treatments

import (
	"context"

	errWrap "medisuite-api/common/errors"
	errConsts "medisuite-api/constants/errors"
	treatmentdb "medisuite-api/pkg/db/treatments"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ITreatmentRepo interface {
	Create(ctx context.Context, req treatmentdb.CreateTreatmentParams) (*treatmentdb.Treatment, error)
	Update(ctx context.Context, req treatmentdb.UpdateTreatmentParams) (*treatmentdb.Treatment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context) ([]treatmentdb.Treatment, error)
	FindById(ctx context.Context, id uuid.UUID) (*treatmentdb.Treatment, error)
	FindByName(ctx context.Context, name string) (*treatmentdb.Treatment, error)
}

type TreatmentRepo struct {
	tq *treatmentdb.Queries
}

func NewTreatmentRepo(tq *treatmentdb.Queries) ITreatmentRepo {
	return &TreatmentRepo{tq: tq}
}

// Repository method for creating a new service.
func (r *TreatmentRepo) Create(ctx context.Context, req treatmentdb.CreateTreatmentParams) (*treatmentdb.Treatment, error) {
	treatment, err := r.tq.CreateTreatment(ctx, req)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &treatment, nil
}

// Repository method for updating a service.
func (r *TreatmentRepo) Update(ctx context.Context, req treatmentdb.UpdateTreatmentParams) (*treatmentdb.Treatment, error) {
	treatment, err := r.tq.UpdateTreatment(ctx, req)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &treatment, nil
}

// Repository method for deleting a service.
func (r *TreatmentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.tq.DeleteTreatment(ctx, id)
	if err != nil {
		return errWrap.WrapError(errConsts.ErrSQLError)
	}
	return nil
}

// Repository method for finding all services.
func (r *TreatmentRepo) FindAll(ctx context.Context) ([]treatmentdb.Treatment, error) {
	treatments, err := r.tq.FindTreatments(ctx)
	if err != nil {
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return treatments, nil
}

// Repository method for finding a service by ID.
func (r *TreatmentRepo) FindById(ctx context.Context, id uuid.UUID) (*treatmentdb.Treatment, error) {
	treatment, err := r.tq.FindTreatmentById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			// treatment not found is not an error here, return nil
			return nil, nil
		}
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &treatment, nil
}

// Repository method for finding a service by name.
func (r *TreatmentRepo) FindByName(ctx context.Context, name string) (*treatmentdb.Treatment, error) {
	treatment, err := r.tq.FindTreatmentByName(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			// treatment not found is not an error here, return nil
			return nil, nil
		}
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return &treatment, nil
}
