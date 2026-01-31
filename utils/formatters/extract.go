package formatters

// ExtractFormatter is a function type that formats input data.
type ExtractFormatter func(value map[string]any) map[string]any

// Extract picks out the keys from value that exist in the formatted output.
// This is useful for extracting only the relevant fields from a value
// based on what a formatter produces.
//
// Example:
//
//	value := map[string]any{
//		"gasPrice": "0x1",
//		"maxFeePerGas": "0x2",
//		"customField": "ignored",
//	}
//	formatter := func(v map[string]any) map[string]any {
//		return map[string]any{
//			"gasPrice": v["gasPrice"],
//		}
//	}
//	extracted := Extract(value, formatter)
//	// extracted = {"gasPrice": "0x1"}
func Extract(value map[string]any, format ExtractFormatter) map[string]any {
	if format == nil {
		return map[string]any{}
	}

	result := make(map[string]any)
	formatted := format(value)

	extract(value, formatted, result)

	return result
}

// extract recursively extracts keys from value that exist in formatted.
func extract(value, formatted, result map[string]any) {
	for key, formattedValue := range formatted {
		// If the key exists in value, add it to result
		if v, ok := value[key]; ok {
			result[key] = v
		}

		// Recursively handle nested objects
		if formattedValue != nil {
			if nestedFormatted, ok := formattedValue.(map[string]any); ok {
				extract(value, nestedFormatted, result)
			}
		}
	}
}

// ExtractKeys extracts only the specified keys from a map.
func ExtractKeys(value map[string]any, keys []string) map[string]any {
	result := make(map[string]any)
	for _, key := range keys {
		if v, ok := value[key]; ok {
			result[key] = v
		}
	}
	return result
}

// OmitKeys returns a map with the specified keys omitted.
func OmitKeys(value map[string]any, keys []string) map[string]any {
	result := make(map[string]any)
	keySet := make(map[string]bool)
	for _, key := range keys {
		keySet[key] = true
	}

	for key, v := range value {
		if !keySet[key] {
			result[key] = v
		}
	}
	return result
}
