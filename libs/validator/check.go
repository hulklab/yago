package validator

type Rule struct {
	Params  []string
	Method  interface{}
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
