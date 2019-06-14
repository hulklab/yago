package validator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Rule struct {
	Params  []string
	Method  int
	On      []string
	Min     float64
	Max     float64
	Pattern string
	Message string
}

type Label map[string]string

func (l *Label) Get(key string) string {
	if value, ok := (*l)[key]; ok {
		return value
	}
	return key
}

func ValidateHttp(c *gin.Context, action string, labels Label, rules []Rule) (bool, error) {
	for _, rule := range rules {
		actionMatch := false
		if len(rule.On) == 0 {
			actionMatch = true
		} else {
			for _, a := range rule.On {
				if a == action {
					actionMatch = true
					break
				}
			}
		}

		if actionMatch {
			switch rule.Method {
			case Required:
				for _, p := range rule.Params {
					pv, exist := c.Get(p)
					if !exist {
						return false, fmt.Errorf("%s 不存在", labels.Get(p))
					}
					if valid, err := (RequiredValidator{}).Check(pv); !valid {
						return false, getErr(labels.Get(p), err, rule.Message)
					}
				}
			case String:
				for _, p := range rule.Params {
					pv, exist := c.Get(p)
					if !exist {
						continue
					}
					if valid, err := (StringValidator{Min: int(rule.Min), Max: int(rule.Max)}).Check(pv); !valid {
						return false, getErr(labels.Get(p), err, rule.Message)
					}
				}
			case Int:
				for _, p := range rule.Params {
					pv, exist := c.Get(p)
					if !exist {
						continue
					}
					pvInt, err := strconv.Atoi(pv.(string))
					if err != nil {
						return false, fmt.Errorf("%s 不是个整数", labels.Get(p))
					}
					if valid, err := (IntValidator{Min: int(rule.Min), Max: int(rule.Max)}).Check(pvInt); !valid {
						return false, getErr(labels.Get(p), err, rule.Message)
					}
				}
			case Float:
				for _, p := range rule.Params {
					pv, exist := c.Get(p)
					if !exist {
						continue
					}

					pvFloat, err := strconv.ParseFloat(pv.(string), 64)
					if err != nil {
						return false, fmt.Errorf("%s 不是个浮点数", labels.Get(p))
					}
					if valid, err := (FloatValidator{Min: rule.Min, Max: rule.Max}).Check(pvFloat); !valid {
						return false, getErr(labels.Get(p), err, rule.Message)
					}
				}
			case JSON:
				for _, p := range rule.Params {
					pv, exist := c.Get(p)
					if !exist {
						continue
					}
					if valid, err := (JSONValidator{}).Check(pv); !valid {
						return false, getErr(labels.Get(p), err, rule.Message)
					}
				}
			case IP:
				for _, p := range rule.Params {
					pv, exist := c.Get(p)
					if !exist {
						continue
					}
					if valid, err := (IPValidator{}).Check(pv); !valid {
						return false, getErr(labels.Get(p), err, rule.Message)
					}
				}
			case Match:
				for _, p := range rule.Params {
					pv, exist := c.Get(p)
					if !exist {
						continue
					}
					if valid, err := (MatchValidator{Pattern: rule.Pattern}).Check(pv); !valid {
						return false, getErr(labels.Get(p), err, rule.Message)
					}
				}
			}
		}
	}
	return true, nil
}

func getErr(label string, err error, message string) error {
	if message == "" {
		return fmt.Errorf("%s %s", label, err)
	}
	return fmt.Errorf("%s %s", label, message)

}
