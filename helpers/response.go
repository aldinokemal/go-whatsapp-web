package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Meta    interface{} `json:"meta"`
	Results interface{} `json:"results"`
}

func RespondJSON(w *gin.Context, status int, payload interface{}, message ...string) {
	fmt.Println("status ", status)
	var res ResponseData

	res.Code = status
	if len(message) > 0 {
		res.Message = message[0]
	}
	//res.Meta = utils.ResponseMessage(status)
	res.Results = payload

	w.SecureJSON(status, res)
}
