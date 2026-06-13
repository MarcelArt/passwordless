package common

import (
	"fmt"
	"strings"
)

type DynamicRawQuery struct {
	baseQuery string
	wheres    []string
	values    []any
}

// baseQuery need to consist of a %s placeholder for the where clause
func NewDynamicRawQuery(baseQuery string, initialValues ...any) *DynamicRawQuery {
	return &DynamicRawQuery{baseQuery: baseQuery, values: initialValues}
}

func (d *DynamicRawQuery) AddWhere(q string, v ...any) *DynamicRawQuery {
	d.values = append(d.values, v...)
	d.wheres = append(d.wheres, q)
	return d
}

func (d *DynamicRawQuery) Build() (string, []any) {
	wheres := strings.Join(d.wheres, "\n")
	return fmt.Sprintf(d.baseQuery, wheres), d.values
}
