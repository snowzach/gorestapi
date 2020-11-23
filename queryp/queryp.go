package queryp

import (
	"encoding/json"
	"fmt"
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
	Logic FilterLogic
	Op    FilterOp
	Field Field
	Value interface{}

	SubFilter Filter
}

type FilterOp int
type FilterLogic int

const (
	FilterTypeNotFound FilterType = iota
	FilterTypeEquals
	FilterTypeString
	FilterTypeTime
	FilterTypeBool
	FilterTypeCompare
)

const (
	FilterOpEquals FilterOp = iota
	FilterOpNotEquals
	FilterOpLessThan
	FilterOpLessThanEqual
	FilterOpGreaterThan
	FilterOpGreaterThanEqual
	FilterOpILike
	FilterOpNotILike
)

const (
	FilterLogicAnd FilterLogic = iota
	FilterLogicOr
)

type Field = string // Alias
type FilterType int
type FilterFieldTypes map[Field]FilterType

type Sort []SortField

type SortField struct {
	Field Field
	Desc  bool
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

func (op FilterOp) MarshalJSON() ([]byte, error) {
	switch op {
	case FilterOpEquals:
		return []byte(`"="`), nil
	case FilterOpNotEquals:
		return []byte(`"!="`), nil
	case FilterOpLessThan:
		return []byte(`"<"`), nil
	case FilterOpLessThanEqual:
		return []byte(`"<="`), nil
	case FilterOpGreaterThan:
		return []byte(`">"`), nil
	case FilterOpGreaterThanEqual:
		return []byte(`">="`), nil
	case FilterOpILike:
		return []byte(`"=~"`), nil
	case FilterOpNotILike:
		return []byte(`"!=~"`), nil
	default:
		return []byte(fmt.Sprintf(`"?!:%d"`, op)), nil
	}
}

func (logic FilterLogic) MarshalJSON() ([]byte, error) {
	switch logic {
	case FilterLogicAnd:
		return []byte(`"AND"`), nil
	case FilterLogicOr:
		return []byte(`"OR"`), nil
	default:
		return []byte(fmt.Sprintf(`"?!:%d"`, logic)), nil
	}
}

func (qp *QueryParameters) String() string {
	b, _ := json.Marshal(qp)
	return string(b)
}

func (qp *QueryParameters) PrettyString() string {
	b, _ := json.MarshalIndent(qp, "", "  ")
	return string(b)
}
