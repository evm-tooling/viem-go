package utils

import (
	"math/big"
	"reflect"

	json "github.com/goccy/go-json"
)

// Stringify is a JSON marshaling function that handles *big.Int values
// by converting them to strings, similar to how viem handles BigInt in JavaScript.
//
// Example:
//
//	type Data struct {
//		Value *big.Int `json:"value"`
//	}
//	result, _ := Stringify(Data{Value: big.NewInt(123456789012345)})
//	// result: {"value":"123456789012345"}
func Stringify(v any) (string, error) {
	converted := convertBigInts(v)
	bytes, err := json.Marshal(converted)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// StringifyIndent is like Stringify but with indentation.
func StringifyIndent(v any, prefix, indent string) (string, error) {
	converted := convertBigInts(v)
	bytes, err := json.MarshalIndent(converted, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// StringifyBytes returns the JSON as bytes instead of a string.
func StringifyBytes(v any) ([]byte, error) {
	converted := convertBigInts(v)
	return json.Marshal(converted)
}

// convertBigInts recursively converts *big.Int values to strings.
func convertBigInts(v any) any {
	if v == nil {
		return nil
	}

	// Handle *big.Int directly
	if bi, ok := v.(*big.Int); ok {
		if bi == nil {
			return nil
		}
		return bi.String()
	}

	// Handle big.Int directly
	if bi, ok := v.(big.Int); ok {
		return bi.String()
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return convertBigInts(val.Elem().Interface())

	case reflect.Struct:
		result := make(map[string]any)
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if !field.IsExported() {
				continue
			}
			// Use json tag if present
			name := field.Tag.Get("json")
			if name == "" {
				name = field.Name
			} else if name == "-" {
				continue
			}
			// Handle ",omitempty" and other json tag options
			if idx := len(name); idx > 0 {
				for j, c := range name {
					if c == ',' {
						name = name[:j]
						break
					}
				}
			}
			result[name] = convertBigInts(val.Field(i).Interface())
		}
		return result

	case reflect.Map:
		result := make(map[string]any)
		iter := val.MapRange()
		for iter.Next() {
			key := iter.Key()
			var keyStr string
			if key.Kind() == reflect.String {
				keyStr = key.String()
			} else {
				keyStr = reflect.ValueOf(key.Interface()).String()
			}
			result[keyStr] = convertBigInts(iter.Value().Interface())
		}
		return result

	case reflect.Slice, reflect.Array:
		result := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = convertBigInts(val.Index(i).Interface())
		}
		return result

	default:
		return v
	}
}
