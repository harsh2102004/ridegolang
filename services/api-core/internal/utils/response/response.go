package response

import (
	"org/api-core/internal/utils/constants"

	"github.com/gin-gonic/gin"
)

// 🔹 Standard API response structure
type APIResponse struct {
	StatusCode int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Code       string      `json:"code,omitempty"`
}

// 🔹 Internal function (single source of truth)
func send(c *gin.Context, status int, message string, data interface{}, err interface{}, code string) {
	resp := APIResponse{
		StatusCode: status,
		Message:    message,
	}

	if data != nil {
		resp.Data = data
	}

	if err != nil {
		resp.Error = err
	}

	if code != "" {
		resp.Code = code
	}

	c.JSON(status, resp)
}


func OK(c *gin.Context, message string, data interface{}) {
	send(c, constants.SUCCESS, message, data, nil, "")
}

func Created(c *gin.Context, message string, data interface{}) {
	send(c, constants.CREATED, message, data, nil, "")
}

func NoContent(c *gin.Context) {
	c.Status(constants.NOCONTENT)
}


func BadRequest(c *gin.Context, message string, err interface{}) {
	send(c, constants.BADREQUEST, message, nil, err, "")
}

func Unauthorized(c *gin.Context, message string) {
	send(c, constants.UNAUTHORIZED, message, nil, nil, "")
}

func Forbidden(c *gin.Context, message string) {
	send(c, constants.FORBIDDEN, message, nil, nil, "")
}

func NotFound(c *gin.Context, message string) {
	send(c, constants.NOTFOUND, message, nil, nil, "")
}


func InternalServerError(c *gin.Context, message string, err interface{}) {
	send(c, constants.INTERNALSERVERERROR, message, nil, err, "")
}