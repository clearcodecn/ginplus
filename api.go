package ginplus

type PageData struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}

type ApiResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

func success(data interface{}) *ApiResponse {
	return &ApiResponse{
		Data:    data,
		Success: true,
		Message: "",
		Code:    200,
	}
}

func successPage(data interface{}, count int64) *ApiResponse {
	return &ApiResponse{
		Data: &PageData{
			Data:  data,
			Count: count,
		},
		Success: true,
		Message: "",
		Code:    200,
	}
}

func fail(code int, message string) *ApiResponse {
	return &ApiResponse{
		Success: false,
		Message: message,
		Code:    code,
	}
}

func (c *Context) Success(data interface{}) {
	c.JSON(200, success(data))
}

func (c *Context) Fail(code int, message string) {
	c.JSON(200, fail(code, message))
}

func (c *Context) SuccessPage(data interface{}, count int64) {
	c.JSON(200, successPage(data, count))
}
