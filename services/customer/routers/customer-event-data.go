package routers

import (
	"net/http"
	// "log"
	// "encoding/json"

	// "localdomain/customer/core"
	"localdomain/customer/models"

	"github.com/gin-gonic/gin"
)

// FindAllCustomerEventData comment
func (ctx *RouteContext) FindAllCustomerEventData() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			eventSeries []*models.CustomerEventDataInfo
			results     []*map[string]interface{}
		)

		customerID := ginCtx.Param("customerID")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		eventSeries, err = models.FindAllCustomerEventData(db, customerID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		for _, evtData := range eventSeries {
			info, err := models.FetchDetailInfo(db, evtData)
			if err != nil {
				break
			}
			results = append(results, &info)
		}
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		// jsonString, _ := json.Marshal(results)
		// log.Println(string(jsonString[:]))

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   (models.PaginationOptions{}).NoDelta(),
			"results": results,
		})
	}
}
