package http

import (
	"catalog-backend/config"
	"catalog-backend/internal/delivery/http/handler"
	"catalog-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	productHandler    *handler.ProductHandler
	sectionHandler    *handler.SectionHandler
	contactHandler    *handler.ContactHandler
	uploadHandler     *handler.UploadHandler
	storeHandler      *handler.StoreHandler
	StoreImageHandler *handler.StoreImageHandler
}

func NewRouter(
	productUC *usecase.ProductUseCase,
	sectionUC *usecase.SectionUseCase,
	contactUC *usecase.ContactUseCase,
	uploadCfg *config.UploadConfig,
	storeUC *usecase.StoreUseCase,
) *Router {
	return &Router{
		productHandler:    handler.NewProductHandler(productUC),
		sectionHandler:    handler.NewSectionHandler(sectionUC),
		contactHandler:    handler.NewContactHandler(contactUC),
		uploadHandler:     handler.NewUploadHandler(productUC, uploadCfg.Path, uploadCfg.MaxSize),
		storeHandler:      handler.NewStoreHandler(storeUC),
		StoreImageHandler: handler.NewStoreImageHandler(storeUC, uploadCfg.Path, uploadCfg.MaxSize),
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine, config *config.Config) {
	// Группа API
	api := engine.Group("/api/v1")

	// Публичные маршруты (для клиентов)
	public := api.Group("/public")
	{
		public.GET("/products", r.productHandler.GetProducts)
		public.GET("/products/:id", r.productHandler.GetProduct)
		public.GET("/sections", r.sectionHandler.GetSections)
		public.GET("/sections/:id", r.sectionHandler.GetSection)
		public.GET("/contacts", r.contactHandler.GetContacts)
		public.GET("/store-info", r.storeHandler.GetStoreInfo)
	}

	// Административные маршруты (для управления)
	admin := api.Group("/admin")
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
	}

	engine.Static("/uploads", config.Upload.Path)

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}
