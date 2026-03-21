package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"soporte/internal/core/domain"
	"soporte/internal/delivery/http/middlewares"
)

func bindJSON[T any](c *gin.Context) (T, bool) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		fail(c, domain.ValidationError(formatValidationError(err), err))
		return request, false
	}
	return request, true
}

func bindQuery[T any](c *gin.Context) (T, bool) {
	var query T
	if err := c.ShouldBindQuery(&query); err != nil {
		fail(c, domain.ValidationError(formatValidationError(err), err))
		return query, false
	}
	return query, true
}

func getID(c *gin.Context) (int, bool) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		fail(c, err)
		return 0, false
	}
	return id, true
}

func parseID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil {
		return 0, domain.ValidationError("id must be a valid integer", err)
	}

	if id <= 0 {
		return 0, domain.ValidationError("id must be greater than 0", nil)
	}

	return id, nil
}

func json(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"data":       data,
		"request_id": middlewares.GetRequestID(c),
	})
}

func created(c *gin.Context, location string, data any) {
	c.Header("Location", location)
	c.JSON(http.StatusCreated, gin.H{
		"data":       data,
		"request_id": middlewares.GetRequestID(c),
	})
}

func list(c *gin.Context, data any, total int64, limit, offset int) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"meta": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
		"request_id": middlewares.GetRequestID(c),
	})
}

func fail(c *gin.Context, err error) {
	c.JSON(domain.StatusCode(err), gin.H{
		"error": gin.H{
			"code":    domain.ErrorCode(err),
			"message": domain.ErrorMessage(err),
		},
		"request_id": middlewares.GetRequestID(c),
	})
}
