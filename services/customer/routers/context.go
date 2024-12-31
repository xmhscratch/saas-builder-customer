package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"

	// "bytes"
	"reflect"

	"localdomain/customer/core"
	"localdomain/customer/models"

	"github.com/gin-gonic/gin"
	lru "github.com/hnlq715/golang-lru"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// RouteContext comment
type RouteContext struct {
	Config          *core.Config
	SessionKeyVault *lru.Cache
}

// GetDatabase comment
func (ctx *RouteContext) GetDatabase(ginCtx *gin.Context) (db *gorm.DB, err error) {
	var organizationID string

	organizationID = ginCtx.Param("organizationID")
	if len(organizationID) == 0 {
		organizationID = ginCtx.Request.Header.Get("x-organization-id")
	}
	if len(organizationID) == 0 {
		err = errors.New("organization not found")
		return nil, err
	}

	customerDb := fmt.Sprintf("system_customer_%s", organizationID)
	return core.NewDatabase(ctx.Config, customerDb)
}

// _parseForm comment
func (ctx *RouteContext) _parseForm(ginCtx *gin.Context) (infoValues map[string]interface{}, err error) {
	contentType := ginCtx.Request.Header.Get("content-type")

	// contentLength, err := models.ParseInt(ginCtx.Request.Header.Get("content-length"))
	// if err != nil {
	// 	ctx.Error(err)(ginCtx)
	// 	return
	// }

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		ctx.Error(err)(ginCtx)
		return
	}

	infoValues = make(map[string]interface{})

	switch true {
	case strings.HasPrefix(mediaType, "multipart/"):
		{
			mpR := multipart.NewReader(ginCtx.Request.Body, params["boundary"])
			c := make(chan map[string]interface{})

			go func(c chan map[string]interface{}) {
				// arrKeyNames := map[string]bool{}

				for {
					p, err := mpR.NextPart()
					if err == io.EOF {
						break
					}
					if err != nil {
						return
					}
					defer p.Close()

					fKey := p.FormName()
					fVal, err := io.ReadAll(p)
					if err == io.EOF {
						break
					}
					if err != nil {
						return
					}

					if infoValues[fKey] != nil {
						oldFVal := infoValues[fKey]

						if reflect.TypeOf([]([]uint8){}) == reflect.TypeOf(oldFVal) {
							infoValues[fKey] = append(infoValues[fKey].([]([]uint8)), fVal)
						} else {
							// if oldFVal, err = models.ParseString(oldFVal); err != nil {
							// 	continue
							// }
							infoValues[fKey] = []([]uint8){oldFVal.([]uint8), fVal}
						}
						// if !arrKeyNames[fKey] {
						// 	arrKeyNames[fKey] = true
						// }
					} else {
						infoValues[fKey] = fVal
					}
				}

				// for arrName, _ := range arrKeyNames {
				// 	// var serializeByts []byte
				// 	// if serializeByts, err = json.Marshal(infoValues[arrName]); err != nil {
				// 	// 	continue
				// 	// }
				// 	// log.Println(string(serializeByts[:]))
				// 	infoValues[arrName] = infoValues[arrName]
				// }

				c <- infoValues
			}(c)

			infoValues = <-c
			break
		}

	case strings.HasPrefix(contentType, "application/json"):
		{
			err = ginCtx.ShouldBindJSON(&infoValues)
			break
		}

	default:
		{
			ginCtx.Request.ParseForm()
			for k, v := range ginCtx.Request.PostForm {
				if len(v) == 1 {
					infoValues[k] = v[0]
				} else {
					infoValues[k] = v
				}
			}
			break
		}
	}

	// log.Println(infoValues)
	return infoValues, err
}

// _parseAttributeFilter comment
func (ctx *RouteContext) _parseAttributeFilter(ginCtx *gin.Context) (attrFilter *models.AttributeFilter, err error) {
	var (
		isVisibleOnFront *bool
		isVisibleInList  *bool
	)
	sortBy := ginCtx.QueryMap("s")
	codeNameFilter := ginCtx.QueryMap("a")
	entityTypeFilter := ginCtx.QueryMap("e")

	if ginCtx.Query("f") != "" {
		v, err := models.ParseBool(ginCtx.Query("f"))
		if err != nil {
			return attrFilter, err
		}
		isVisibleOnFront = &v
	}

	if ginCtx.Query("g") != "" {
		v, err := models.ParseBool(ginCtx.Query("g"))
		if err != nil {
			return attrFilter, err
		}
		isVisibleInList = &v
	}

	attrFilter = &models.AttributeFilter{
		EntityTypes: entityTypeFilter,
		CodeNames:   codeNameFilter,
		SortBy:      sortBy,
	}
	if isVisibleOnFront != nil {
		attrFilter.IsVisibleOnFront = null.BoolFromPtr(isVisibleOnFront)
	}
	if isVisibleInList != nil {
		attrFilter.IsVisibleInList = null.BoolFromPtr(isVisibleInList)
	}

	return attrFilter, err
}

// _parsePagination comment
func (ctx *RouteContext) _parsePagination(ginCtx *gin.Context) (opts *models.PaginationOptions, err error) {
	pageNum, err := strconv.Atoi(ginCtx.DefaultQuery("p", "1"))
	if err != nil {
		return nil, err
	}

	pageSize, err := strconv.Atoi(ginCtx.DefaultQuery("ps", "20"))
	if err != nil {
		return nil, err
	}

	opts = &models.PaginationOptions{
		Size:   pageSize,
		Number: pageNum,
	}
	return opts, err
}

func (ctx *RouteContext) _dispatchEvent(organizationID string, eventName string, results map[string]interface{}) (err error) {
	objects, err := json.Marshal(results)
	if err != nil {
		return err
	}
	postValues := &url.Values{}
	postValues.Set("_json", string(objects))

	urlString := core.BuildString("http://", ctx.Config.ClusterHostNames.Integration, "/dispatch/", eventName)
	return core.CreateAPIRequest(ctx.Config, organizationID, urlString, postValues)
}
