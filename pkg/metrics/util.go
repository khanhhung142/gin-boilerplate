package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func getCidFromReq(g *gin.Context) string {
	cid := g.GetHeader("X-Company-Id")
	if cid != "" {
		return cid
	}

	if cid = g.Request.URL.Query().Get("cid"); cid != "" {
		return cid
	}

	if cid = g.Request.PostFormValue("cid"); cid != "" {
		return cid
	}

	req := g.Request

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	var mapKeyValue = make(map[string]interface{})
	err = json.Unmarshal(body, &mapKeyValue)
	if err == nil && mapKeyValue["cid"] != "" {
		return fmt.Sprint(mapKeyValue["cid"])
	}

	return ""
}

// func getRespCodeFromBody(r *gin.Context) (int, string) {

// 	dataByte := r.
// 	var mapBody ghttp.DefaultHandlerResponse
// 	err := json.Unmarshal(dataByte, &mapBody)
// 	if err == nil {
// 		return mapBody.Code, ""
// 	}

// 	return gcode.CodeUnknown.Code(), gcode.CodeUnknown.Message()
// }
