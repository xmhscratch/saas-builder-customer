package routers

import (
	"context"
	"encoding/json"

	// "log"
	"net/http"

	"localdomain/customer/core"
	"localdomain/customer/models"

	"github.com/gin-gonic/gin"
	redisV8 "github.com/go-redis/redis/v8"
)

type SessionInfo struct {
	Cookie         interface{} `json:"cookie"`
	HomeURL        string      `json:"homeURL"`
	User           string      `json:"user"`
	OrganizationID string      `json:"organizationId"`
}

// SetupOrganization comment
func (ctx *RouteContext) SetupOrganization() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		organizationID := ginCtx.Param("organizationID")

		newDbName := core.BuildString("system_customer", "_", organizationID)
		_, err := core.NewDatabase(ctx.Config, newDbName)
		if err != nil {
			ginCtx.JSON(http.StatusOK, gin.H{})
			return
		}

		_, err = models.NewCustomerSearch(ctx.Config, organizationID, true)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{})
	}
}

// SyncOrganizationCustomers comment
func (ctx *RouteContext) SyncOrganizationCustomers() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		// organizationID := ginCtx.Param("organizationID")

		// err := ctx._syncSearchImport(ginCtx, organizationID)
		// if err != nil {
		// 	ctx.Error(err)(ginCtx)
		// 	return
		// }

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customerEvents, err := models.FindAllCustomerEvents(db, &models.CustomerEventFilter{}, &models.PaginationOptions{})
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		for _, customerEvent := range customerEvents {
			eventData := &models.CustomerEventData{
				EventName: customerEvent.EventName,
			}
			err = eventData.Sync(db, ctx.Config)
			if err != nil {
				break
			}
		}
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{})
	}
}

// Noop comment
func (ctx *RouteContext) Noop() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.JSON(http.StatusOK, gin.H{})
	}
}

// AdminBypassRequest comment
func (ctx *RouteContext) CheckAdminRequest() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		organizationID := ginCtx.Request.Header.Get("x-organization-id")
		sessionKey := ginCtx.Request.Header.Get("x-ssid")

		isAdmin, _ := ctx.isAdminRequest(sessionKey, organizationID)
		ginCtx.Set("isAdmin", isAdmin)
		ginCtx.Next()
	}
}

func (ctx *RouteContext) isAdminRequest(sessionKey string, organizationID string) (bool, error) {
	var (
		err       error
		activeOrg string
	)
	redisSessionClient := redisV8.NewClient(&redisV8.Options{
		Addr:         core.BuildString(ctx.Config.Redis.Host, ":", ctx.Config.Redis.Port),
		Password:     ctx.Config.Redis.Password,
		DB:           1,
		MinIdleConns: 0,
		MaxConnAge:   0,
	})
	background := context.Background()

	if activeOrg, ok := ctx.SessionKeyVault.Get(sessionKey); !ok {
		var sessionInfo *SessionInfo
		sessionVal, err := redisSessionClient.Get(background, core.BuildString("session_", sessionKey)).Result()
		if err != nil {
			return false, err
		}
		if err = json.Unmarshal([]byte(sessionVal), &sessionInfo); err != nil {
			return false, err
		}
		activeOrg = sessionInfo.OrganizationID
		if activeOrg == organizationID {
			ctx.SessionKeyVault.Add(sessionKey, activeOrg)
			return true, err
		}
		return false, err
	} else {
		if activeOrg != organizationID {
			return false, err
		} else {
			return true, err
		}
	}

	if activeOrg == organizationID {
		return true, err
	}
	return false, err
}
