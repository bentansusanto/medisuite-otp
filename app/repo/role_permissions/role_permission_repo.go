package role_permissions

import (
	"context"
	"log/slog"

	errWrap "medisuite-api/common/errors"
	errConsts "medisuite-api/constants/errors"
	rolepermissiondb "medisuite-api/pkg/db/role_permissions"

	"github.com/google/uuid"
)

type IRolePermissionRepo interface {
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]rolepermissiondb.GetRolePermissionsRow, error)
}

type RolePermissionRepo struct {
	rpq *rolepermissiondb.Queries
}

func NewRolePermissionRepo(rpq *rolepermissiondb.Queries) IRolePermissionRepo {
	return &RolePermissionRepo{rpq: rpq}
}

// Repository method for getting all permissions for a specific role
func (r *RolePermissionRepo) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]rolepermissiondb.GetRolePermissionsRow, error) {
	permissions, err := r.rpq.GetRolePermissions(ctx, roleID)
	if err != nil {
		slog.Error("Error getting role permissions", "error", err, "role_id", roleID)
		return nil, errWrap.WrapError(errConsts.ErrSQLError)
	}
	return permissions, nil
}
