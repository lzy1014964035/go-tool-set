package http

import (
	"encoding/json"
	"net/http"
	"github.com/lzy1014964035/go-tool-set/service"
)

type ResponseWriter = http.ResponseWriter
type Request        = http.Request
type ReuqestData    = service.ToMap
type OtherData      = service.Any

// 添加路由
func AddUrl(url string, callableFunction func(http.ResponseWriter, *http.Request, ReuqestData, OtherData), otherData OtherData){
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request){
		var data ReuqestData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}
		// 获取 JSON 中的字段值
		callableFunction(w, r, data, otherData);
	})
}

// 创建HTTP服务器
func MakeHttpService(listenPort string){
	service.Fail(http.ListenAndServe(":" + listenPort, nil))
}

// 创建HTTPS服务器
func MakeHttpsService(listenPort string, sslPemPath string, sslKeyPath string){
	service.Fail(http.ListenAndServeTLS(":"+listenPort, sslPemPath, sslKeyPath, nil))
}

// 返回内容
func ReturnSuccess(w http.ResponseWriter, message string, data service.ToMap){
	ReturnData(w, "success", "200", message, data);
}

// 返回失败
func ReturnFail(w http.ResponseWriter, message string, data service.ToMap){
	ReturnData(w, "fail", "405", message, data);
}

// 返回结构
func ReturnData(w http.ResponseWriter, status string, code string, message string, data service.ToMap){
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(service.JsonEncode(service.ToMap{
		"status": "fail",
		"code": code,
		"result": service.ToMap{
			"msg": message,
			"data": data,
		},
	})));
}