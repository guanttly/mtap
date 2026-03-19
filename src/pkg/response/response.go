// Package response 提供统一的HTTP响应格式封装
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// R 统一响应结构
type R struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResult 分页结果
type PageResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewPageResult 创建分页结果
func NewPageResult[T any](items []T, total int64, page, pageSize int) *PageResult[T] {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &PageResult[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// OK 成功响应（无数据）
func OK(c *gin.Context) {
	c.JSON(http.StatusOK, R{
		Code:    0,
		Message: "success",
	})
}

// OKWithData 成功响应（带数据）
func OKWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, R{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created 创建成功
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, R{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Fail 失败响应（通用）
func Fail(c *gin.Context, httpStatus int, code bizErr.Code, detail string) {
	c.JSON(httpStatus, R{
		Code:    int(code),
		Message: bizErr.MessageOf(code),
		Data:    detail,
	})
}

// FailWithError 根据 BizError 自动响应
func FailWithError(c *gin.Context, err error) {
	if be, ok := err.(*bizErr.BizError); ok {
		httpStatus := codeToHTTPStatus(be.Code)
		c.JSON(httpStatus, R{
			Code:    int(be.Code),
			Message: be.Message,
			Data:    be.Detail,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, R{
		Code:    int(bizErr.ErrInternal),
		Message: bizErr.MessageOf(bizErr.ErrInternal),
		Data:    err.Error(),
	})
}

// codeToHTTPStatus 错误码到HTTP状态码映射
func codeToHTTPStatus(code bizErr.Code) int {
	switch code {
	case bizErr.ErrUnauthorized:
		return http.StatusUnauthorized
	case bizErr.ErrForbidden:
		return http.StatusForbidden
	case bizErr.ErrNotFound:
		return http.StatusNotFound
	case bizErr.ErrInvalidParam:
		return http.StatusBadRequest
	case bizErr.ErrRateLimit:
		return http.StatusTooManyRequests
	case bizErr.ErrDuplicate, bizErr.ErrConflict:
		return http.StatusConflict
	case bizErr.ErrTimeout:
		return http.StatusGatewayTimeout
	default:
		if code >= 2000 && code < 8000 {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}
}
