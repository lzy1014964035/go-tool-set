package config

import (
	"reflect"
	"github.com/lzy1014964035/go-tool-set/service"
)

// 获取配置
func Get(filePath string, configPath string) service.Any {
	service.Dump("获取配置");
	configMap := service.FileGetContentYml(filePath);
	configPathArray := service.Explode(".", configPath);
	result := getPathValue(configMap, configPathArray);
	return result;
}

// 根据路径获取值
func getPathValue(data service.Any, pathArray []string) service.Any {
	if len(pathArray) == 0 {
		return data
	}
	key := pathArray[0]
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Map:
		if value.MapIndex(reflect.ValueOf(key)).IsValid() {
			return getPathValue(value.MapIndex(reflect.ValueOf(key)).Interface(), pathArray[1:])
		}
	case reflect.Struct:
		field := value.FieldByName(key)
		if field.IsValid() {
			return getPathValue(field.Interface(), pathArray[1:])
		}
	}

	return nil
}