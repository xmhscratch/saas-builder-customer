package routers

import (
	// "log"
	"net/http"
	"strconv"

	"gopkg.in/guregu/null.v4"

	"localdomain/customer/models"

	"github.com/gin-gonic/gin"
)

// FindAllAttributes comment
func (ctx *RouteContext) FindAllAttributes() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			attrs            []*models.Attribute
			isVisibleOnFront *bool
			isVisibleInList  *bool
			isConfigurable   *bool
			results          []*map[string]interface{}
		)

		pageNum, err := strconv.Atoi(ginCtx.DefaultQuery("p", "1"))
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		pageSize, err := strconv.Atoi(ginCtx.DefaultQuery("ps", "-1"))
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		codeNameFilter := ginCtx.QueryMap("a")

		if ginCtx.Query("f") != "" {
			v, err := models.ParseBool(ginCtx.Query("f"))
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}
			isVisibleOnFront = &v
		}

		if ginCtx.Query("g") != "" {
			v, err := models.ParseBool(ginCtx.Query("g"))
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}
			isVisibleInList = &v
		}

		if ginCtx.Query("c") != "" {
			v, err := models.ParseBool(ginCtx.Query("c"))
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}
			isConfigurable = &v
		}

		opts := &models.PaginationOptions{
			Size:   pageSize,
			Number: pageNum,
		}

		attrFilter := &models.AttributeFilter{
			CodeNames: codeNameFilter,
		}
		if isVisibleOnFront != nil {
			attrFilter.IsVisibleOnFront = null.BoolFromPtr(isVisibleOnFront)
		}
		if isVisibleInList != nil {
			attrFilter.IsVisibleInList = null.BoolFromPtr(isVisibleInList)
		}
		if isConfigurable != nil {
			attrFilter.IsConfigurable = null.BoolFromPtr(isConfigurable)
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		attrs, err = models.FindAllAttributes(db, attrFilter, opts)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		for _, attr := range attrs {
			info, err := attr.FetchDetailInfo(db)
			if err != nil {
				ctx.Error(err)(ginCtx)
				break
			}
			results = append(results, &info)
		}
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		total := int(models.CountAllAttributes(db, attrFilter))

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   opts.BuildDelta(pageNum, total),
			"results": results,
		})
	}
}

// FindOneAttribute comment
func (ctx *RouteContext) FindOneAttribute() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		codeName := ginCtx.Param("codeName")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		attribute, err := models.FindOneAttribute(db, codeName)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var results map[string]interface{}
		results, err = attribute.FetchDetailInfo(db)
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

// NewAttribute comment
func (ctx *RouteContext) NewAttribute() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err        error
			infoValues map[string]interface{}
			results    map[string]interface{}
		)

		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		infoValues, err = ctx._parseForm(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		attribute := &models.Attribute{
			IsVisibleOnFront: null.BoolFrom(true),
			IsConfigurable:   null.BoolFrom(true),
			// IsUserDefined:    null.BoolFrom(true),
		}
		err = attribute.Normalize(infoValues)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = attribute.Register(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		results, err = attribute.FetchDetailInfo(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = ctx._dispatchEvent(organizationID, "customer-attribute/create", results)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	}
}

// UpdateAttribute comment
func (ctx *RouteContext) UpdateAttribute() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err        error
			infoValues map[string]interface{}
			results    map[string]interface{}
		)

		organizationID := ginCtx.Request.Header.Get("x-organization-id")
		codeName := ginCtx.Param("codeName")

		infoValues, err = ctx._parseForm(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		attribute, err := models.FindOneAttribute(db, codeName)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = attribute.Update(db, infoValues)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		results, err = attribute.FetchDetailInfo(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = ctx._dispatchEvent(organizationID, "customer-attribute/update", results)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	}
}

// DeleteAttribute comment
func (ctx *RouteContext) DeleteAttribute() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var err error

		organizationID := ginCtx.Request.Header.Get("x-organization-id")
		codeName := ginCtx.Param("codeName")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		attribute, err := models.FindOneAttribute(db, codeName)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = attribute.Delete(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = ctx._dispatchEvent(organizationID, "customer-attribute/delete", map[string]interface{}{"codeName": codeName})
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": true,
		})
	}
}
