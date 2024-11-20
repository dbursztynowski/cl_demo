package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// A simple JSON Schema structure for demonstration purposes
type JSONSchema struct {
	Type       string                `json:"type"`
	Properties map[string]JSONSchema `json:"properties,omitempty"`
	Items      *JSONSchema           `json:"items,omitempty"`
}

// Function to generate template from JSON Schema
func generateTemplateFromSchema(schema JSONSchema, key string) string {
	var sb strings.Builder

	switch schema.Type {
	case "object":
		sb.WriteString("{\n")
		obj_len := len(schema.Properties)
		i := 1
		for key, prop := range schema.Properties {
			sb.WriteString(fmt.Sprintf("  \"%s\": ", key))
			sb.WriteString(generateTemplateFromSchema(prop, key))

			if i == obj_len {
				sb.WriteString("\n")
			} else {
				sb.WriteString(",\n")
			}
			i++
		}
		sb.WriteString("}")
	case "array":
		sb.WriteString("[\n")
		if schema.Items != nil {
			sb.WriteString(generateTemplateFromSchema(*schema.Items, key))
		}
		sb.WriteString("\n]")
	case "string":
		sb.WriteString("\"{{ .")
		sb.WriteString(key)
		sb.WriteString(" }}\"")
	case "integer":
		sb.WriteString("{{ .")
		sb.WriteString(key)
		sb.WriteString(" }}")
	case "number":
		sb.WriteString("{{ .")
		sb.WriteString(key)
		sb.WriteString(" }}")
	}

	return sb.String()
}

func ExtractSubJSON(data []byte, key string) ([]byte, error) {

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if subJSON, found := result[key]; found {
		return json.Marshal(subJSON)
	}

	return nil, fmt.Errorf("key %s not found", key)
}
