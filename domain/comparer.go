package domain

import (
	"fmt"
	"reflect"
)

type CompareMode string

type CompareModes []CompareMode

func NewCompareModes(modes []string) CompareModes {
	var result CompareModes
	for _, mode := range modes {
		result = append(result, CompareMode(mode))
	}

	return result
}

type ErrorCode string

type ErrorResult struct {
	Key       string
	LHS       YAMLEntry
	RHS       YAMLEntry
	ErrorCode ErrorCode
}

func (er ErrorResult) FindNilSide() string {
	if er.LHS.Type == "null" {
		return "LHS"
	}

	if er.RHS.Type == "null" {
		return "RHS"
	}

	return ""
}

type ErrorResults []ErrorResult

func (er ErrorResults) IsEmpty() bool {
	return len(er) == 0
}

func TypeUnmatchedResult(key string, lhs any, rhs any) ErrorResult {
	return ErrorResult{
		Key:       key,
		LHS:       NewYAMLEntry(lhs),
		RHS:       NewYAMLEntry(rhs),
		ErrorCode: ErrorTypeUnmatched,
	}
}

func ValueUnmatchedResult(key string, lhs any, rhs any) ErrorResult {
	return ErrorResult{
		Key:       key,
		LHS:       NewYAMLEntry(lhs),
		RHS:       NewYAMLEntry(rhs),
		ErrorCode: ErrorValueUnmatched,
	}
}

func KeyNotFoundResult(key string, lhs any, rhs any) ErrorResult {
	return ErrorResult{
		Key:       key,
		LHS:       NewYAMLEntry(lhs),
		RHS:       NewYAMLEntry(rhs),
		ErrorCode: ErrorKeyNotFound,
	}
}

func IndexNotFoundResult(key string, lhs any, rhs any) ErrorResult {
	return ErrorResult{
		Key:       key,
		LHS:       NewYAMLEntry(lhs),
		RHS:       NewYAMLEntry(rhs),
		ErrorCode: ErrorIndexNotFound,
	}
}

type Results []ErrorResult

type YAMLEntry struct {
	Type  string
	Value string
}

func NewYAMLEntry(entry any) YAMLEntry {
	if entry == nil {
		return YAMLEntry{
			Type:  "null",
			Value: "null",
		}
	}

	var typeString string
	switch reflect.TypeOf(entry).Kind() {
	case reflect.Array, reflect.Slice:
		typeString = "array"
	case reflect.Map:
		typeString = "map"
	case reflect.String:
		typeString = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		typeString = "int"
	case reflect.Float32, reflect.Float64:
		typeString = "float"
	case reflect.Bool:
		typeString = "bool"
	case reflect.Struct:
		typeString = "object"
	default:
		typeString = reflect.TypeOf(entry).Kind().String()
	}

	return YAMLEntry{
		Type:  typeString,
		Value: fmt.Sprintf("%v", entry),
	}
}
