package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

const (
	Required = iota
	Number
	String
	JSON
	IP
	Match
)

type Validator interface {
	Check(v interface{}) (bool, error)
}

type StringValidator struct {
	Min int
	Max int
}

func (v StringValidator) Check(value interface{}) (bool, error) {
	str, ok := value.(string)
	if !ok {
		return false, errors.New("不是个字符串")
	}
	strLen := len(str)
	if strLen == 0 {
		return false, fmt.Errorf("不能为空")
	}

	if v.Min != 0 && strLen < v.Min {
		return false, fmt.Errorf("至少应该有 %v 个字符长", v.Min)
	}

	if v.Max != 0 && v.Max >= v.Min && strLen > v.Max {
		return false, fmt.Errorf("最大不能超过 %v 个字符长", v.Max)
	}
	return true, nil
}

type NumberValidator struct {
	Min int
	Max int
}

func (v NumberValidator) Check(value interface{}) (bool, error) {
	num, ok := value.(int)

	if !ok {
		return false, errors.New("不是个整数")
	}

	if v.Min != 0 && num < v.Min {
		return false, fmt.Errorf("最小不能小于 %v", v.Min)
	}

	if v.Max != 0 && v.Max >= v.Min && num > v.Max {
		return false, fmt.Errorf("最大不能大于 %v", v.Max)
	}

	return true, nil
}

type IPValidator struct {
}

func (v IPValidator) Check(value interface{}) (bool, error) {
	isIP, _ := regexp.Match(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`, value.([]byte))
	if isIP == false {
		return isIP, errors.New("不是个IP地址")
	}
	return true, nil
}

type JSONValidator struct {
}

func (v JSONValidator) Check(value interface{}) (bool, error) {
	if !json.Valid(value.([]byte)) {
		return false, errors.New("不是一个有效的JSON串")
	}
	return true, nil
}

type MatchValidator struct {
	Pattern string
}

func (v MatchValidator) Check(value interface{}) (bool, error) {
	matched, _ := regexp.Match(v.Pattern, value.([]byte))
	if !matched {
		return false, errors.New("不合法")
	}
	return true, nil
}

type RequiredValidator struct {
}

func (v RequiredValidator) Check(value interface{}) (bool, error) {
	str := value.(string)
	if str == "" {
		return false, errors.New("不能为空")
	}
	return true, nil
}
