package routers

import (
	"fmt"
	"net/http"
	"time"

	"localdomain/customer/core"

	"github.com/gin-gonic/gin"
	lru "github.com/hnlq715/golang-lru"
)

// NewContext comment
func NewContext(cfg *core.Config) (ctx *RouteContext, err error) {
	var sessionKeyVault *lru.Cache
	if sessionKeyVault, err = lru.NewWithExpire(5000, 15*60*time.Second); err != nil {
		return nil, err
	}
	ctx = &RouteContext{
		Config:          cfg,
		SessionKeyVault: sessionKeyVault,
	}
	return ctx, err
}

// Init comment
func (ctx *RouteContext) Init(router *gin.Engine) {
	// c := cors.New(cors.Options{
	// 	// AllowedOrigins:   []string{"*"},
	// 	// AllowCredentials: true,
	// 	// Enable Debugging for testing, consider disabling in production
	// 	Debug: true,
	// })
	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(func(ginCtx *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			ginCtx.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		ginCtx.AbortWithStatus(http.StatusInternalServerError)
	}))
	router.Use(ctx.CheckAdminRequest())

	// router.Use(ctx.CORS())
	router.OPTIONS("/*urlPath", ctx.Ping())

	// customer search
	router.GET("/customers", ctx.FindAllCustomers())
	router.GET("/customer/:customerID/groups", ctx.FindAllCustomerGroups())
	// customer operation
	router.POST("/customers", ctx.NewCustomer())
	router.PUT("/customers", ctx.NewCustomer())
	router.POST("/customers/import", ctx.ImportCustomer())
	router.PUT("/customers/import", ctx.ImportCustomer())
	router.POST("/customers/export", ctx.ExportCustomer())
	router.PUT("/customers/export", ctx.ExportCustomer())
	router.PATCH("/customers/export", ctx.ExportCustomer())
	router.GET("/customer/:customerID", ctx.FindOneCustomer())
	router.POST("/customer/:customerID", ctx.UpdateCustomer())
	router.PUT("/customer/:customerID", ctx.UpdateCustomer())
	router.PATCH("/customer/:customerID", ctx.UpdateCustomer())
	router.DELETE("/customer/:customerID", ctx.DeleteCustomer())
	router.DELETE("/customer/:customerID/attribute/:codeName", ctx.DeleteCustomerAttributeValue())
	// customer info
	router.POST("/customer-info/by-email", ctx.FindOneCustomerByEmail())
	router.POST("/customer-info/by-email-with-shopify", ctx.FindOneCustomerByEmailWithShopify())

	// customer attribute search
	router.GET("/attributes", ctx.FindAllAttributes())
	router.GET("/attribute/:codeName", ctx.FindOneAttribute())
	// customer attribute operation
	router.POST("/attributes", ctx.NewAttribute())
	router.PUT("/attributes", ctx.NewAttribute())
	router.POST("/attribute/:codeName", ctx.UpdateAttribute())
	router.PATCH("/attribute/:codeName", ctx.UpdateAttribute())
	router.DELETE("/attribute/:codeName", ctx.DeleteAttribute())

	// group search
	router.GET("/groups", ctx.FindAllGroups())
	router.GET("/group/:groupID", ctx.FindOneGroup())
	router.GET("/group/:groupID/children", ctx.FindAllGroupChildren())
	router.GET("/group/:groupID/descendants", ctx.FindAllGroupDescendants())
	router.GET("/group/:groupID/paths", ctx.FindAllGroupPaths())
	// group operation
	router.POST("/groups", ctx.NewGroup())
	router.PUT("/groups", ctx.NewGroup())
	router.POST("/group/:groupID", ctx.UpdateGroup())
	router.PATCH("/group/:groupID", ctx.UpdateGroup())
	router.DELETE("/group/:groupID", ctx.DeleteGroup())

	// customer series event
	router.GET("/customer-events", ctx.FindAllCustomerEvents())
	router.GET("/customer-event/:eventName", ctx.FindOneCustomerEvent())
	router.POST("/customer-events", ctx.NewCustomerEvent())
	router.PUT("/customer-events", ctx.NewCustomerEvent())
	router.PUT("/customer-event/:eventName", ctx.UpdateCustomerEvent())
	router.PATCH("/customer-event/:eventName", ctx.UpdateCustomerEvent())
	router.DELETE("/customer-event/:eventName", ctx.DeleteCustomerEvent())
	router.GET("/customer-event-data/:customerID", ctx.FindAllCustomerEventData())

	router.POST("/acknowledgement/earned-achievements", ctx.UpdateEarnedAchievementInfo())
	router.GET("/acknowledgement/earned-achievements/:customerID", ctx.FindAllCustomerEarnedAchievements())
	router.POST("/acknowledgement/earned-achievements/:customerID", ctx.AckCustomerEarnedAchievements())
	router.GET("/acknowledgement/earned-achievement-fulfillments/:customerID", ctx.FindAllCustomerEarnedAchievementFulfillments())
	router.POST("/acknowledgement/earned-achievement-fulfillments/:customerID", ctx.AckCustomerEarnedAchievementFulfillments())

	// group category
	router.GET("/group-categories", ctx.FindAllGroupCategories())
	router.GET("/group-category/:groupCategoryID", ctx.FindOneGroupCategory())

	// customer's group
	router.GET("/customer-groups/:customerID", ctx.FindAllCustomerGroups())

	// attribute's enumeration
	router.GET("/enumerations", ctx.FindAllEnumerations())
	router.GET("/enumeration/:enumID", ctx.FindOneEnumeration())

	// gift request
	router.GET("/gift-requests", ctx.FindAllGiftRequests())
	router.POST("/gift-requests/:giftID", ctx.NewGiftRequest())
	router.PUT("/gift-requests/:giftID", ctx.NewGiftRequest())
	router.DELETE("/gift-requests/:giftID", ctx.DeleteGiftRequestsByGiftID())
	router.GET("/gift-request/:requestID", ctx.FindOneGiftRequest())
	router.POST("/gift-request/:requestID", ctx.UpdateGiftRequest())
	router.PATCH("/gift-request/:requestID", ctx.UpdateGiftRequest())
	router.DELETE("/gift-request/:requestID", ctx.DeleteGiftRequest())

	// complete gift request
	router.POST("/gift-request-complete/:requestID", ctx.CompleteGiftRequest())

	// point balance history
	router.GET("/balances/:customerID", ctx.FindAllPointBalanceHistories())
	router.POST("/balances/:customerID", ctx.NewPointBalanceHistory())
	router.PUT("/balances/:customerID", ctx.NewPointBalanceHistory())

	// statistics reporting
	router.GET("/stats/customer/sign-up", ctx.ReportRecentMonthSignup())
	router.GET("/stats/customer/detail/:customerID", ctx.ReportCustomerDetailStatistics())

	// internal operation
	router.POST("/_setup/:organizationID", ctx.SetupOrganization())
	router.POST("/_sync/:organizationID", ctx.SyncOrganizationCustomers())

	// router.POST("/_ping", ctx.Noop())
	// router.POST("/_series/:organizationID/:customerID", ctx.SyncOrganizationCustomerEventDataSeries())

	// default
	router.NoRoute(ctx.NoRoute())
}
