package http

import (
	"catalog-backend/config"
	"catalog-backend/internal/delivery/http/handler"
	"catalog-backend/internal/delivery/http/middleware"
	"catalog-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	productHandler     *handler.ProductHandler
	sectionHandler     *handler.SectionHandler
	contactHandler     *handler.ContactHandler
	uploadHandler      *handler.UploadHandler
	storeHandler       *handler.StoreHandler
	StoreImageHandler  *handler.StoreImageHandler
	customerHandler    *handler.CustomerHandler
	reservationHandler *handler.ReservationHandler
	adminHandler       *handler.AdminHandler
	adminUC            *usecase.AdminUseCase
}

func NewRouter(
	productUC *usecase.ProductUseCase,
	sectionUC *usecase.SectionUseCase,
	contactUC *usecase.ContactUseCase,
	uploadCfg *config.UploadConfig,
	storeUC *usecase.StoreUseCase,
	customerUC *usecase.CustomerUseCase,
	reservationUC *usecase.ReservationUseCase,
	adminUC *usecase.AdminUseCase,
) *Router {
	return &Router{
		productHandler:     handler.NewProductHandler(productUC),
		sectionHandler:     handler.NewSectionHandler(sectionUC),
		contactHandler:     handler.NewContactHandler(contactUC),
		uploadHandler:      handler.NewUploadHandler(productUC, uploadCfg.Path, uploadCfg.MaxSize),
		storeHandler:       handler.NewStoreHandler(storeUC),
		StoreImageHandler:  handler.NewStoreImageHandler(storeUC, uploadCfg.Path, uploadCfg.MaxSize),
		customerHandler:    handler.NewCustomerHandler(customerUC),
		reservationHandler: handler.NewReservationHandler(reservationUC),
		adminHandler:       handler.NewAdminHandler(adminUC),
		adminUC:            adminUC,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine, config *config.Config) {
	api := engine.Group("/api/v1")

	// Аутентификация (публично)
	auth := api.Group("/auth")
	{
		auth.GET("/is-first", r.adminHandler.IsFirstAdmin)
		auth.POST("/register", r.adminHandler.Register)
		auth.POST("/login", r.adminHandler.Login)
	}

	// Публичные маршруты (для клиентов)
	public := api.Group("/public")
	{
		public.GET("/products", r.productHandler.GetProducts)
		public.GET("/products/:id", r.productHandler.GetProduct)
		public.GET("/sections", r.sectionHandler.GetSections)
		public.GET("/sections/:id", r.sectionHandler.GetSection)
		public.GET("/contacts", r.contactHandler.GetContacts)
		public.GET("/store-info", r.storeHandler.GetStoreInfo)
		public.POST("/customers", r.customerHandler.GetOrCreate)
		public.POST("/reservations", r.reservationHandler.CreateReservation)
	}

	// Административные маршруты (требуют JWT)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(r.adminUC))
	{
		admin.POST("/products", r.productHandler.CreateProduct)
		admin.PUT("/products/:id", r.productHandler.UpdateProduct)
		admin.DELETE("/products/:id", r.productHandler.DeleteProduct)

		admin.POST("/sections", r.sectionHandler.CreateSection)
		admin.PUT("/sections/:id", r.sectionHandler.UpdateSection)
		admin.DELETE("/sections/:id", r.sectionHandler.DeleteSection)

		admin.POST("/contacts", r.contactHandler.CreateContact)
		admin.PUT("/contacts/:id", r.contactHandler.UpdateContact)
		admin.DELETE("/contacts/:id", r.contactHandler.DeleteContact)

		admin.POST("/products/:id/image", r.uploadHandler.UploadProductImage)

		admin.PUT("/store-info", r.storeHandler.UpdateStoreInfo)
		admin.POST("/store/image", r.StoreImageHandler.UploadStoreImage)

		admin.GET("/reservations", r.reservationHandler.GetReservations)
		admin.PUT("/reservations/:id", r.reservationHandler.UpdateReservationStatus)

		// Управление администраторами
		admin.GET("/admins", r.adminHandler.GetAdmins)
		admin.POST("/admins", r.adminHandler.CreateAdmin)
		admin.PUT("/admins/:id", r.adminHandler.UpdateAdmin)
		admin.DELETE("/admins/:id", r.adminHandler.DeleteAdmin)
	}

	engine.Static("/uploads", config.Upload.Path)

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}
