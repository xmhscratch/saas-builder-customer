package routers

import (
	"net/http"

	"gopkg.in/guregu/null.v4"

	"localdomain/customer/models"

	"github.com/gin-gonic/gin"
)

// FindAllEnumerations comment
func (ctx *RouteContext) FindAllEnumerations() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		keyword := ginCtx.Query("s")

		enumFilter := &models.EnumerationFilter{
			SeachKeyword: null.StringFrom(keyword),
		}

		opts := &models.PaginationOptions{
			Size:   1000,
			Number: 1,
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var enums []*models.Enumeration

		enums, err = models.FindAllEnumerations(db, enumFilter)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var results []*map[string]interface{}

		for _, enum := range enums {
			info, err := enum.FetchDetailInfo(db)
			if err != nil {
				break
			}
			results = append(results, &info)
		}
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		total := int(models.CountAllEnumerations(db, enumFilter))

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   opts.BuildDelta(1, total),
			"results": results,
		})
	}
}

// FindOneEnumeration comment
func (ctx *RouteContext) FindOneEnumeration() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		enumID := ginCtx.Param("enumID")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		enum, err := models.FindOneEnumeration(db, enumID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var results map[string]interface{}
		results, err = enum.FetchDetailInfo(db)
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

// // NewAttribute comment
// func (ctx *RouteContext) NewAttribute() gin.HandlerFunc {
// 	return func(ginCtx *gin.Context) {
// 		var err error

// 		entityType := ginCtx.DefaultPostForm("entityType", "varchar")
// 		codeName := ginCtx.DefaultPostForm("codeName", "")
// 		label := ginCtx.DefaultPostForm("label", "")
// 		description := ginCtx.DefaultPostForm("description", "")
// 		defaultValue := ginCtx.DefaultPostForm("defaultValue", "")

// 		var isUserDefined bool
// 		isUserDefined, err = models.ParseBool(ginCtx.DefaultPostForm("isUserDefined", "true"))
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		var isRequired bool
// 		isRequired, err = models.ParseBool(ginCtx.DefaultPostForm("isRequired", "false"))
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		db, err := ctx.GetDatabase(ginCtx)
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		attribute := &models.Attribute{
// 			EntityType:    null.StringFrom(entityType),
// 			CodeName:      null.StringFrom(codeName),
// 			Label:         null.StringFrom(label),
// 			Description:   null.StringFrom(description),
// 			DefaultValue:  null.StringFrom(defaultValue),
// 			IsUserDefined: null.BoolFrom(isUserDefined),
// 			IsRequired:    null.BoolFrom(isRequired),
// 		}
// 		newID := core.GenerateID()
// 		attribute.ID = null.StringFrom(newID)

// 		err = attribute.Register(db)
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		ginCtx.JSON(http.StatusOK, attribute)
// 	}
// }
