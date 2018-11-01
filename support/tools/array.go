package tools

import (
	"fmt"
	"reflect"
	"strings"
)

func Find(findStr string, array []string) bool {

	for _, val := range array {
		if val == findStr {
			return true
		}
	}
	return false
}

func IsArr(params interface{}) bool {

	rfValue := reflect.ValueOf(params)
	if rfValue.Kind() == reflect.Slice || rfValue.Kind() == reflect.Array {
		return true
	}
	return false
}

func Split(param interface{}, separator string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(param), "[]"), " ", separator, -1)
}
