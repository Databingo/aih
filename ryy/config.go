package main

var (
	config map[string]interface{}
)

func initConfig() {
	config = map[string]interface{}{
		"tab_width":     float64(4),
		"tab_to_spaces": true,
	}
}

func configGet(key string, b *Buffer) string {
	if v, ok := config[key]; ok {
		if vv, ok := v.(string); ok {
			return vv
		}
	}
	return ""
}

func configGetBool(key string, b *Buffer) bool {
	if v, ok := config[key]; ok {
		if vv, ok := v.(bool); ok {
			return vv
		}
	}
	return false
}

func configGetNumber(key string, b *Buffer) float64 {
	if v, ok := config[key]; ok {
		if vv, ok := v.(float64); ok {
			return vv
		}
	}
	return 0
}

func configSet(key string, value interface{}) {
	config[key] = value
}
