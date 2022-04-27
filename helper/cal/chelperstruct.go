package cal

import (
	"reflect"

	"github.com/gogf/gf/v2/util/gconv"
)

//反射获取字段名 按顺序
func GetField(info interface{}) (dataArr []string) {
	t := reflect.TypeOf(info)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	dataArr = make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		valueName := t.Field(i).Name
		dataArr = append(dataArr, valueName)
	}
	return
}

//反射获取字段值 按顺序
func GetValue(info interface{}) (dataArr []string) {
	v := reflect.ValueOf(info)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	dataArr = make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		valueName := gconv.String(v.Field(i).Interface())
		dataArr = append(dataArr, valueName)
	}
	return
}

//反射获取有Tag的字段名 按顺序
func GetTagField(info interface{}, tagname string) (dataArr []string) {
	t := reflect.TypeOf(info)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	dataArr = make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		if tagValue := t.Field(i).Tag.Get(tagname); tagValue != "" {
			dataArr = append(dataArr, tagValue)
		}
	}
	return
}

//反射获取有Tag的字段值 按顺序
func GetTagValue(info interface{}, tagname string) (dataArr []string) {
	v := reflect.ValueOf(info)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	dataArr = make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		if tagValue := t.Field(i).Tag.Get(tagname); tagValue != "" {
			valueName := gconv.String(v.Field(i).Interface())
			dataArr = append(dataArr, valueName)
		}
	}
	return
}
