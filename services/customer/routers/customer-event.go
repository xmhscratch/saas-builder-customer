package routers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	null "gopkg.in/guregu/null.v4"

	// "localdomain/customer/core"
	"localdomain/customer/models"
)

// FindAllCustomerEvents comment
func (ctx *RouteContext) FindAllCustomerEvents() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			customerEvents []*models.CustomerEvent
			results        []*map[string]interface{}
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

		customerEventFilter := &models.CustomerEventFilter{}

		opts := &models.PaginationOptions{
			Size:   pageSize,
			Number: pageNum,
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customerEvents, err = models.FindAllCustomerEvents(db, customerEventFilter, opts)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		for _, evt := range customerEvents {
			info, err := evt.FetchDetailInfo(db)
			if err != nil {
				break
			}
			results = append(results, &info)
		}
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		total := int(models.CountAllCustomerEvents(db, customerEventFilter))

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   opts.BuildDelta(pageNum, total),
			"results": results,
		})
	}
}

// FindOneCustomerEvent comment
func (ctx *RouteContext) FindOneCustomerEvent() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		eventName := ginCtx.Param("eventName")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customerEvent, err := models.FindOneCustomerEvent(db, eventName)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var results map[string]interface{}
		results, err = customerEvent.FetchDetailInfo(db)
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

// NewCustomerEvent comment
func (ctx *RouteContext) NewCustomerEvent() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err     error
			results map[string]interface{}
		)

		infoValues, err := ctx._parseForm(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var (
			eventName               string    = ""
			title                   string    = ""
			synchronizationInterval int64     = int64(5000)
			createdAt               time.Time = time.Now()
		)

		eventName, err = models.ParseString(infoValues["eventName"])
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		if eventName == "" {
			ctx.Error(errors.New("EventName is required attribute"))(ginCtx)
			return
		}

		if infoValues["title"] != "" {
			title, err = models.ParseString(infoValues["title"])
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}
		}

		if (infoValues["synchronizationInterval"] != "") && (infoValues["synchronizationInterval"] != nil) {
			synchronizationInterval = int64(float64(infoValues["synchronizationInterval"].(float64)))
			if synchronizationInterval < int64(3600000) {
				synchronizationInterval = int64(3600000)
			}
		}

		customerEvent := &models.CustomerEvent{
			EventName:               null.StringFrom(eventName),
			Title:                   null.StringFrom(title),
			SynchronizationInterval: null.IntFrom(synchronizationInterval),
			TimestampModel: models.TimestampModel{
				CreatedAt: null.TimeFrom(createdAt),
			},
		}

		if infoValues["description"] != "" {
			description, err := models.ParseString(infoValues["description"])
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}
			customerEvent.Description = null.StringFrom(description)
		}

		if infoValues["aggregateMethod"] != "" {
			aggregateMethod, err := models.ParseString(infoValues["aggregateMethod"])
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}
			customerEvent.AggregateMethod = null.StringFrom(aggregateMethod)
		} else {
			customerEvent.AggregateMethod = null.StringFromPtr(nil)
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = customerEvent.Register(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results, err = customerEvent.FetchDetailInfo(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"context": gin.H{},
			"results": results,
		})
	}
}

// UpdateCustomerEvent comment
func (ctx *RouteContext) UpdateCustomerEvent() gin.HandlerFunc {
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

		eventName := ginCtx.Param("eventName")
		if eventName == "" {
			ctx.Error(errors.New("EventName is required attribute"))(ginCtx)
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customerEvent, err := models.FindOneCustomerEvent(db, eventName)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = customerEvent.Update(db, infoValues)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results, err = customerEvent.FetchDetailInfo(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	}
}

// DeleteCustomerEvent comment
func (ctx *RouteContext) DeleteCustomerEvent() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var err error

		eventName := ginCtx.Param("eventName")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customerEvent, err := models.FindOneCustomerEvent(db, eventName)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = customerEvent.Delete(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": true,
		})
	}
}
