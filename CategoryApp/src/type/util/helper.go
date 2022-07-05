package util

import (
	categoryType "CategoryApp/src/type"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func GetClient() http.Client {
	return http.Client{
		Transport: &http.Transport{
			Proxy:                  nil,
			DialContext:            nil,
			Dial:                   nil,
			DialTLSContext:         nil,
			DialTLS:                nil,
			TLSClientConfig:        nil,
			TLSHandshakeTimeout:    0,
			DisableKeepAlives:      false,
			DisableCompression:     false,
			MaxIdleConns:           0,
			MaxIdleConnsPerHost:    0,
			MaxConnsPerHost:        0,
			IdleConnTimeout:        0,
			ResponseHeaderTimeout:  0,
			ExpectContinueTimeout:  0,
			TLSNextProto:           nil,
			ProxyConnectHeader:     nil,
			GetProxyConnectHeader:  nil,
			MaxResponseHeaderBytes: 0,
			WriteBufferSize:        0,
			ReadBufferSize:         0,
			ForceAttemptHTTP2:      false,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
}

func CheckCategoryModel(category categoryType.Category) (bool, *Error) {
	return true, nil
}

func IsValidUUID(uuidStr string) bool {
	if _, err := uuid.Parse(uuidStr); err != nil {
		return false
	}
	return true
}

func ValidatePaginationFilters(page, pageSize string, maxLimit int) (int64, int64) {
	max := int64(maxLimit)

	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil || pageInt < 0 {
		pageInt = 0
	}

	pageSizeInt, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil || pageSizeInt <= 0 {
		pageSizeInt = max
	}
	return pageInt, pageSizeInt
}

func ValidateSortingFilters(entity any, sortingArea, SortingDirection string) (string, int) {
	sort := ""
	var direction int

	if strings.EqualFold(SortingDirection, "asc") || strings.EqualFold(SortingDirection, "1") {
		direction = 1
	} else if strings.EqualFold(SortingDirection, "desc") || strings.EqualFold(SortingDirection, "dsc") || strings.EqualFold(SortingDirection, "-1") {
		direction = -1
	} else {
		direction = 0
	}

	v := reflect.ValueOf(entity)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if sortingArea == typeOfS.Field(i).Name {
			sort = strings.ToLower(typeOfS.Field(i).Name)
			break
		}
	}

	return sort, direction
}

func CreateEqualFilter(value, field string) interface{} {
	return bson.M{"$regex": primitive.Regex{
		Pattern: value,
		Options: "i",
	}}
}

var operators []string = []string{"<=", ">=", "==", "||"}

func CreateFilter(model any, filters string) map[string]interface{} {
	responseFilter := map[string]interface{}{}
	v := reflect.ValueOf(model)
	typeOfS := v.Type()

	filterArr := strings.Split(filters, ",")
	for _, filter := range filterArr {
		for i := 0; i <= len(operators); i++ {
			if strings.Contains(filter, operators[i]) {
				splittedCurrentFilter := strings.Split(filter, operators[i])
				field := splittedCurrentFilter[0]
				value := splittedCurrentFilter[len(splittedCurrentFilter)-1]
				for j := 0; j < v.NumField(); j++ {
					if field == typeOfS.Field(j).Name {
						if strings.Contains(filter, ">=") {
							responseFilter[field] = bson.M{"$gte": value}
							break
						} else if strings.Contains(filter, "<=") {
							if prevFilter, exist := responseFilter[field]; !exist {
								responseFilter[field] = bson.M{"$lte": value}
							} else {
								fmt.Println(prevFilter)
								responseFilter["$and"] = bson.A{
									bson.M{field: value},
									prevFilter,
								}
							}
							break
						} else if strings.Contains(filter, "==") {
							responseFilter[field] = bson.M{"$regex": primitive.Regex{
								Pattern: value,
								Options: "i",
							}}
							break
						} else if strings.Contains(filter, "||") {
							break
						}
					}
				}
			}
		}
	}
	return nil
}

func ValidateForUserHandlerId(id string) *Error {
	if id == "" {
		return PathVariableNotFound.ModifyApplicationName("category handler").ModifyErrorCode(4050)
	}

	if !IsValidUUID(id) {
		return PathVariableIsNotValid.ModifyApplicationName("category handler").ModifyErrorCode(4051)
	}
	return nil
}
