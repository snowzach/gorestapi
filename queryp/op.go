package queryp

import "fmt"

type FilterOp int

const (
	FilterOpEquals FilterOp = iota
	FilterOpNotEquals

	FilterOpLessThan
	FilterOpLessThanEqual
	FilterOpGreaterThan
	FilterOpGreaterThanEqual

	FilterOpLike
	FilterOpNotLike
	FilterOpILike
	FilterOpNotILike

	FilterOpRegexp
	FilterOpNotRegexp
	FilterOpIRegexp
	FilterOpNotIRegexp

	FilterOpSymEquals      = "="
	FilterOpSymNotEquals   = "!="
	FilterOpSymLessThan    = "<"
	FilterOpSymGreaterThan = ">"

	FilterOpSymLessThanEqual     = "<="
	FilterOpSymLessThanEqual2    = "=<"
	FilterOpSymGreaterThanEqual  = ">="
	FilterOpSymGreaterThanEqual2 = "=>"

	FilterOpSymLike     = "=~"
	FilterOpSymNotLike  = "!=~"
	FilterOpSymILike    = "=~~"
	FilterOpSymNotILike = "!=~~"

	FilterOpSymRegexp     = ":"
	FilterOpSymNotRegexp  = "!:"
	FilterOpSymIRegexp    = ":~"
	FilterOpSymNotIRegexp = "!:~"
)

var (
	FilterOpSymToFilterOp = map[string]FilterOp{
		FilterOpSymEquals:            FilterOpEquals,
		FilterOpSymNotEquals:         FilterOpNotEquals,
		FilterOpSymLessThan:          FilterOpLessThan,
		FilterOpSymGreaterThan:       FilterOpGreaterThan,
		FilterOpSymLessThanEqual:     FilterOpLessThanEqual,
		FilterOpSymLessThanEqual2:    FilterOpLessThanEqual,
		FilterOpSymGreaterThanEqual:  FilterOpGreaterThanEqual,
		FilterOpSymGreaterThanEqual2: FilterOpGreaterThanEqual,
		FilterOpSymLike:              FilterOpLike,
		FilterOpSymNotLike:           FilterOpNotLike,
		FilterOpSymILike:             FilterOpILike,
		FilterOpSymNotILike:          FilterOpNotILike,
		FilterOpSymRegexp:            FilterOpRegexp,
		FilterOpSymNotRegexp:         FilterOpNotRegexp,
		FilterOpSymIRegexp:           FilterOpIRegexp,
		FilterOpSymNotIRegexp:        FilterOpNotIRegexp,
	}
	FilterOpToFilterOpSym = map[FilterOp]string{
		FilterOpEquals:           FilterOpSymEquals,
		FilterOpNotEquals:        FilterOpSymNotEquals,
		FilterOpLessThan:         FilterOpSymLessThan,
		FilterOpGreaterThan:      FilterOpSymGreaterThan,
		FilterOpLessThanEqual:    FilterOpSymLessThanEqual,
		FilterOpGreaterThanEqual: FilterOpSymGreaterThanEqual,
		FilterOpLike:             FilterOpSymLike,
		FilterOpNotLike:          FilterOpSymNotLike,
		FilterOpILike:            FilterOpSymILike,
		FilterOpNotILike:         FilterOpSymNotILike,
		FilterOpRegexp:           FilterOpSymRegexp,
		FilterOpNotRegexp:        FilterOpSymNotRegexp,
		FilterOpIRegexp:          FilterOpSymIRegexp,
		FilterOpNotIRegexp:       FilterOpSymNotIRegexp,
	}
)

func (op FilterOp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + op.String() + `"`), nil
}

func (op FilterOp) String() string {
	if filterOptSym, found := FilterOpToFilterOpSym[op]; found {
		return filterOptSym
	}
	return fmt.Sprintf(`OP(%d):`, op)
}
