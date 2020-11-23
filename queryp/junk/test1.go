package main

import (
	"fmt"
	"log"

	"github.com/snowzach/gorestapi/queryp"
)

func main() {

	q, err := queryp.ParseQuery("field=value&((another=<value|yet=another1|limit=weee))|third=value&limit=10&option=beans&sort=test,-another")
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}

	// fmt.Println(len(q[0], len))

	fmt.Println(q.PrettyString())

}
