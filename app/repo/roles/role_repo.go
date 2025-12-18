package roles

import (
	"context"
	"log/slog"

	errWrap "medisuite-api/common/errors"
	errConsts "medisuite-api/constants/errors"
	roledb "medisuite-api/pkg/db/roles"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type IRoleRepo interface {
	FindRoleById(ctx context.Context, id uuid.UUID) (*roledb.Role, error)
}

type RoleRepo struct {
	rq *roledb.Queries
}

func NewRoleRepo(rq *roledb.Queries) IRoleRepo {
	return &RoleRepo{rq: rq}
}

func (r *RoleRepo) FindRoleById(ctx context.Context, roleID uuid.UUID) (*roledb.Role, error) {
	slog.Info("FindRoleById called", "role_id", roleID)

	row, err := r.rq.GetRoleByID(ctx, roleID)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Warn("Role not found in DB", "role_id", roleID)
			return nil, nil
		}
		slog.Error("GetRoleByID error", "role_id", roleID, "err", err)
		return nil, errWrap.WrapError(errConsts.ErrRoleNotFound)
	}

	role := &roledb.Role{
		ID:              row.ID,
		Name:            row.Name,
		Code:            row.Code,
		Level:           row.Level,
		Description:     row.Description,
		CanSelfRegister: row.CanSelfRegister,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}

	slog.Info("Role found", "role_id", role.ID, "code", role.Code, "can_self_register", role.CanSelfRegister)
	return role, nil
}
