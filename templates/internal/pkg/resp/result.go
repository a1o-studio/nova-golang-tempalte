package resp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:model Result
type Result[T any] struct {
	Code    int `json:"code" example:"200"` // 业务状态码
	Message any `json:"message"`            // 消息
	Data    T   `json:"data"`               // 返回结果
} //	@name	Result

type HttpError struct {
	Code    int `json:"code"`    // 错误码
	Message any `json:"message"` // 错误消息
} //	@name	HttpError

type ErrorOption func(*Result[any])

func WithError(err AppError) ErrorOption {
	return func(r *Result[any]) {
		r.Code = err.Code
		r.Message = err.Message
	}
}

func WithMessage(msg any) ErrorOption {
	return func(r *Result[any]) {
		r.Message = msg
	}
}

func WithCode(code int) ErrorOption {
	return func(r *Result[any]) {
		r.Code = code
	}
}

func Success[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Result[T]{
		Code:    0, // code=0 表示成功
		Message: "success",
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, httpStatus, code int, message any) {
	// Abort 后不会执行后续中间件与 handler 的后续操作
	c.Abort()
	c.JSON(httpStatus, Result[any]{
		Code:    code,
		Message: message,
	})
}

// 服务端内部错误 500
func ServerError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrServerError.Code,
		Message: ErrServerError.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusInternalServerError, res.Code, res.Message)
}

// 404 错误
func NotFoundError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrNotFound.Code,
		Message: ErrNotFound.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusNotFound, res.Code, res.Message)
}

func WrapNotFoundError(options ...ErrorOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		NotFoundError(c, options...)
	}
}

// 请求参数错误
func InvalidError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrBadRequest.Code,
		Message: ErrBadRequest.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusBadRequest, res.Code, res.Message)
}

// 超时
func TimeoutError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrGatewayTimeout.Code,
		Message: ErrGatewayTimeout.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusGatewayTimeout, res.Code, res.Message)
}

// 参数校验失败
func FailedValidationError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrStatusUnprocessableEntity.Code,
		Message: ErrStatusUnprocessableEntity.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusUnprocessableEntity, res.Code, res.Message)
}

func ForbiddenError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrForbidden.Code,
		Message: ErrForbidden.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusForbidden, res.Code, res.Message)
}

func UnauthorizedError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrUnauthorized.Code,
		Message: ErrUnauthorized.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusUnauthorized, res.Code, res.Message)
}

func ConflictError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrConflict.Code,
		Message: ErrConflict.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusConflict, res.Code, res.Message)
}

// 不支持的请求方法
func MethodNotAllowedError(c *gin.Context, options ...ErrorOption) {
	message := fmt.Sprintf("The %s method is not supported for this resource", c.Request.Method)
	res := Result[any]{
		Code:    ErrMethodNotAllowed.Code,
		Message: message,
	}
	Error(c, http.StatusMethodNotAllowed, res.Code, res.Message)
}

func TooManyRequestsError(c *gin.Context, options ...ErrorOption) {
	res := Result[any]{
		Code:    ErrTooManyRequests.Code,
		Message: ErrTooManyRequests.Message,
	}
	updateErrorOption(&res, options...)
	Error(c, http.StatusTooManyRequests, res.Code, res.Message)
}

func WrapMethodNotAllowedError(options ...ErrorOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		MethodNotAllowedError(c, options...)
	}
}

func updateErrorOption(res *Result[any], options ...ErrorOption) {
	for _, option := range options {
		option(res)
	}
}

type ErrorMapping struct {
	Target error
	Action func(c *gin.Context, err error)
}

func HandleErrors(c *gin.Context, err error, mappings []ErrorMapping) {
	for _, mapping := range mappings {
		if errors.Is(err, mapping.Target) {
			mapping.Action(c, err)
			return
		}
	}
	// 如果没有匹配到任何错误，返回服务器错误
	ServerError(c, WithMessage("An unexpected error occurred"))
}
