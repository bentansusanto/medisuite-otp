package treatments

import (
	"time"

	"github.com/google/uuid"
)

type TreatmentDTO struct {
	CategoryID    uuid.UUID `json:"category_id" validation:"required"`
	NameTreatment string    `json:"name_treatment" validation:"required"`
	Description   string    `json:"description" validation:"required"`
	Thumbnail     string    `json:"thumbnail" validation:"required"`
	Price         float64   `json:"price" validation:"required"`
	Duration      int32     `json:"duration" validation:"required"`
	IsActive      bool      `json:"is_active" validation:"required"`
}

type TreatmentResponse struct {
	ID            uuid.UUID `json:"id"`
	CategoryID    uuid.UUID `json:"category_id"`
	NameTreatment string    `json:"name_treatment"`
	Description   string    `json:"description"`
	Thumbnail     string    `json:"thumbnail"`
	Price         float64   `json:"price"`
	Duration      int32     `json:"duration"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type CategoryDTO struct {
	NameCategory string `json:"name_category" validation:"required"`
}

type CategoryResponse struct {
	ID           uuid.UUID `json:"id"`
	NameCategory string    `json:"name_category"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
