package modules

import (
	"fmt"
	"strings"
	"time"
)

// Common helper functions shared across modules

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	if n, ok := val.(int); ok {
		return n != 0
	}
	if n, ok := val.(int64); ok {
		return n != 0
	}
	if s, ok := val.(string); ok {
		return s != ""
	}
	return true
}

func areEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Try direct comparison
	switch aVal := a.(type) {
	case string:
		if bVal, ok := b.(string); ok {
			return aVal == bVal
		}
	case int:
		if bVal, ok := b.(int); ok {
			return aVal == bVal
		}
		if bVal, ok := b.(int64); ok {
			return int64(aVal) == bVal
		}
		if bVal, ok := b.(float64); ok {
			return float64(aVal) == bVal
		}
	case int64:
		if bVal, ok := b.(int64); ok {
			return aVal == bVal
		}
		if bVal, ok := b.(int); ok {
			return aVal == int64(bVal)
		}
		if bVal, ok := b.(float64); ok {
			return float64(aVal) == bVal
		}
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal == bVal
		}
		if bVal, ok := b.(int); ok {
			return aVal == float64(bVal)
		}
		if bVal, ok := b.(int64); ok {
			return aVal == float64(bVal)
		}
	case bool:
		if bVal, ok := b.(bool); ok {
			return aVal == bVal
		}
	}

	return false
}

func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	if n, ok := val.(int); ok {
		return fmt.Sprintf("%d", n)
	}
	if n, ok := val.(int64); ok {
		return fmt.Sprintf("%d", n)
	}
	if n, ok := val.(float64); ok {
		return fmt.Sprintf("%g", n)
	}
	if b, ok := val.(bool); ok {
		if b {
			return "true"
		}
		return "false"
	}
	return fmt.Sprintf("%v", val)
}

func toLowerCase(val interface{}) string {
	s := toString(val)
	return strings.ToLower(s)
}

func randInt(n int) int {
	// Simple pseudo-random implementation
	return int(time.Now().UnixNano()) % n
}

func forthicError(msg string) error {
	return fmt.Errorf("%s", msg)
}

func toInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}
