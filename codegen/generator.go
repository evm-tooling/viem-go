package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/ChefBingbong/viem-go/abi"
)

// Generator generates Go code from an ABI.
type Generator struct {
	packageName  string
	contractName string
	abi          *abi.ABI
	abiJSON      []byte
}

// NewGenerator creates a new code generator.
func NewGenerator(packageName, contractName string, abiJSON []byte) (*Generator, error) {
	parsedABI, err := abi.Parse(abiJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return &Generator{
		packageName:  packageName,
		contractName: contractName,
		abi:          parsedABI,
		abiJSON:      abiJSON,
	}, nil
}

// Generate generates the Go code for the contract.
func (g *Generator) Generate() ([]byte, error) {
	data := g.buildTemplateData()

	tmpl, err := template.New("contract").Funcs(templateFuncs).Parse(contractTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Return unformatted if formatting fails (for debugging)
		return buf.Bytes(), fmt.Errorf("failed to format generated code: %w (returning unformatted)", err)
	}

	return formatted, nil
}

// TemplateData holds all data needed for code generation.
type TemplateData struct {
	PackageName  string
	ContractName string
	ABIJSON      string
	Functions    []FunctionData
	Events       []EventData
	HasEvents    bool
}

// FunctionData holds data for a single function.
type FunctionData struct {
	Name            string
	GoName          string
	Inputs          []ParamData
	Outputs         []ParamData
	IsReadOnly      bool
	StateMutability string
	Signature       string
}

// EventData holds data for a single event.
type EventData struct {
	Name      string
	GoName    string
	Inputs    []ParamData
	Signature string
}

// ParamData holds data for a parameter.
type ParamData struct {
	Name   string
	GoName string
	Type   string
	GoType string
}

// buildTemplateData builds the data structure for templates.
func (g *Generator) buildTemplateData() TemplateData {
	data := TemplateData{
		PackageName:  g.packageName,
		ContractName: g.contractName,
		ABIJSON:      escapeBackticks(string(g.abiJSON)),
		HasEvents:    len(g.abi.Events) > 0,
	}

	// Process functions
	for _, fn := range g.abi.Functions {
		fnData := FunctionData{
			Name:            fn.Name,
			GoName:          toExportedName(fn.Name),
			IsReadOnly:      fn.IsReadOnly(),
			StateMutability: fn.StateMutability.String(),
			Signature:       fn.Signature,
		}

		for i, input := range fn.Inputs {
			name := input.Name
			if name == "" {
				name = fmt.Sprintf("arg%d", i)
			}
			fnData.Inputs = append(fnData.Inputs, ParamData{
				Name:   name,
				GoName: toLowerCamelCase(name),
				Type:   input.Type,
				GoType: solidityToGoType(input.Type),
			})
		}

		for i, output := range fn.Outputs {
			name := output.Name
			if name == "" {
				name = fmt.Sprintf("ret%d", i)
			}
			fnData.Outputs = append(fnData.Outputs, ParamData{
				Name:   name,
				GoName: toLowerCamelCase(name),
				Type:   output.Type,
				GoType: solidityToGoType(output.Type),
			})
		}

		data.Functions = append(data.Functions, fnData)
	}

	// Process events
	for _, ev := range g.abi.Events {
		evData := EventData{
			Name:      ev.Name,
			GoName:    toExportedName(ev.Name),
			Signature: ev.Signature,
		}

		for i, input := range ev.Inputs {
			name := input.Name
			if name == "" {
				name = fmt.Sprintf("arg%d", i)
			}
			evData.Inputs = append(evData.Inputs, ParamData{
				Name:   name,
				GoName: toExportedName(name),
				Type:   input.Type,
				GoType: solidityToGoType(input.Type),
			})
		}

		data.Events = append(data.Events, evData)
	}

	return data
}

// solidityToGoType converts a Solidity type to a Go type.
func solidityToGoType(solType string) string {
	// Handle arrays
	if strings.HasSuffix(solType, "[]") {
		elemType := strings.TrimSuffix(solType, "[]")
		return "[]" + solidityToGoType(elemType)
	}

	// Handle fixed arrays
	if matched, _ := regexp.MatchString(`\[\d+\]$`, solType); matched {
		re := regexp.MustCompile(`^(.+)\[(\d+)\]$`)
		matches := re.FindStringSubmatch(solType)
		if len(matches) == 3 {
			return "[" + matches[2] + "]" + solidityToGoType(matches[1])
		}
	}

	// Handle basic types
	switch {
	case solType == "address":
		return "common.Address"
	case solType == "bool":
		return "bool"
	case solType == "string":
		return "string"
	case solType == "bytes":
		return "[]byte"
	case strings.HasPrefix(solType, "bytes"):
		// bytes1 through bytes32
		size := strings.TrimPrefix(solType, "bytes")
		return "[" + size + "]byte"
	case strings.HasPrefix(solType, "uint"):
		size := strings.TrimPrefix(solType, "uint")
		if size == "" || size == "256" {
			return "*big.Int"
		}
		sizeInt := 0
		fmt.Sscanf(size, "%d", &sizeInt)
		if sizeInt <= 8 {
			return "uint8"
		} else if sizeInt <= 16 {
			return "uint16"
		} else if sizeInt <= 32 {
			return "uint32"
		} else if sizeInt <= 64 {
			return "uint64"
		}
		return "*big.Int"
	case strings.HasPrefix(solType, "int"):
		size := strings.TrimPrefix(solType, "int")
		if size == "" || size == "256" {
			return "*big.Int"
		}
		sizeInt := 0
		fmt.Sscanf(size, "%d", &sizeInt)
		if sizeInt <= 8 {
			return "int8"
		} else if sizeInt <= 16 {
			return "int16"
		} else if sizeInt <= 32 {
			return "int32"
		} else if sizeInt <= 64 {
			return "int64"
		}
		return "*big.Int"
	case strings.HasPrefix(solType, "tuple"):
		// For tuples, we would need more complex handling
		// For now, return interface{}
		return "interface{}"
	default:
		return "interface{}"
	}
}

// toExportedName converts a name to an exported Go identifier.
func toExportedName(name string) string {
	if name == "" {
		return "Value"
	}
	// Handle special cases
	name = strings.ReplaceAll(name, "_", " ")
	words := strings.Fields(name)
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	result := strings.Join(words, "")
	
	// Ensure first character is uppercase
	runes := []rune(result)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}

// toLowerCamelCase converts a name to lowerCamelCase.
func toLowerCamelCase(name string) string {
	if name == "" {
		return "value"
	}
	exported := toExportedName(name)
	runes := []rune(exported)
	if len(runes) > 0 {
		runes[0] = unicode.ToLower(runes[0])
	}
	return string(runes)
}

// escapeBackticks escapes backticks in a string for use in raw string literals.
func escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "` + \"`\" + `")
}

// Template functions
var templateFuncs = template.FuncMap{
	"join": strings.Join,
	"add": func(a, b int) int {
		return a + b
	},
}
