package store

import (
	"errors"
	"net/url"
	"strings"

	"github.com/spf13/cast"
)

// ErrNotFound is a standard not found error
var ErrNotFound = errors.New("not found")

type Results struct {
	Count   int64       `json:"count"`
	Results interface{} `json:"results"`
}

const (
	FilterTypeEquals   FilterType = "="
	FilterTypeLike     FilterType = "LIKE"
	FilterTypeILike    FilterType = "ILIKE"
	FilterTypeCompare  FilterType = "COMPARE"
	FilterTypeNotFound FilterType = ""
)

type FilterType string
type FilterValues map[string][]string       // Define field filter values
type FilterFieldTypes map[string]FilterType // Define field to filter type

type SortFields []string
type Sort struct {
	Field string
	Desc  bool
}
type SortValues []*Sort
type Options map[string]struct{} // Just a lookup of string

func MakeOptions(options []string) Options {
	o := make(Options)
	for _, value := range options {
		o[value] = struct{}{}
	}
	return o
}

func (o Options) HasOption(option string) bool {
	_, found := o[option]
	return found
}

type FindQueryParameters struct {
	PreFilter          FilterValues `json:"pre_filter"`
	PreFilterInclusive Options      `json:"pre_filter_inclusive"`
	Filter             FilterValues `json:"filter"`
	FilterInclusive    Options      `json:"filter_and"`
	Sort               SortValues   `json:"sort"`
	Options            Options      `json:"options"`
	Limit              int          `json:"limit"`
	Offset             int          `json:"offset"`
}

// ParseURLValuesToFindQueryParameters handles turning query paramters into filter values
func ParseURLValuesToFindQueryParameters(values url.Values) *FindQueryParameters {

	fqp := &FindQueryParameters{
		PreFilter:          make(FilterValues),
		PreFilterInclusive: make(Options),
		Filter:             make(FilterValues),
		FilterInclusive:    make(Options),
		Sort:               make(SortValues, 0),
		Options:            make(Options),
	}

	for key, value := range values {
		switch key {
		case "limit":
			fqp.Limit = cast.ToInt(value[0])
		case "offset":
			fqp.Offset = cast.ToInt(value[0])
		case "sort":
			if len(values) == 0 {
				continue
			}
			for _, sortField := range strings.Split(value[0], ",") {
				if index := strings.Index(sortField, ":"); index > -1 {
					if index < len(sortField)-1 && sortField[index+1] == 'd' {
						fqp.Sort = append(fqp.Sort, &Sort{Field: sortField[:index], Desc: true})
					} else {
						fqp.Sort = append(fqp.Sort, &Sort{Field: sortField[:index]})
					}
				} else {
					fqp.Sort = append(fqp.Sort, &Sort{Field: sortField})
				}
			}
		case "inclusive":
			fqp.FilterInclusive = MakeOptions(value)
		case "option":
			fqp.Options = MakeOptions(value)
		default:
			fqp.Filter[key] = append(fqp.Filter[key], value...)
		}
	}
	return fqp

}

// Replace will replace the filter with a new filter
func (fv FilterValues) Replace(field string, values ...string) {
	fv[field] = values
}

// Add will add a new filter to the filter values
func (fv FilterValues) Add(field string, values ...string) {
	if currentValues, found := fv[field]; found {
		fv[field] = append(currentValues, values...)
	} else {
		fv[field] = values
	}
}

func (fft FilterFieldTypes) FindFilterType(search string) (string, FilterType) {

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
