package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lzy1014964035/go-tool-set/service"
	"github.com/lzy1014964035/go-tool-set/service/log"
)

type Connect = websocket.Conn

// 路径组
var pathMap = make(map[string]func(*websocket.Conn, service.Any, service.Any))

// 设置路径
func SetPath(dealPath string, callable func(*websocket.Conn, service.Any, service.Any)){
	service.Dump(dealPath, callable);
	pathMap[dealPath] = callable
}

// 挂起服务
func MakeWSService(servicePort string) {
	service.Dump("挂起服务 端口：" + servicePort);
	http.HandleFunc("/ws", makeConnect)
	http.ListenAndServe(":"+servicePort, nil)
}

// 连接时的闭包处理
func makeConnect(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WebSocket 连接
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 允许所有来源的连接
			return true
		},
	}
	// 升级并返回连接
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.MakeLogError("WebSocket 升级失败:", err);
		return
	}

	// 处理请求
	dealRequest(connect, w, r);
}

// 处理请求
func dealRequest(connect *websocket.Conn, w http.ResponseWriter, r *http.Request){
	for {
		// 读取客户端发送的消息
		_, message, err := connect.ReadMessage()
		if err != nil {
			log.MakeLogError("读取消息失败:", err, message)
			continue;
		}
		messageString := string(message); // 转为字符串
		messageJsonData := service.JsonDecode(messageString); // 转为JSON
		if(messageJsonData == nil){
			dealRequestString(connect, w, r, messageString)
		}else{
			dealRequestJson(connect, w, r, messageJsonData)
		}
	}
}

// 处理请求字符串
func dealRequestString(connect *websocket.Conn, w http.ResponseWriter, r *http.Request, messageString string){
	service.Dump("接收string类型消息", messageString)
	SendToCli(connect, "不处理字符串类型数据，请传入json字符串", "501", "/showMessage", service.ToMap{"requestData": messageString});
}

// 处理请求JSON
func dealRequestJson(connect *websocket.Conn, w http.ResponseWriter, r *http.Request, messageJsonData service.ToMap){
	dealPath := messageJsonData["deal_path"].(string);
	data := messageJsonData["data"].(service.Any);
	// 如果没有传路径
	if(dealPath == ""){
		SendFailToCli(connect, "缺少处理路径 deal_path", "504", nil)
		return;
	}
	// 如果回调不存在
	if(pathMap[dealPath] == nil){
		SendFailToCli(connect, "路径" + dealPath + "在服务端并不存在", "503", nil)
		return;
	}
	// 取闭包
	pathFunction := pathMap[dealPath]
	pathFunction(connect, data, messageJsonData);
}

// 发送信息到
func SendToCli(connect *websocket.Conn, dealPath string, message string, code string, data service.Any) {
	reponseData := service.ToMap{
		"code": code,
		"deal_path": dealPath,
		"message": message,
		"data": data,
	};
	err := connect.WriteJSON(reponseData);
	if err != nil {
		log.MakeLogError("推送到websocket cli失败:", service.ToMap{"reponseData": reponseData, "err": err})
	}
}

// 发送失败消息
func SendFailToCli(connect *websocket.Conn, message string, code string, data service.Any) {
	SendToCli(connect, "/showMessage", message, code, data);
}

// 发送成功消息
func SendSuccessToCli(connect *websocket.Conn, dealPath string, message string, code string, data service.Any){
	SendToCli(connect, dealPath, message, code, data);
}