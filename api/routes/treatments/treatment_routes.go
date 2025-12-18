package treatments

import (
	"medisuite-api/api/handler"
	"medisuite-api/app/repo"
	"medisuite-api/common/middlewares"

	"github.com/gin-gonic/gin"
)

type ITreatmentRoute interface {
	Run()
}

type TreatmentRoute struct {
	h handler.IHandler
	g *gin.RouterGroup
	r repo.IRepo
}

func NewTreatmentRoute(handler handler.IHandler, group *gin.RouterGroup, repo repo.IRepo) *TreatmentRoute {
	return &TreatmentRoute{
		h: handler,
		g: group,
		r: repo,
	}
}

func (r *TreatmentRoute) Run() {
	groups := r.g.Group("/treatments")
	{
		// routes
		groups.GET("/find_all", r.h.TreatmentHandler().FindAllTreatment)
		groups.GET("/:id", r.h.TreatmentHandler().FindByIdTreatment)
		groups.POST("/create", middlewares.AuthMiddleware(), middlewares.RequirePermission(r.r, "treatment", "create"), r.h.TreatmentHandler().CreateTreatment)
		groups.PUT("/update/:id", middlewares.AuthMiddleware(), middlewares.RequirePermission(r.r, "treatment", "update"), r.h.TreatmentHandler().UpdateTreatment)
		groups.DELETE("/delete/:id", middlewares.AuthMiddleware(), middlewares.RequirePermission(r.r, "treatment", "delete"), r.h.TreatmentHandler().DeleteTreatment)
	}
}
