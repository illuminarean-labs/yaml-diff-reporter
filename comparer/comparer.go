package comparer

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/YangTaeyoung/yaml-diff-reporter/domain"
	"github.com/samber/lo"
)

const (
	Type  domain.CompareMode = "type"
	Key   domain.CompareMode = "key"
	Index domain.CompareMode = "index"
	Value domain.CompareMode = "value"
)

type Comparer interface {
	Compare(currentKey string, lhs any, rhs any)
	Results() *domain.ErrorResults
}

func New(config Config) Comparer {
	c := comparer{
		results: &domain.ErrorResults{},
		config:  config,
	}

	return c
}

type Config struct {
	IgnoredKeys []string
	Modes       domain.CompareModes
}

type comparer struct {
	results *domain.ErrorResults
	config  Config
}

func mapKey(parent string, key string) string {
	if parent == "" {
		return key
	}

	return fmt.Sprintf("%s.%s", parent, key)
}

func sliceKey(parent string, index int) string {
	return fmt.Sprintf("%s[%d]", parent, index)
}

func (c comparer) Compare(parent string, lhs any, rhs any) {
	if lo.Contains(c.config.IgnoredKeys, parent) {
		return
	}

	lhsType := reflect.TypeOf(lhs)
	rhsType := reflect.TypeOf(rhs)

	if lo.Contains(c.config.Modes, Type) && lhsType != rhsType {
		*c.results = append(*c.results, domain.TypeUnmatchedResult(parent, lhs, rhs))
		return
	}

	switch lhs.(type) {
	case map[string]any:
		lhsMap, _ := lhs.(map[string]any)
		rhsMap, _ := rhs.(map[string]any)
		c.compareMap(parent, lhsMap, rhsMap)
	case []any:
		lhsArr, _ := lhs.([]any)
		rhsArr, _ := rhs.([]any)
		c.compareSlice(parent, lhsArr, rhsArr)
	default:
		if lo.Contains(c.config.Modes, Value) && !reflect.DeepEqual(lhs, rhs) {
			*c.results = append(*c.results, domain.ValueUnmatchedResult(parent, lhs, rhs))
		}
	}
}

func (c comparer) compareMap(parent string, lhs map[string]any, rhs map[string]any) {
	var (
		ok     bool
		key    string
		lhsVal any
		rhsVal any
	)
	visited := make(map[string]bool)

	for key, lhsVal = range lhs {
		nextKey := mapKey(parent, key)
		if lo.Contains(c.config.IgnoredKeys, nextKey) {
			continue
		}

		visited[key] = true
		rhsVal, ok = rhs[key]
		if !ok {
			if lo.Contains(c.config.Modes, Key) {
				*c.results = append(*c.results, domain.KeyNotFoundResult(nextKey, lhsVal, nil))
			}

			continue
		}

		c.Compare(nextKey, lhsVal, rhsVal)
	}
	for key, rhsVal = range rhs {
		nextKey := mapKey(parent, key)
		if lo.Contains(c.config.IgnoredKeys, nextKey) {
			continue
		}

		if visited[key] {
			continue
		}

		lhsVal, ok = lhs[key]
		if !ok {
			if lo.Contains(c.config.Modes, Key) {
				*c.results = append(*c.results, domain.KeyNotFoundResult(nextKey, nil, rhsVal))
			}
		}
	}
}

func (c comparer) compareSlice(parent string, lhs []any, rhs []any) {
	var (
		idx    int
		lhsVal any
		rhsVal any
	)

	visited := make(map[int]bool)
	for idx, lhsVal = range lhs {
		nextKey := sliceKey(parent, idx)
		if lo.Contains(c.config.IgnoredKeys, nextKey) {
			continue
		}

		visited[idx] = true
		if len(rhs) <= idx {
			if lo.Contains(c.config.Modes, Index) {
				*c.results = append(*c.results, domain.IndexNotFoundResult(nextKey, lhsVal, nil))
			}

			continue
		}
		rhsVal = rhs[idx]

		c.Compare(nextKey, lhsVal, rhsVal)
	}

	for idx, rhsVal = range rhs {
		nextKey := sliceKey(parent, idx)
		if lo.Contains(c.config.IgnoredKeys, nextKey) {
			continue
		}
		if visited[idx] {
			continue
		}

		if len(lhs) <= idx {
			if lo.Contains(c.config.Modes, Index) {
				*c.results = append(*c.results, domain.IndexNotFoundResult(nextKey, nil, rhsVal))
			}
		}
	}
}

func (c comparer) Results() *domain.ErrorResults {
	sort.SliceStable(*c.results, func(i, j int) bool {
		return (*c.results)[i].ErrorCode < (*c.results)[j].ErrorCode
	})

	return c.results
}
