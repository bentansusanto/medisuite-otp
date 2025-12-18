package routes

// import (
// 	"medisuite-api/api/handler"
// 	"medisuite-api/app/repo"
// 	"medisuite-api/common/middlewares"

// 	"github.com/gin-gonic/gin"
// )

// // Example of how to use RBAC middleware in your routes

// func SetupCategoryRoutes(router *gin.RouterGroup, h handler.IHandler, repo repo.IRepo) {
// 	categories := router.Group("/categories")
// 	{
// 		// Public endpoint - no authentication required
// 		categories.GET("", h.CategoryHandler().FindAllCategory)
// 		categories.GET("/:id", h.CategoryHandler().FindByIdCategory)

// 		// Protected endpoints - require authentication and specific permissions
// 		categories.POST("",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "category", "create"),
// 			h.CategoryHandler().CreateCategory,
// 		)

// 		categories.PUT("/:id",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "category", "update"),
// 			h.CategoryHandler().UpdateCategory,
// 		)

// 		categories.DELETE("/:id",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "category", "delete"),
// 			h.CategoryHandler().DeleteCategory,
// 		)
// 	}
// }

// func SetupTreatmentRoutes(router *gin.RouterGroup, h handler.IHandler, repo repo.IRepo) {
// 	treatments := router.Group("/treatments")
// 	{
// 		// Public endpoints
// 		treatments.GET("", h.TreatmentHandler().FindAllTreatment)
// 		treatments.GET("/:id", h.TreatmentHandler().FindByIdTreatment)

// 		// Protected endpoints - require authentication and specific permissions
// 		treatments.POST("",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "treatment", "create"),
// 			h.TreatmentHandler().CreateTreatment,
// 		)

// 		treatments.PUT("/:id",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "treatment", "update"),
// 			h.TreatmentHandler().UpdateTreatment,
// 		)

// 		treatments.DELETE("/:id",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "treatment", "delete"),
// 			h.TreatmentHandler().DeleteTreatment,
// 		)
// 	}
// }

// func SetupUserRoutes(router *gin.RouterGroup, h handler.IHandler, repo repo.IRepo) {
// 	users := router.Group("/users")
// 	{
// 		// Protected endpoints - require specific roles
// 		users.GET("",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequireRole(repo, "owner", "admin"),
// 			h.UserHandler().GetAllUsers,
// 		)

// 		users.POST("",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "user", "create"),
// 			h.UserHandler().CreateUser,
// 		)

// 		users.PUT("/:id",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "user", "update"),
// 			h.UserHandler().UpdateUser,
// 		)

// 		users.DELETE("/:id",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "user", "delete"),
// 			h.UserHandler().DeleteUser,
// 		)
// 	}
// }

// func SetupPatientRoutes(router *gin.RouterGroup, h handler.IHandler, repo repo.IRepo) {
// 	patients := router.Group("/patients")
// 	{
// 		// All patient endpoints require authentication
// 		patients.Use(middlewares.AuthMiddleware())

// 		// View patients - doctors, admin, owner can view
// 		patients.GET("",
// 			middlewares.RequirePermission(repo, "patient", "read"),
// 			// h.PatientHandler().FindAllPatients,
// 		)

// 		patients.GET("/:id",
// 			middlewares.RequirePermission(repo, "patient", "read"),
// 			// h.PatientHandler().FindPatientById,
// 		)

// 		// Create patient - doctors, admin, owner can create
// 		patients.POST("",
// 			middlewares.RequirePermission(repo, "patient", "create"),
// 			// h.PatientHandler().CreatePatient,
// 		)

// 		// Update patient - doctors, admin, owner can update
// 		patients.PUT("/:id",
// 			middlewares.RequirePermission(repo, "patient", "update"),
// 			// h.PatientHandler().UpdatePatient,
// 		)

// 		// Delete patient - only admin and owner
// 		patients.DELETE("/:id",
// 			middlewares.RequirePermission(repo, "patient", "delete"),
// 			// h.PatientHandler().DeletePatient,
// 		)
// 	}
// }

// func SetupAdminRoutes(router *gin.RouterGroup, h handler.IHandler, repo repo.IRepo) {
// 	admin := router.Group("/admin")
// 	{
// 		// All admin routes require owner or admin role
// 		admin.Use(middlewares.AuthMiddleware())
// 		admin.Use(middlewares.RequireRole(repo, "owner", "admin"))

// 		// Admin dashboard
// 		admin.GET("/dashboard", func(c *gin.Context) {
// 			c.JSON(200, gin.H{"message": "Admin dashboard"})
// 		})

// 		// System settings
// 		admin.GET("/settings", func(c *gin.Context) {
// 			c.JSON(200, gin.H{"message": "System settings"})
// 		})
// 	}
// }

// // Example: Combining multiple middleware
// func SetupAppointmentRoutes(router *gin.RouterGroup, h handler.IHandler, repo repo.IRepo) {
// 	appointments := router.Group("/appointments")
// 	{
// 		// Patients can create and view their own appointments
// 		appointments.POST("",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "appointment", "create"),
// 			// h.AppointmentHandler().CreateAppointment,
// 		)

// 		// Doctors and admin can view all appointments
// 		appointments.GET("",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "appointment", "read"),
// 			// h.AppointmentHandler().FindAllAppointments,
// 		)

// 		// Update appointment
// 		appointments.PUT("/:id",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "appointment", "update"),
// 			// h.AppointmentHandler().UpdateAppointment,
// 		)

// 		// Delete appointment
// 		appointments.DELETE("/:id",
// 			middlewares.AuthMiddleware(),
// 			middlewares.RequirePermission(repo, "appointment", "delete"),
// 			// h.AppointmentHandler().DeleteAppointment,
// 		)
// 	}
// }
