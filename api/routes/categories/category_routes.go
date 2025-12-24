package categories

import (
	"medisuite-api/api/handler"
	"medisuite-api/app/repo"
	"medisuite-api/common/middlewares"

	"github.com/gin-gonic/gin"
)

type ICategoryRoute interface {
	Run()
}

type CategoryRoute struct {
	h handler.IHandler
	g *gin.RouterGroup
	r repo.IRepo
}

func NewCategoryRoute(handler handler.IHandler, group *gin.RouterGroup, repo repo.IRepo) *CategoryRoute {
	return &CategoryRoute{
		h: handler,
		g: group,
		r: repo,
	}
}

func (r *CategoryRoute) Run() {
	groups := r.g.Group("categories")
	{
		// routes
		groups.GET("/find_all", r.h.CategoryHandler().FindAllCategory)
		groups.GET("/:id", r.h.CategoryHandler().FindByIdCategory)
		groups.POST("/create", middlewares.AuthMiddleware(), middlewares.RequirePermission(r.r, "category", "create"), r.h.CategoryHandler().CreateCategory)
		groups.PUT("/update/:id", middlewares.AuthMiddleware(), middlewares.RequirePermission(r.r, "category", "update"), r.h.CategoryHandler().UpdateCategory)
		groups.DELETE("/delete/:id", middlewares.AuthMiddleware(), middlewares.RequirePermission(r.r, "category", "delete"), r.h.CategoryHandler().DeleteCategory)
	}
}
