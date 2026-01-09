package nodes

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ConditionNode branches workflow based on a condition.
type ConditionNode struct{}

// NewConditionNode creates a new condition node.
func NewConditionNode() *ConditionNode {
	return &ConditionNode{}
}

// Type returns the node type identifier.
func (n *ConditionNode) Type() string {
	return "condition"
}

// Execute evaluates a condition and returns the branch result.
func (n *ConditionNode) Execute(ctx context.Context, input NodeInput, services ServiceProvider) (NodeOutput, error) {
	// Get condition configuration
	field, _ := input.Config["field"].(string)
	operator, _ := input.Config["operator"].(string)
	expectedValue := input.Config["value"]

	// Resolve the field value
	var actualValue interface{}

	// Check for settings fields
	if strings.HasPrefix(field, "settings.") {
		settingsField := strings.TrimPrefix(field, "settings.")
		settings := services.GetSettingsManager().Get()
		actualValue = getSettingsField(settings, settingsField)
	} else {
		actualValue = input.Config["actualValue"]
	}

	// Evaluate condition
	result := evaluateCondition(actualValue, operator, expectedValue)

	return NewNodeOutput(map[string]interface{}{
		"result": result,
		"branch": boolToString(result),
	}), nil
}

// getSettingsField extracts a field value from settings using reflection.
func getSettingsField(settings interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(settings)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			tagName := strings.Split(jsonTag, ",")[0]
			if tagName == fieldName {
				return v.Field(i).Interface()
			}
		}
		if strings.EqualFold(field.Name, fieldName) {
			return v.Field(i).Interface()
		}
	}

	return nil
}

// evaluateCondition evaluates a condition with the given operator.
func evaluateCondition(actual interface{}, operator string, expected interface{}) bool {
	switch operator {
	case "equals", "==", "eq":
		return reflect.DeepEqual(actual, expected)
	case "notEquals", "!=", "ne":
		return !reflect.DeepEqual(actual, expected)
	case "contains":
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		return strings.Contains(actualStr, expectedStr)
	case "startsWith":
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		return strings.HasPrefix(actualStr, expectedStr)
	case "endsWith":
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		return strings.HasSuffix(actualStr, expectedStr)
	case "isEmpty":
		return isEmptyValue(actual)
	case "isNotEmpty":
		return !isEmptyValue(actual)
	case "greaterThan", ">", "gt":
		return compareNumeric(actual, expected) > 0
	case "lessThan", "<", "lt":
		return compareNumeric(actual, expected) < 0
	default:
		return isTruthy(actual)
	}
}

func isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case []interface{}:
		return len(val) == 0
	case map[string]interface{}:
		return len(val) == 0
	}
	return false
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val != "" && val != "false" && val != "0"
	case int:
		return val != 0
	case float64:
		return val != 0
	}
	return true
}

func compareNumeric(a, b interface{}) int {
	aFloat := toFloat(a)
	bFloat := toFloat(b)
	if aFloat < bFloat {
		return -1
	}
	if aFloat > bFloat {
		return 1
	}
	return 0
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}
	return 0
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
