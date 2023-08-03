package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lzy1014964035/go-tool-set/service"
	"github.com/lzy1014964035/go-tool-set/service/log"
)

type ConnectData struct {
	Id  string // 连接ID
	Connect    *websocket.Conn // 连接对象
	LastReuqestTime int64  // 最后的请求时间 时间戳
}

type Connect = websocket.Conn
type RequestData = service.Any

var connectPool = make(map[string]ConnectData);

// 路径组
var pathMap = make(map[string]func(ConnectData, service.Any))

// 设置路径
func SetPath(dealPath string, callable func(ConnectData, service.Any)){
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
	// 加入连接池
	connectData := addConnectInPool(connect);
	// 处理请求
	dealRequest(connectData, w, r);
}

// 将连接添加进连接池
func addConnectInPool(connect *websocket.Conn) ConnectData {
	var connectId string = string(service.GetTimeNanoseconds());
	connectPool[connectId] = ConnectData{
		Id: connectId,
		Connect: connect,
		LastReuqestTime: service.GetTimeSeconds(),
	};
	return connectPool[connectId]
}

// 删除池中的连接
func deleteConnectFromPool(connectId string){
	delete(connectPool, connectId)
}

// 处理请求
func dealRequest(connectData ConnectData, w http.ResponseWriter, r *http.Request){
	for {
		// 读取客户端发送的消息
		_, message, err := connectData.Connect.ReadMessage()
		connectData.LastReuqestTime = service.GetTimeSeconds()
		if err != nil {
			log.MakeLogError("读取消息失败:", err, message)
			continue;
		}
		messageString := string(message); // 转为字符串
		messageJsonData := service.JsonDecode(messageString); // 转为JSON
		if(messageJsonData == nil){
			dealRequestString(connectData, w, r, messageString)
		}else{
			dealRequestJson(connectData, w, r, messageJsonData)
		}
	}
}

// 处理请求字符串
func dealRequestString(connectData ConnectData, w http.ResponseWriter, r *http.Request, messageString string){
	service.Dump("接收string类型消息", messageString)
	SendToCli(connectData, "不处理字符串类型数据，请传入json字符串", "501", "/showMessage", service.ToMap{"requestData": messageString});
}

// 处理请求JSON
func dealRequestJson(connectData ConnectData, w http.ResponseWriter, r *http.Request, messageJsonData service.ToMap){
	dealPath := messageJsonData["deal_path"].(string);
	// data := messageJsonData["data"].(service.Any);
	// 如果没有传路径
	if(dealPath == ""){
		SendFailToCli(connectData, "缺少处理路径 deal_path", "504", nil)
		return;
	}
	// 如果回调不存在
	if(pathMap[dealPath] == nil){
		SendFailToCli(connectData, "路径" + dealPath + "在服务端并不存在", "503", nil)
		return;
	}
	// 取闭包
	pathFunction := pathMap[dealPath]
	pathFunction(connectData, messageJsonData);
}

// 发送信息到
func SendToCli(connectData ConnectData, dealPath string, message string, code string, data service.Any) {
	reponseData := service.ToMap{
		"code": code,
		"deal_path": dealPath,
		"message": message,
		"data": data,
	};
	err := connectData.Connect.WriteJSON(reponseData);
	if err != nil {
		log.MakeLogError("推送到websocket cli失败:", service.ToMap{"reponseData": reponseData, "err": err})
	}
}

// 发送失败消息
func SendFailToCli(connectData ConnectData, message string, code string, data service.Any) {
	SendToCli(connectData, "/showMessage", message, code, data);
}

// 发送成功消息
func SendSuccessToCli(connectData ConnectData, dealPath string, message string, data service.Any){
	SendToCli(connectData, dealPath, message, "200", data);
}