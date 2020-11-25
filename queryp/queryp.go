package queryp

import (
	"bytes"
	"encoding/json"
	"strings"
)

type QueryParameters struct {
	Filter  Filter  `json:"filter"`
	Sort    Sort    `json:"sort"`
	Options Options `json:"options"`
	Limit   int     `json:"limit"`
	Offset  int     `json:"offset"`
}

type Filter []FilterTerm

type FilterTerm struct {
	Logic FilterLogic `json:"logic"`
	Op    FilterOp    `json:"op"`
	Field Field       `json:"field"`
	Value string      `json:"value"`

	SubFilter Filter `json:"sub_filter,omitempty"`
}

const (
	FilterTypeNotFound FilterType = iota
	FilterTypeSimple
	FilterTypeString
	FilterTypeNumeric
	FilterTypeTime
	FilterTypeBool
)

type Field = string // Alias
type FilterType int
type FilterFieldTypes map[Field]FilterType

type SortFields []string

type Sort []SortTerm

type SortTerm struct {
	Field Field `json:"field"`
	Desc  bool  `json:"desc"`
}

type Options map[string]struct{} // Just a lookup of string

func (o Options) HasOption(option string) bool {
	_, found := o[option]
	return found
}

func (fft FilterFieldTypes) FindFilterType(search string) (Field, FilterType) {

	if filterType, found := fft[search]; found {
		return search, filterType
	}

	// Search for the suffix in the list of filters
	search = "." + search
	for filter, filterType := range fft {
		if strings.HasSuffix(filter, search) {
			return filter, filterType
		}
	}
	return "", FilterTypeNotFound

}

func (o Options) MarshalJSON() ([]byte, error) {
	ret := make([]string, 0, len(o))
	for key := range o {
		ret = append(ret, key)
	}
	return json.Marshal(ret)
}

func (qp *QueryParameters) String() string {
	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	if err := e.Encode(qp); err != nil {
		panic(err)
	}
	return b.String()
}

func (qp *QueryParameters) PrettyString() string {

	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	if err := e.Encode(qp); err != nil {
		panic(err)
	}
	return b.String()

}
