package queryp

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

// Handles parsing query requests with complex matching and precedence
var (
	stringQueryParser = regexp.MustCompile("(&|^|\\|)(\\(*)([a-zA-Z_.]+)(!=~~|=~~|!=~|=~|!:~|!:|:~|:|<=|=<|>=|=>|!=|<|>|=)([^&|\\)|$]+)(\\)*)")

	ErrCouldNotParse = errors.New("could not parse")
)

func ParseRawQuery(rq string) (*QueryParameters, error) {
	q, err := url.PathUnescape(rq)
	if err != nil {
		return nil, err
	}
	return ParseQuery(q)
}

// ParseQuery converts a string into query parameters
// This loosely follows standard HTTP URL encoding
func ParseQuery(q string) (*QueryParameters, error) {

	qp := &QueryParameters{
		Sort:    make(Sort, 0),
		Options: make(Options),
	}

	if len(q) == 0 {
		return qp, nil
	}

	matches := stringQueryParser.FindAllStringSubmatch(q, -1)
	if len(matches) == 0 {
		return nil, ErrCouldNotParse
	}

	// for i, match := range matches {
	// 	fmt.Printf("%d: %v\n", i, match)
	// }

	// Recursive parse function
	var parsedChars int
	var pos int
	var found bool
	var err error
	var parse func(depth int) (Filter, error)
	parse = func(depth int) (Filter, error) {
		start := true

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

			var logic FilterLogic // Default is START if omitted
			if start {
				logic = FilterLogicStart
				start = false
			} else {
				logic, found = FilterLogicSymToFilterLogic[m[1]]
				if !found {
					return nil, fmt.Errorf("invalid filter logic: %s", m[1])
				}
			}

			op, found := FilterOpSymToFilterOp[m[4]]
			if !found {
				return nil, fmt.Errorf("invalid filter logic: %s", m[4])
			}

			// If we have a paren we haven't traversed down into
			if len(m[2]) > depth {
				subFilter, err := parse(depth + 1 + len(m[2]) - len(m[6]))
				if err != nil {
					return nil, err
				}
				filter = append(filter, FilterTerm{
					Logic:     logic,
					SubFilter: subFilter, // Parse, handle redundant parens
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
								qp.Sort = append(qp.Sort, SortTerm{Field: sortField[1:], Desc: true})
							} else {
								qp.Sort = append(qp.Sort, SortTerm{Field: sortField})
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
				return filter, nil
			}
			pos++
		}

		return filter, nil
	}

	// Parse the filter
	qp.Filter, err = parse(0)
	if err != nil {
		return nil, err
	}

	// If the entire string was not parsed, return error
	if parsedChars != len(q) {
		return nil, ErrCouldNotParse
	}

	return qp, nil

}
