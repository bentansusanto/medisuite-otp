package repo

import (
	categoryRepo "medisuite-api/app/repo/categories"
	rolePermissionRepo "medisuite-api/app/repo/role_permissions"
	roleRepo "medisuite-api/app/repo/roles"
	treatmentRepo "medisuite-api/app/repo/treatments"
	userRepo "medisuite-api/app/repo/users"
)

type IRepo interface {
	UserRepo() userRepo.IUserRepo
	RoleRepo() roleRepo.IRoleRepo
	RolePermissionRepo() rolePermissionRepo.IRolePermissionRepo
	CategoryRepo() categoryRepo.ICategoryRepo
	TreatmentRepo() treatmentRepo.ITreatmentRepo
}

type Repo struct {
	store Store
}

func NewRepo(store Store) IRepo {
	return &Repo{store: store}
}

func (r *Repo) UserRepo() userRepo.IUserRepo {
	q := r.store.Queries()
	return userRepo.NewUserRepo(q.Users, q.Sessions)
}

func (r *Repo) RoleRepo() roleRepo.IRoleRepo {
	q := r.store.Queries()
	return roleRepo.NewRoleRepo(q.Roles)
}

func (r *Repo) RolePermissionRepo() rolePermissionRepo.IRolePermissionRepo {
	q := r.store.Queries()
	return rolePermissionRepo.NewRolePermissionRepo(q.RolePermissions)
}

func (r *Repo) CategoryRepo() categoryRepo.ICategoryRepo {
	q := r.store.Queries()
	return categoryRepo.NewCategoryRepo(q.Categories)
}

func (r *Repo) TreatmentRepo() treatmentRepo.ITreatmentRepo {
	q := r.store.Queries()
	return treatmentRepo.NewTreatmentRepo(q.Treatments)
}
