package qppg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/snowzach/gorestapi/queryp"
	"github.com/stretchr/testify/assert"
)

func TestQueryParser(t *testing.T) {

	q, err := queryp.ParseQuery("field=value&((another=<value|yet=another1|limit=weee))|third=value&limit=10&option=beans&sort=test,-another")
	assert.Nil(t, err)
	fmt.Println(q.PrettyString())

	var queryClause strings.Builder
	var queryParams = []interface{}{}

	filterFields := queryp.FilterFieldTypes{
		"another":           queryp.FilterTypeString,
		"yet":               queryp.FilterTypeString,
		"field":             queryp.FilterTypeString,
		"limit":             queryp.FilterTypeNumeric,
		"third":             queryp.FilterTypeString,
		"thing.id":          queryp.FilterTypeString,
		"thing.name":        queryp.FilterTypeString,
		"thing.description": queryp.FilterTypeString,
	}

	err = FilterQuery(filterFields, q.Filter, &queryClause, &queryParams)
	assert.Nil(t, err)

}
