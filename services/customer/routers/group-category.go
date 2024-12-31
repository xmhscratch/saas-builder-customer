package routers

import (
	"net/http"
	// "strconv"

	// "gopkg.in/guregu/null.v4"

	// "localdomain/customer/core"
	"localdomain/customer/models"

	"github.com/gin-gonic/gin"
)

// FindAllGroupCategories comment
func (ctx *RouteContext) FindAllGroupCategories() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			results         []*map[string]interface{}
			groupCategories []*models.GroupCategory
		)

		groupCategoryFilter := &models.GroupCategoryFilter{}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		groupCategories, err = models.FindAllGroupCategories(db, groupCategoryFilter)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		for _, groupCategory := range groupCategories {
			info, err := groupCategory.FetchDetailInfo(db)
			if err != nil {
				break
			}
			results = append(results, &info)
		}
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		// total := models.CountAllGroupCategories(db, groupCategoryFilter)

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   (models.PaginationOptions{}).NoDelta(),
			"results": results,
		})
	}
}

// FindOneGroupCategory comment
func (ctx *RouteContext) FindOneGroupCategory() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			results map[string]interface{}
		)

		groupCategoryID := ginCtx.Param("groupCategoryID")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		groupCategory, err := models.FindOneGroupCategory(db, groupCategoryID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results, err = groupCategory.FetchDetailInfo(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   (models.PaginationOptions{}).NoDelta(),
			"results": results,
		})
	}
}
