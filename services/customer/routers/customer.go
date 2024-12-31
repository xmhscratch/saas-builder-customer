package routers

import (
	"encoding/json"
	"net/http"

	// "strconv"
	"errors"
	// "log"
	"time"

	"localdomain/customer/models"

	dateparse "github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	null "gopkg.in/guregu/null.v4"
)

// FindAllCustomers comment
func (ctx *RouteContext) FindAllCustomers() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err     error
			total   int
			results []*map[string]interface{}
		)

		userID := ginCtx.Request.Header.Get("x-user-id")
		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		opts, err := ctx._parsePagination(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		pageNum := opts.GetPageNumber()
		pageSize := opts.GetLimit()
		pageOffset := opts.GetOffset()

		// attrFilter, err := ctx._parseAttributeFilter(ginCtx)
		// if err != nil {
		// 	ctx.Error(err)(ginCtx)
		// 	return
		// }
		isForce, err := models.ParseBool(ginCtx.DefaultQuery("f", "0"))
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		searcher, err := models.NewCustomerSearch(ctx.Config, organizationID, isForce)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		q, err := models.BuildFilterAttributeTextQuery(db, ginCtx.DefaultQuery("q", ""))
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		aq := ginCtx.DefaultQuery("aq", "")
		fq := ginCtx.Query("fq")

		findResultIDs, total, err := searcher.SearchIDs(q, aq, fq, pageOffset, pageSize)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		for _, customerID := range findResultIDs {
			customer, err := models.FindOneCustomer(db, customerID)
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}
			if customer.ID.IsZero() {
				total--
				continue
			}
			if !customer.TimestampModel.DeletedAt.IsZero() {
				total--
				continue
			}
			// info, err := customer.FetchDetailInfo(db, attrFilter)
			info, err := customer.FetchDetailInfo(db, nil)
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}

			info["userId"] = userID
			info["organizationId"] = organizationID

			results = append(results, &info)
		}

		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   opts.BuildDelta(pageNum, total),
			"results": results,
		})
	}
}

// FindOneCustomer comment
func (ctx *RouteContext) FindOneCustomer() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		customerID := ginCtx.Param("customerID")

		userID := ginCtx.Request.Header.Get("x-user-id")
		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		isForce, err := models.ParseBool(ginCtx.DefaultQuery("f", "0"))
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		_, err = models.NewCustomerSearch(ctx.Config, organizationID, isForce)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customer, err := models.FindOneCustomer(db, customerID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		if customer.ID.IsZero() {
			ctx.Error(errors.New("customer not found"))(ginCtx)
			return
		}

		attrFilter, err := ctx._parseAttributeFilter(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		var results map[string]interface{}
		results, err = customer.FetchDetailInfo(db, attrFilter)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results["userId"] = userID
		results["organizationId"] = organizationID

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   (models.PaginationOptions{}).NoDelta(),
			"results": results,
		})
	}
}

// FindOneCustomerByEmail comment
func (ctx *RouteContext) FindOneCustomerByEmail() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err     error
			results map[string]interface{}
		)

		userID := ginCtx.Request.Header.Get("x-user-id")
		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		emailAddress := ginCtx.PostForm("emailAddress")

		if emailAddress == "" {
			ctx.Error(errors.New("email address is empty"))(ginCtx)
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customer, err := models.FindOneCustomerByEmail(db, emailAddress)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		if customer.ID.IsZero() {
			ctx.Error(errors.New("customer not found"))(ginCtx)
			return
		}

		attrFilter, err := ctx._parseAttributeFilter(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results, err = customer.FetchDetailInfo(db, attrFilter)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results["userId"] = userID
		results["organizationId"] = organizationID

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   (models.PaginationOptions{}).NoDelta(),
			"results": results,
		})
	}
}

// FindOneCustomerByEmailWithShopify comment
func (ctx *RouteContext) FindOneCustomerByEmailWithShopify() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err     error
			results map[string]interface{}
		)

		userID := ginCtx.Request.Header.Get("x-user-id")
		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		emailAddress := ginCtx.PostForm("emailAddress")
		if emailAddress == "" {
			ctx.Error(errors.New("email address is empty"))(ginCtx)
			return
		}

		shopifyCustomerID := ginCtx.PostForm("shopifyCustomerId")
		if shopifyCustomerID == "" {
			ctx.Error(errors.New("shopify customer identifier is empty"))(ginCtx)
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customer, err := models.FindOneCustomerByEmailWithShopify(db, emailAddress, shopifyCustomerID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}
		// log.Println(customer)
		if customer.ID.IsZero() {
			ctx.Error(errors.New("customer not found"))(ginCtx)
			return
		}

		attrFilter, err := ctx._parseAttributeFilter(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results, err = customer.FetchDetailInfo(db, attrFilter)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results["userId"] = userID
		results["organizationId"] = organizationID

		ginCtx.JSON(http.StatusOK, gin.H{
			"delta":   (models.PaginationOptions{}).NoDelta(),
			"results": results,
		})
	}
}

// NewCustomer comment
func (ctx *RouteContext) NewCustomer() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err        error
			results    map[string]interface{}
			infoValues map[string]interface{}
		)

		defer func() {
			if err != nil {
				ctx.Error(err)(ginCtx)
				return
			}
			ginCtx.JSON(http.StatusOK, gin.H{
				"results": results,
			})
		}()

		userID := ginCtx.Request.Header.Get("x-user-id")
		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		infoValues, err = ctx._parseForm(ginCtx)
		if err != nil {
			return
		}

		// jsonString, _ := json.Marshal(infoValues)
		// log.Println(string(jsonString[:]))
		// log.Println(models.ParseString(infoValues["emailAddress"]))
		// log.Println(models.ParseString(infoValues["test1"]))

		attrFilter, err := ctx._parseAttributeFilter(ginCtx)
		if err != nil {
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			return
		}

		customerID, err := models.ParseString(infoValues["id"])
		if err != nil {
			return
		}

		emailAddress, err := models.ParseString(infoValues["emailAddress"])
		if err != nil {
			return
		}

		shopifyCustomerID, err := models.ParseString(infoValues["shopifyCustomerId"])
		if err != nil {
			return
		}

		var createdAt time.Time
		if createdAtValue, err := models.ParseString(infoValues["createdAt"]); err != nil {
			return
		} else {
			if createdAt, err = dateparse.ParseAny(createdAtValue); err != nil {
				createdAtValue = time.Now().Format(time.RFC850)
				if createdAt, err = dateparse.ParseAny(createdAtValue); err != nil {
					return
				}
			}
		}

		inputPassword, err := models.ParseString(infoValues["password"])
		if err != nil {
			return
		}

		customer := &models.Customer{
			DefaultGormModel: models.DefaultGormModel{
				ID: null.StringFrom(customerID),
			},
			EmailAddress:      null.StringFrom(emailAddress),
			ShopifyCustomerID: null.StringFrom(shopifyCustomerID),
			TimestampModel: models.TimestampModel{
				CreatedAt: null.TimeFrom(createdAt),
			},
		}

		var forceCreate bool = false
		isAdmin, _ := ginCtx.Get("isAdmin")
		if isAdmin == false {
			forceCreate, err = models.ParseBool(ginCtx.DefaultQuery("fc", "0"))
			if err != nil {
				return
			}
		} else {
			forceCreate = true
		}

		err = customer.Register(db, inputPassword, forceCreate)
		if err != nil {
			return
		}

		delete(infoValues, "id")
		delete(infoValues, "emailAddress")
		delete(infoValues, "shopifyCustomerId")
		delete(infoValues, "createdAt")
		delete(infoValues, "updatedAt")
		// log.Println(infoValues)
		err = customer.SetAttributes(db, infoValues)
		if err != nil {
			return
		}

		results, err = customer.FetchDetailInfo(db, attrFilter)
		if err != nil {
			return
		}

		_, err = models.NewCustomerSearch(ctx.Config, organizationID, false)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results["userId"] = userID
		results["organizationId"] = organizationID

		err = ctx._dispatchEvent(organizationID, "customer/create", map[string]interface{}{"customerId": customerID})
		if err != nil {
			return
		}
	}
}

// UpdateCustomer comment
func (ctx *RouteContext) UpdateCustomer() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err        error
			results    map[string]interface{}
			infoValues map[string]interface{}
		)

		userID := ginCtx.Request.Header.Get("x-user-id")
		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		infoValues, err = ctx._parseForm(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customerID := ginCtx.Param("customerID")
		// emailAddress := infoValues.Get("emailAddress")

		attrFilter, err := ctx._parseAttributeFilter(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		b, err := json.Marshal(infoValues)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = json.Unmarshal(b, &infoValues)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customer, err := models.FindOneCustomer(db, customerID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = customer.SetAttributes(db, infoValues)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = customer.Update(db, map[string]interface{}{})
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results, err = customer.FetchDetailInfo(db, attrFilter)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		_, err = models.NewCustomerSearch(ctx.Config, organizationID, false)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		results["userId"] = userID
		results["organizationId"] = organizationID

		err = ctx._dispatchEvent(organizationID, "customer/update", map[string]interface{}{"customerId": customerID})
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	}
}

// DeleteCustomer comment
func (ctx *RouteContext) DeleteCustomer() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var err error

		// userID := ginCtx.Request.Header.Get("x-user-id")
		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		customerID := ginCtx.Param("customerID")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		customer, err := models.FindOneCustomer(db, customerID)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = customer.Delete(db)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		_, err = models.NewCustomerSearch(ctx.Config, organizationID, false)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		err = ctx._dispatchEvent(organizationID, "customer/delete", map[string]interface{}{"customerId": customerID})
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"results": true,
		})
	}
}

// ImportCustomer comment
func (ctx *RouteContext) ImportCustomer() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err      error
			importer *models.CustomerImporter
			driveID  string
			fileID   string
		)

		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		infoValues, err := ctx._parseForm(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		importer, err = models.NewCustomerImporter(ctx.Config, organizationID, db)
		defer importer.Dispose()

		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		if driveID, err = models.ParseString(infoValues["driveId"]); err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		if fileID, err = models.ParseString(infoValues["fileId"]); err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		if err = importer.AddFile(driveID, fileID); err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{})
	}
}

// ExportCustomer comment
func (ctx *RouteContext) ExportCustomer() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err      error
			exporter *models.CustomerExporter
		)

		organizationID := ginCtx.Request.Header.Get("x-organization-id")

		db, err := ctx.GetDatabase(ginCtx)
		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		exporter, err = models.NewCustomerExporter(ctx.Config, organizationID, db)
		defer exporter.Dispose()

		if err != nil {
			ctx.Error(err)(ginCtx)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{})
	}
}

// // ImportCustomer comment
// func (ctx *RouteContext) ImportCustomer() gin.HandlerFunc {
// 	return func(ginCtx *gin.Context) {
// 		var (
// 			err        error
// 			results    map[string]interface{}
// 			infoValues map[string]interface{}
// 		)

// 		userID := ginCtx.Request.Header.Get("x-user-id")
// 		organizationID := ginCtx.Request.Header.Get("x-organization-id")

// 		infoValues, err = ctx._parseForm(ginCtx)
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		attrFilter, err := ctx._parseAttributeFilter(ginCtx)
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		db, err := ctx.GetDatabase(ginCtx)
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		emailAddress, err := models.ParseString(infoValues["emailAddress"])
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		customer := &models.Customer{
// 			EmailAddress: null.StringFrom(emailAddress),
// 		}

// 		err = customer.Register(db)
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		delete(infoValues, "emailAddress")

// 		err = customer.SetAttributes(db, infoValues)
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// 		results, err = customer.FetchDetailInfo(db, attrFilter)
// 		if err != nil {
// 			ctx.Error(err)(ginCtx)
// 			return
// 		}

// _, err = models.NewCustomerSearch(ctx.Config, organizationID, false)
// if err != nil {
// 	ctx.Error(err)(ginCtx)
// 	return
// }

// 		results["userId"] = userID
// 		results["organizationId"] = organizationID

// 		ginCtx.JSON(http.StatusOK, gin.H{
// 			"results": results,
// 		})
// 	}
// }
