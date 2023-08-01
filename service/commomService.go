package service

import (
	"fmt"
	"reflect"
	"encoding/json"
	"os"
	"time"
	"strings"
	"io/ioutil"
	"math"
	"github.com/ghodss/yaml"
)


type Any any

type Object = interface{}
type Array = []interface{}

// map别名
type ToMap = map[string]interface{}
type TM = map[string]interface{}

// array别名
type ToArray = []interface{}
type TA = []interface{}


// 打印全部的内容
func Dump (data ...interface{}) {
	switch v := reflect.ValueOf(data);
	v.Kind() {
	default:
		// 打印内容
		ForeachArray(data, func (k,v interface{}){
			fmt.Println(v);
		});
	}
}

// 遍历数组
// array 遍历的数组
// callable（k 下标, v 值） 回调函数
func ForeachArray[Any any] (array []Any, callable func(interface{}, interface{})) {
	for i := 0; i < len(array); i++ {
		// printing the days of the week
		callable(i, array[i])
	}
}


// 数据结构转JSON
func JsonEncode[Any any] (data Any) (string) {
	jsonBytes,err := json.Marshal(data)
	if(err != nil){
		Dump("jsonEncode异常", err);
	}
	jsonString := string(jsonBytes)
	return jsonString;

	// prettyJSON, err := json.MarshalIndent(data, "", "  ")
	// if(err != nil){
	// 	self.Dump("jsonDecode异常", err);
	// }
	// return prettyJSON;
}

// json字符串转结构
func JsonDecode(jsonString string) (Object) {
	var data Object
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		Dump("jsonDecode异常", err)
	}
	return data;
}

// 检查和创建文件夹
func CheckAndCreateDir(dirPath string) {
	// 创建日志文件
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			// 创建目录失败
			panic(err)
		}
	}
}

// 简化和创建文件
func CheckAndCreateFile(filePath string) {
	// 路由解析，剔除最后的节点，获取文件夹路径
	filePathArray := Explode("/", filePath);
	dirPathArray := ArrayRemoveEnd(filePathArray);
	// 如果解析出文件夹路由，则对文件夹路径检查和创建
	if(dirPathArray != nil){ 
		dirPathString := Implode("/", dirPathArray);
		CheckAndCreateDir(dirPathString);
	}
	// 创建日志文件
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 文件不存在，创建文件
		file, err := os.Create(filePath)
		if err != nil {
			Dump("创建文件失败:", filePath, err)
			return
		}
		defer file.Close()
		Dump("文件创建成功", filePath)
	}
}

// 获取时分秒日期格式时间
func MakeDateWithYMDHIS(timestmap ...int64) (string) {
	var timestamp int64
	// 如果 timestmap 参数提供了值，则使用提供的值
	if len(timestmap) > 0 {
		timestamp = timestmap[0]
	} else {
		// 否则，默认为当前时间的时间戳
		timestamp = time.Now().Unix()
	}
	currentTime := time.Unix(timestamp, 0)
	currentDate := currentTime.Format("2006-01-02 15:04:05")
	return currentDate;
}

// 获取日期格式时间
func MakeDateWithYMD(timestmap ...int64) (string) {
	var timestamp int64
	// 如果 timestmap 参数提供了值，则使用提供的值
	if len(timestmap) > 0 {
		timestamp = timestmap[0]
	} else {
		// 否则，默认为当前时间的时间戳
		timestamp = time.Now().Unix()
	}
	currentTime := time.Unix(timestamp, 0)
	currentDate := currentTime.Format("2006-01-02")
	return currentDate;
}

// 获取时间戳
func GetTimeSeconds() (int64) {
	currentTime := time.Now()
	timestampSeconds := currentTime.Unix()
	return timestampSeconds
}

// 获取纳秒时间戳
func GetTimeNanoseconds() (int64) {
	currentTime := time.Now()
	timestampSeconds := currentTime.UnixNano()
	return timestampSeconds
}

// 数组成字符串
func Implode(sign string, stringArray []string) (string) {
	result := strings.Join(stringArray, sign)
	return result;
}

// 将字符串解析成数组
func Explode(sign string, dealString string) ([]string) {
	result := strings.Split(dealString, sign)
	return result;
}

// 获取数组最后一个
func ArrayEnd[Any any](arrayData []Any) (Any) {
	last := arrayData[len(arrayData)-1]
	return last
}


// 剔除数组最后一个
func ArrayRemoveEnd[Any any](arrayData []Any) ([]Any) {
	if len(arrayData) <= 1 {
		return nil
	}
	result := arrayData[:len(arrayData)-1]
	return result
}

// 获取文件的内容
func FileGetContent(filePath string) (string) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Dump("无法读取文件:", err)
		return ""
	}
	// 将字节切片转换为字符串
	text := string(content)
	return text;
}

// ymal转json
func MakeYmlToJson(yamlContent []byte) ([]byte, error) {
	jsonContent, err := yaml.YAMLToJSON(yamlContent)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}

// 获取yml文件内容
func FileGetContentYml(filePath string) (Object) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Dump("无法读取文件:", err)
		return nil
	}
	// 转JSON
	jsonContent,_ := MakeYmlToJson(content);

	// 转成map
	result := JsonDecode(string(jsonContent));

	return result;
}

// 四舍五入
func Round(num float64, decimals int) float64 {
	scale := math.Pow(10, float64(decimals))
	return math.Round(num*scale) / scale
}

// 系统睡眠，参数单位为秒
func Sleep(sleepSecond float64) {
	Dump(333);
	sleepSecond = Round(sleepSecond, 2)
	sleepDuration := time.Duration(sleepSecond * 1000) * time.Millisecond
	time.Sleep(sleepDuration) // 程序休眠
}