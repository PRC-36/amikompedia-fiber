package util

type BaseResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func ConstructBaseResponse(response BaseResponse) (BaseResponse, int) {

	resp := BaseResponse{
		Code:   response.Code,
		Status: response.Status,
		Data:   response.Data,
	}

	return resp, response.Code
}
