package nodes

import "fmt"

// NodeOutput holds the output of a node execution.
// It provides field access for variable interpolation.
type NodeOutput struct {
	fields map[string]interface{}
}

// NewNodeOutput creates a new node output with the given fields.
func NewNodeOutput(fields map[string]interface{}) NodeOutput {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	return NodeOutput{fields: fields}
}

// EmptyOutput returns an empty node output.
func EmptyOutput() NodeOutput {
	return NodeOutput{fields: make(map[string]interface{})}
}

// GetField returns the string value of a field for variable interpolation.
func (o NodeOutput) GetField(name string) string {
	if o.fields == nil {
		return ""
	}
	val, ok := o.fields[name]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

// GetFieldPath navigates a nested path like ["headers", "Content-Type"] and returns the value.
func (o NodeOutput) GetFieldPath(path []string) string {
	if len(path) == 0 || o.fields == nil {
		return ""
	}

	// Start with the first path element
	current, ok := o.fields[path[0]]
	if !ok {
		return ""
	}

	// Navigate through remaining path elements
	for _, key := range path[1:] {
		switch v := current.(type) {
		case map[string]interface{}:
			current, ok = v[key]
			if !ok {
				return ""
			}
		case map[string]string:
			val, exists := v[key]
			if !exists {
				return ""
			}
			return val
		default:
			return "" // Can't navigate further
		}
	}

	return fmt.Sprintf("%v", current)
}

// Get returns the raw value of a field.
func (o NodeOutput) Get(name string) (interface{}, bool) {
	if o.fields == nil {
		return nil, false
	}
	val, ok := o.fields[name]
	return val, ok
}

// Set sets a field value.
func (o *NodeOutput) Set(name string, value interface{}) {
	if o.fields == nil {
		o.fields = make(map[string]interface{})
	}
	o.fields[name] = value
}

// Preview returns a short preview of the output for debugging/UI.
func (o NodeOutput) Preview() string {
	if o.fields == nil {
		return ""
	}

	// Try common field names for preview
	for _, key := range []string{"text", "rawText", "processedText", "output", "message", "stdout"} {
		if val, ok := o.fields[key]; ok {
			str := fmt.Sprintf("%v", val)
			if len(str) > 100 {
				return str[:100] + "..."
			}
			return str
		}
	}

	return ""
}

// Fields returns a copy of all fields.
func (o NodeOutput) Fields() map[string]interface{} {
	if o.fields == nil {
		return make(map[string]interface{})
	}
	copy := make(map[string]interface{}, len(o.fields))
	for k, v := range o.fields {
		copy[k] = v
	}
	return copy
}
