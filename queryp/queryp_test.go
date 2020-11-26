package queryp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParser(t *testing.T) {

	q, err := ParseQuery("field=value&((another=<value|yet=another1|limit=weee))|third=value&limit=10&option=beans&sort=test,-another")
	assert.Nil(t, err)
	assert.Equal(t, `{"filter":[{"logic":"START","op":"=","field":"field","value":"value"},{"logic":"AND","op":"=","field":"","value":"","sub_filter":[{"logic":"START","op":"<=","field":"another","value":"value"},{"logic":"OR","op":"=","field":"yet","value":"another1"},{"logic":"OR","op":"=","field":"limit","value":"weee"}]},{"logic":"OR","op":"=","field":"third","value":"value"}],"sort":[{"field":"test","desc":false},{"field":"another","desc":true}],"options":["beans"],"limit":10,"offset":0}`+"\n", q.String())

}
