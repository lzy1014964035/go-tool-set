package log

import (
	"os"
	"github.com/lzy1014964035/go-tool-set/service"
)

var logBasePath = "./logs/"
type LogService struct {}


func MakeLog (logFileName string, logData ...service.Any) {
	// 组合日志路径
	logPath := logBasePath + service.MakeDateWithYMD() + "/" + logFileName;
	// 检查和创建日志文件
	service.CheckAndCreateFile(logPath);
	// 打开日志文件
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		service.Dump("打开日志文件异常", logPath, err);
	}
	// 结束时关闭
	defer file.Close()
	// 写入日志
	dateTime := service.MakeDateWithYMDHIS();
	logData = append([]service.Any{"---------" + dateTime + "---------"}, logData...) // 追加装饰用头部
	logData = append(logData, "-------------------------------------") // 追加装饰用尾部
	// 轮询写入日志
	service.ForeachArray(logData, func(k, v interface{}) {
		logDataStr := service.JsonEncode(v);
		_, err = file.WriteString(logDataStr + "\r\n")
		if err != nil {
			service.Dump("文件写入失败", logPath, err);
		}
	})
}

// 创建错误的日志文件
func MakeLogError(logData ...service.Any) {
	MakeLog("error.log", logData...);
}



