package abi

import (
	"fmt"
	"strings"
)

// FormatAbiItem returns the human-readable signature of an ABI item.
// Works for functions, events, and errors.
//
// Example:
//
//	sig := FormatAbiItem(function) // "transfer(address,uint256)"
func FormatAbiItem(item any) (string, error) {
	switch v := item.(type) {
	case Function:
		return formatFunction(v), nil
	case *Function:
		return formatFunction(*v), nil
	case Event:
		return formatEvent(v), nil
	case *Event:
		return formatEvent(*v), nil
	case Error:
		return formatError(v), nil
	case *Error:
		return formatError(*v), nil
	default:
		return "", fmt.Errorf("unsupported ABI item type: %T", item)
	}
}

// FormatAbiParams formats a list of parameters as a comma-separated string.
// If includeName is true, parameter names are included.
//
// Example:
//
//	params := FormatAbiParams(function.Inputs, false) // "address,uint256"
//	params := FormatAbiParams(function.Inputs, true)  // "address to, uint256 amount"
func FormatAbiParams(params []Parameter, includeName bool) string {
	parts := make([]string, len(params))
	for i, param := range params {
		parts[i] = formatParam(param, includeName)
	}
	if includeName {
		return strings.Join(parts, ", ")
	}
	return strings.Join(parts, ",")
}

// formatFunction formats a function to its signature.
func formatFunction(fn Function) string {
	return fmt.Sprintf("%s(%s)", fn.Name, FormatAbiParams(fn.Inputs, false))
}

// formatEvent formats an event to its signature.
func formatEvent(ev Event) string {
	return fmt.Sprintf("%s(%s)", ev.Name, FormatAbiParams(ev.Inputs, false))
}

// formatError formats an error to its signature.
func formatError(err Error) string {
	return fmt.Sprintf("%s(%s)", err.Name, FormatAbiParams(err.Inputs, false))
}

// formatParam formats a single parameter.
func formatParam(param Parameter, includeName bool) string {
	typ := formatParamType(param)
	if includeName && param.Name != "" {
		return typ + " " + param.Name
	}
	return typ
}

// formatParamType formats the type of a parameter, handling tuples.
func formatParamType(param Parameter) string {
	// Handle tuple types
	if strings.HasPrefix(param.Type, "tuple") {
		components := make([]string, len(param.Components))
		for i, comp := range param.Components {
			components[i] = formatParamType(comp)
		}
		// Extract any array suffix (e.g., "tuple[]" -> "[]")
		suffix := strings.TrimPrefix(param.Type, "tuple")
		return fmt.Sprintf("(%s)%s", strings.Join(components, ","), suffix)
	}
	return param.Type
}

// FormatFunctionSignature formats a function name and inputs to a signature.
func FormatFunctionSignature(name string, inputs []Parameter) string {
	return fmt.Sprintf("%s(%s)", name, FormatAbiParams(inputs, false))
}

// FormatEventSignature formats an event name and inputs to a signature.
func FormatEventSignature(name string, inputs []Parameter) string {
	return fmt.Sprintf("%s(%s)", name, FormatAbiParams(inputs, false))
}

// FormatErrorSignature formats an error name and inputs to a signature.
func FormatErrorSignature(name string, inputs []Parameter) string {
	return fmt.Sprintf("%s(%s)", name, FormatAbiParams(inputs, false))
}

// GetFunctionSignature returns the signature of a function in the ABI.
func (a *ABI) GetFunctionSignature(name string) (string, error) {
	fn, ok := a.Functions[name]
	if !ok {
		return "", fmt.Errorf("function %q not found on ABI", name)
	}
	return FormatFunctionSignature(fn.Name, fn.Inputs), nil
}

// GetEventSignature returns the signature of an event in the ABI.
func (a *ABI) GetEventSignature(name string) (string, error) {
	ev, ok := a.Events[name]
	if !ok {
		return "", fmt.Errorf("event %q not found on ABI", name)
	}
	return FormatEventSignature(ev.Name, ev.Inputs), nil
}

// GetErrorSignature returns the signature of an error in the ABI.
func (a *ABI) GetErrorSignature(name string) (string, error) {
	err, ok := a.Errors[name]
	if !ok {
		return "", fmt.Errorf("error %q not found on ABI", name)
	}
	return FormatErrorSignature(err.Name, err.Inputs), nil
}
