package util

type BaseResponse struct {
	Code   int                 `json:"code"`
	Status string              `json:"status"`
	Data   interface{}         `json:"data"`
	Error  []ApiErrorValidator `json:"errors"`
}

func ConstructBaseResponse(response BaseResponse) (BaseResponse, int) {

	resp := BaseResponse{
		Code:   response.Code,
		Status: response.Status,
		Data:   response.Data,
		Error:  response.Error,
	}

	return resp, response.Code
}
