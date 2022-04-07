package postgres

import (
	"strconv"
	"strings"
)

type Field struct {
	Name   string
	Insert string
	Update string
	Arg    interface{}
}

// Builds the values needed to compose an upsert statement
func ComposeUpsert(fields []Field) (string, string, string, []interface{}) {

	names := make([]string, 0)
	inserts := make([]string, 0)
	updates := make([]string, 0)
	args := make([]interface{}, 0)

	for _, field := range fields {
		index := "$#"
		if field.Arg != nil {
			args = append(args, field.Arg)
			index = "$" + strconv.Itoa(len(args))
		}
		if field.Insert != "" {
			names = append(names, field.Name)
			inserts = append(inserts, strings.ReplaceAll(field.Insert, "$#", index))
		}
		if field.Update != "" {
			updates = append(updates, field.Name+" = "+strings.ReplaceAll(field.Update, "$#", index))
		}
	}

	return strings.Join(names, ","), strings.Join(inserts, ","), strings.Join(updates, ","), args

}
