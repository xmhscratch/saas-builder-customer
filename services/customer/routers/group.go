package routers

import (
	"net/http"
	"strconv"

	"gopkg.in/guregu/null.v4"

	"localdomain/customer/core"
	"localdomain/customer/models"

	"github.com/gin-gonic/gin"
)

// FindAllGroups comment
func (ctx *RouteContext) FindAllGroups() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
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

		rootOnly, err := models.ParseBool(ginCtx.DefaultQuery("ro", "true"))
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		codeNameFilter := ginCtx.QueryMap("a")

		opts := &models.PaginationOptions{
			Size:   pageSize,
			Number: pageNum,
		}

		catFilter := &models.GroupFilter{
			CodeNames: codeNameFilter,
			RootOnly:  rootOnly,
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var groups []*models.Group

		groups, err = models.FindAllGroups(db, catFilter, opts)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var results []*map[string]interface{}

		for _, cat := range groups {
			info, err := cat.FetchDetailInfo(db)
			if err != nil {
				break
			}
			results = append(results, &info)
		}
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		total := int(models.CountAllGroups(db, catFilter))

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   opts.BuildDelta(pageNum, total),
			"results": results,
		})
	}
}

// FindOneGroup comment
func (ctx *RouteContext) FindOneGroup() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		groupID := ginCtx.Param("groupID")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		group, err := models.FindOneGroup(db, groupID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var results map[string]interface{}
		results, err = group.FetchDetailInfo(db)
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

// NewGroup comment
func (ctx *RouteContext) NewGroup() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err     error
			results map[string]interface{}
		)

		codeName := ginCtx.DefaultPostForm("codeName", "")
		title := ginCtx.DefaultPostForm("title", "")
		description := ginCtx.DefaultPostForm("description", "")

		group := &models.Group{
			CodeName:    null.StringFrom(codeName),
			Title:       null.StringFrom(title),
			Description: null.StringFrom(description),
		}
		newID := core.GenerateID()
		group.ID = null.StringFrom(newID)

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = group.Register(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results, err = group.FetchDetailInfo(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	}
}

// UpdateGroup comment
func (ctx *RouteContext) UpdateGroup() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err        error
			results    map[string]interface{}
			infoValues map[string]interface{}
		)

		infoValues, err = ctx._parseForm(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		groupID := ginCtx.Param("groupID")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		group, err := models.FindOneGroup(db, groupID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = group.Update(db, infoValues)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results, err = group.FetchDetailInfo(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	}
}

// DeleteGroup comment
func (ctx *RouteContext) DeleteGroup() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err error
		)

		groupID := ginCtx.Param("groupID")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		group, err := models.FindOneGroup(db, groupID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = group.Delete(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": true,
		})
	}
}

// FindAllGroupDescendants comment
func (ctx *RouteContext) FindAllGroupDescendants() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			results []*map[string]interface{}
		)

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		groupID := ginCtx.Param("groupID")
		groups, err := models.FindAllGroupDescendants(db, groupID)

		for _, group := range groups {
			info, err := group.FetchDetailInfo(db)
			if err != nil {
				break
			}
			results = append(results, &info)
		}
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

// FindAllGroupChildren comment
func (ctx *RouteContext) FindAllGroupChildren() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			results []*map[string]interface{}
		)

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		groupID := ginCtx.Param("groupID")
		groups, err := models.FindAllGroupChildren(db, groupID)

		for _, group := range groups {
			info, err := group.FetchDetailInfo(db)
			if err != nil {
				break
			}
			results = append(results, &info)
		}
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

// FindAllGroupPaths comment
func (ctx *RouteContext) FindAllGroupPaths() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			results []*map[string]interface{}
		)

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		groupID := ginCtx.Param("groupID")
		groups, err := models.FindAllGroupPaths(db, groupID)

		for _, group := range groups {
			info, err := group.FetchDetailInfo(db)
			if err != nil {
				break
			}
			results = append(results, &info)
		}
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
