package http

import (
	"net/http"
	"github.com/lzy1014964035/go-tool-set/service"
)

// 添加路由
func AddUrl(url string, callableFunction func(http.ResponseWriter, *http.Request)){
	http.HandleFunc(url, func(ResponseWriter http.ResponseWriter, Request *http.Request){
		callableFunction(ResponseWriter, Request);
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