package queryp

import (
	"errors"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

// Handles parsing query requests with complex matching and precedence
var stringQueryParser = regexp.MustCompile("(&|^|\\|)(\\(*)([a-zA-Z_.]+)(<=|=<|>=|=>|=~|!=~|!=|<|>|=)([^&|\\)|$]+)(\\)*)")

var (
	ErrCouldNotParse = errors.New("could not parse")
)

// ParseQuery converts a string into query parameters
// This loosely follows standard HTTP URL encoding
func ParseQuery(q string) (*QueryParameters, error) {

	matches := stringQueryParser.FindAllStringSubmatch(q, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	// for i, match := range matches {
	// 	fmt.Printf("%d: %v\n", i, match)
	// }

	qp := &QueryParameters{
		Sort:    make(Sort, 0),
		Options: make(Options),
	}

	// Recursive parse function
	var parsedChars int
	var pos int
	var parse func(depth int) Filter
	parse = func(depth int) Filter {
		filter := make([]FilterTerm, 0)
		for pos < len(matches) {

			m := matches[pos]
			// m[0] = matches string
			// m[1] = logic
			// m[2] = open parens
			// m[3] = field
			// m[4] = op
			// m[5] = value
			// m[6] = close parens

			var logic FilterLogic // Default is AND if omitted
			switch m[1] {
			case "&":
				logic = FilterLogicAnd
			case "|":
				logic = FilterLogicOr
			}

			var op FilterOp // Default is EQUALS if omitted
			switch m[4] {
			case "<=":
				op = FilterOpLessThanEqual
			case "=<":
				op = FilterOpGreaterThanEqual
			case ">=":
				op = FilterOpGreaterThanEqual
			case "=>":
				op = FilterOpGreaterThanEqual
			case "=~":
				op = FilterOpILike
			case "!=~":
				op = FilterOpNotILike
			case "!=":
				op = FilterOpNotEquals
			case "<":
				op = FilterOpLessThan
			case ">":
				op = FilterOpGreaterThan
			case "=":
				op = FilterOpEquals
			}

			// If we have a paren we haven't traversed down into
			if len(m[2]) > depth {
				filter = append(filter, FilterTerm{
					Logic:     logic,
					SubFilter: parse(depth + 1 + len(m[2]) - len(m[6])), // Parse, handle redundant parens
				})

			} else {

				parsedChars += len(m[0])

				field, value := m[3], m[5]

				// We will only handle options at depth 0
				if depth == 0 {
					switch field {
					case "limit":
						qp.Limit = cast.ToInt(value)
						pos++
						continue
					case "offset":
						qp.Limit = cast.ToInt(value)
						pos++
						continue
					case "sort":
						for _, sortField := range strings.Split(value, ",") {
							if len(sortField) > 1 && sortField[0] == '-' { // Reverse sort
								qp.Sort = append(qp.Sort, SortField{Field: sortField[1:], Desc: true})
							} else {
								qp.Sort = append(qp.Sort, SortField{Field: sortField})
							}
						}
						pos++
						continue
					case "option":
						for _, optionField := range strings.Split(value, ",") {
							qp.Options[optionField] = struct{}{}
						}
						pos++
						continue
					}
				}

				// Otherwise add a filter option
				filter = append(filter, FilterTerm{
					Logic: logic,
					Op:    op,
					Field: Field(field),
					Value: value,
				})
			}
			if len(m[6]) > 0 { // Close Paren
				return filter
			}
			pos++
		}
		return filter
	}

	// Parse the filter
	qp.Filter = parse(0)

	// If the entire string was not parsed, return error
	if parsedChars != len(q) {
		return nil, ErrCouldNotParse
	}

	return qp, nil

}
