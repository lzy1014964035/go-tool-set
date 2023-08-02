package websocket

import (
	"net/http"
	"time"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lzy1014964035/go-tool-set/service"
	"github.com/lzy1014964035/go-tool-set/service/log"
)

var wg sync.WaitGroup

// 连接对象类型
type connectObject struct {
	id string
	connect *websocket.Conn
	lastCommunicationTime int64
}

//连接池，以连接ID作为下标
var ConnectionPool = make(map[string]connectObject)

// 生成唯一的连接 ID
func generateID() string {
	timestamp := time.Now().UnixNano()
	return strconv.FormatInt(timestamp, 10)
}

// 挂起服务
func MakeService(servicePort string) {
	service.Dump("挂起服务 端口：" + servicePort);
	http.HandleFunc("/ws", handleWebSocket)
	http.ListenAndServe(":"+servicePort, nil)
}

// 连接时的闭包处理
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
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

	// 将连接加入到连接池中
	connectId := generateID();
	ConnectionPool[connectId] = connectObject{
		id: connectId,
		connect: connect,
		lastCommunicationTime: service.GetTimeSeconds(),
	}
	// 携程启动监听连接对象
	go listenConnectReuqest(connectId)
}


// 监听连接的请求
func listenConnectReuqest(connectId string){
	SendToCli(connectId, SendToCliTypeMap["heart"].(string), "/heart", "成功连接-初始化", service.ToMap{
		"text": "first",
	});
	connectData := ConnectionPool[connectId]
	for {
		// 读取客户端发送的消息
		_, message, err := connectData.connect.ReadMessage()
		if err != nil {
			log.MakeLogError("读取消息失败:", err)
			break
		}
		go getCliRequest(connectId, string(message));
	}
}



// 收到cli的推送
func getCliRequest(connectId string, message string){
	// 更新最后的通讯时间
	connectData, ok := ConnectionPool[connectId]
	if(!ok){
		return;
	}
	connectData.lastCommunicationTime = service.GetTimeSeconds();
	// 解析推送的内容
	messageJsonMap := service.JsonDecode(message).(service.ToMap);
	// 检查收到推送的类型和通讯ID
	messageType,_ := messageJsonMap["type"].(string);
	messageId,_ := messageJsonMap["message_id"].(string);
	if(SendToCliTypeMap[messageType] == nil){
		log.MakeLogError("接收信息信息失败，设置的类型异常", service.ToMap{
			"连接的ID": connectId,
			"消息ID": messageId,
			"信息类型": messageType,
			"信息数据结构": messageJsonMap,
		});
		return;
	}
	// 如果类型是心跳，那么在上面一步已经更新了通讯时间了，所以直接返回处理即可
	if(messageType == SendToCliTypeMap["heart"]){
		return;
	}

	// 附加系统的额外内容
	messageJsonMap["system_other_param"] = service.ToMap{
		"time": service.MakeDateWithYMDHIS,
		"ws_connnect_id": connectId,
		"request_id": generateID(),
	}
	// 推送给RPC服务器
}

// 推送给客户端的类型
var SendToCliTypeMap = service.ToMap{
	"heart": "heart", // 心跳
	"message": "message", // 通讯
}

// 发送消息给
func SendToCli[Any any](connectId string, sendType string, sendPath string, sendMessage string, sendData Any){
	// 更新最后的通讯时间
	connectData, ok := ConnectionPool[connectId]
	if !ok {
		return;
	}
	if(SendToCliTypeMap[sendType] == nil){
		log.MakeLogError("推送信息失败，设置的类型异常", service.ToMap{
			"连接的ID": connectId,
			"推送的类型": sendType,
			"发送的信息": sendMessage,
			"发送的数据结构": sendData,
		});
		return
	}
	requestId := generateID();
	sendMessageToCliData := service.ToMap{
		"message_id": requestId,
		"type": sendType,
		"path": sendPath,
		"message": sendMessage,
		"data": sendData,
	};
	// sendMessageToCliDataJsonString := service.JsonEncode(sendMessageToCliData);
	// err := connectData.connect.WriteMessage(websocket.TextMessage, []byte(sendMessageToCliDataJsonString));
	err := connectData.connect.WriteJSON(sendMessageToCliData);
	if err != nil {
		log.MakeLogError("推送到websocket cli失败:", err)
	}
}




