package postgres

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/snowzach/gorestapi/store"
)

// Allow comparing numbers (or dates)
var filterTypeCompareMatch = regexp.MustCompile("^(=|!=|<|<=|>|>=|)([0-9].*|true|false)$")

// FilterQuery will update the queryClause and queryParams with filter values
func (c *Client) FilterQuery(filterFields store.FilterFieldTypes, filters store.FilterValues, filtersInclusive store.Options, queryClause *strings.Builder, queryParams *[]interface{}) error {

	for baseFilter, values := range filters {

		// Lookup the filter
		filter, filterType := filterFields.FindFilterType(baseFilter)
		if filterType == store.FilterTypeNotFound {
			return fmt.Errorf("could not find filter: %s", baseFilter)
		}

		var clauses []string
		for _, value := range values {
			if value == "" {
				return fmt.Errorf("invalid value for %s", filter)
			}
			var op string
			switch filterType {
			case store.FilterTypeCompare:
				matches := filterTypeCompareMatch.FindStringSubmatch(value)
				if len(matches) < 3 {
					return fmt.Errorf("invalid value for %s", filter)
				}
				// The operator
				op = matches[1]
				if len(op) == 0 {
					op = "="
				}
				value = matches[2]
			default:
				if strings.HasPrefix(value, "!=") {
					op = "NOT " + string(filterType)
					value = value[2:]
					if value == "" {
						return fmt.Errorf("invalid value for %s", filter)
					}
				} else {
					op = string(filterType)
				}
			}
			// Build the query
			*queryParams = append(*queryParams, value)
			clauses = append(clauses, filter+" "+op+" $"+strconv.Itoa(len(*queryParams)))

		}

		queryClause.WriteString(" AND (")
		if filtersInclusive.HasOption(baseFilter) || filtersInclusive.HasOption(filter) {
			queryClause.WriteString(strings.Join(clauses, " AND "))
		} else {
			queryClause.WriteString(strings.Join(clauses, " OR "))
		}
		queryClause.WriteString(")")
	}

	return nil

}

// SortQuery will update the queryClause and queryParams with sort values
func (c *Client) SortQuery(sortFields store.SortFields, sortValues store.SortValues, queryClause *strings.Builder, queryParams *[]interface{}) error {

	if len(sortValues) == 0 {
		return nil
	}

	var index int
	for _, sortValue := range sortValues {
		// Search for exact match
		var found bool
		var sortField string
		for _, sortField = range sortFields {
			if sortValue.Field == sortField {
				found = true
				break
			}
		}
		// Check for a matching suffix
		if !found {
			sortValueFieldSuffix := "." + sortValue.Field
			for _, sortField = range sortFields {
				if strings.HasSuffix(sortField, sortValueFieldSuffix) {
					found = true
					break
				}
			}
		}
		// Found a field, build the order by clause
		if found {
			if index == 0 {
				queryClause.WriteString(" ORDER BY ")
			} else {
				queryClause.WriteString(", ")
			}
			queryClause.WriteString(sortField)
			if sortValue.Desc {
				queryClause.WriteString(" DESC")
			}
			index++
		}
	}
	return nil

}
