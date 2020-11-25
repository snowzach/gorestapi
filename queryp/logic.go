package queryp

import "fmt"

type FilterLogic int

const (
	FilterLogicStart FilterLogic = iota
	FilterLogicAnd
	FilterLogicOr

	FilterLogicSymAnd = "&"
	FilterLogicSymOr  = "|"
)

var (
	FilterLogicSymToFilterLogic = map[string]FilterLogic{
		FilterLogicSymAnd: FilterLogicAnd,
		FilterLogicSymOr:  FilterLogicOr,
	}
)

func (logic FilterLogic) MarshalJSON() ([]byte, error) {
	return []byte(`"` + logic.String() + `"`), nil
}

func (logic FilterLogic) String() string {
	switch logic {
	case FilterLogicStart:
		return "START"
	case FilterLogicAnd:
		return "AND"
	case FilterLogicOr:
		return "OR"
	default:
		return fmt.Sprintf(`LOGIC(%d)`, logic)
	}
}
