package routes

import (
	"medisuite-api/api/handler"
	categoryRoutes "medisuite-api/api/routes/categories"
	treatmentRoutes "medisuite-api/api/routes/treatments"
	userRoutes "medisuite-api/api/routes/users"
	"medisuite-api/app/repo"

	"github.com/gin-gonic/gin"
)

type IRoutes interface {
	Serve()
	UserRoutes() userRoutes.IUserRoutes
	CategoryRoutes() categoryRoutes.ICategoryRoute
	TreatmentRoutes() treatmentRoutes.ITreatmentRoute
}

type Routes struct {
	h handler.IHandler
	g *gin.RouterGroup
	r repo.IRepo
}

func NewRoutes(handler handler.IHandler, group *gin.RouterGroup, repo repo.IRepo) IRoutes {
	return &Routes{
		h: handler,
		g: group,
		r: repo,
	}
}

func (r *Routes) Serve() {
	r.UserRoutes().Run()
	r.CategoryRoutes().Run()
	r.TreatmentRoutes().Run()
}

func (r *Routes) UserRoutes() userRoutes.IUserRoutes {
	return userRoutes.NewUserRoutes(r.h, r.g)
}

func (r *Routes) CategoryRoutes() categoryRoutes.ICategoryRoute {
	return categoryRoutes.NewCategoryRoute(r.h, r.g, r.r)
}

func (r *Routes) TreatmentRoutes() treatmentRoutes.ITreatmentRoute {
	return treatmentRoutes.NewTreatmentRoute(r.h, r.g, r.r)
}
