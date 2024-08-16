package mongo

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func GetFloat64(data interface{}, key string) float64 {
	if m, ok := data.(bson.M); ok {
		if val, ok := m[key].(float64); ok {
			return val
		}
	}
	return 0
}

func getFloat32(data interface{}) float32 {
	if val, ok := data.(float64); ok {
		return float32(val)
	}
	return 0
}

func getString(data interface{}) string {
	if val, ok := data.(string); ok {
		return val
	}
	return ""
}

func getFloat64(data interface{}) float64 {
	if val, ok := data.(float64); ok {
		return val
	}
	return 0
}

func safeString(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}

func safeFloat64(val interface{}) float64 {
	if val == nil {
		return 0.0
	}
	return val.(float64)
}

func safeFloat32(val interface{}) float32 {
	if val == nil {
		return 0.0
	}
	return float32(val.(float64))
}

func getBoolean(isRead interface{}) bool {
	switch v := isRead.(type) {
	case string:
		value, err := strconv.ParseBool(v)
		if err == nil {
			return value
		}
	case int:
		return v != 0
	case bool:
		return v
	}
	return false
}
