package routers

import (
	"net/http"

	"localdomain/customer/models"

	"github.com/gin-gonic/gin"
)

// DeleteCustomerAttributeValue comment
func (ctx *RouteContext) DeleteCustomerAttributeValue() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err error
			// results map[string]interface{}
		)

		customerID := ginCtx.Param("customerID")
		codeName := ginCtx.Param("codeName")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = models.DeleteCustomerAttributeValue(db, customerID, codeName)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{})
	}
}
