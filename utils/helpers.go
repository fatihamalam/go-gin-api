package utils

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var GinCtx *gin.Context

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func Paginate(c *gin.Context) (limit, page, offset int) {
	limit = 10
	page = 1

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	offset = (page - 1) * limit
	return
}

func Sorting(c *gin.Context, defaultOrder, defaultSort string, allowedFields map[string]bool) (order string) {
	orderBy := c.DefaultQuery("order_by", defaultOrder)
	sort := c.DefaultQuery("sort", defaultSort)
	if sort != "asc" && sort != "desc" {
		orderBy = defaultSort
	}

	if !allowedFields[orderBy] {
		orderBy = defaultOrder
		sort = defaultSort
	}

	order = orderBy + " " + sort
	return
}

func BuildRelation(c *gin.Context, relations []string) []string {
	var includes []string
	queryIncludes := c.Query("includes")

	if queryIncludes == "" {
		return includes
	}

	for _, include := range strings.Split(queryIncludes, ",") {
		if slices.Contains(relations, include) {
			includes = append(includes, include)
		}
	}

	return includes
}

func ValidationErrorToText(err error) map[string]string {
	var ve validator.ValidationErrors
	out := make(map[string]string)

	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			tag := fe.Tag()

			switch tag {
			case "required":
				out[field] = fmt.Sprintf("The %s field is required", field)
			case "email":
				out[field] = "Invalid email format"
			case "uniqueEmail":
				out[field] = "Email address already exists"
			case "min":
				out[field] = fmt.Sprintf("The min length of %s field is %s character", field, fe.Param())
			case "max":
				out[field] = fmt.Sprintf("The max length of %s field is %s character", field, fe.Param())
			default:
				out[field] = fe.Error()
			}
		}
	}

	return out
}