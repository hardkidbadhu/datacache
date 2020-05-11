package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTotalRecordsInPage(t *testing.T) {

	c := TotalRecordsInPage(2, 10, 20)
	assert.Equal(t, 10, c)

	k := TotalRecordsInPage(2, 10, 6)
	assert.Equal(t, 0, k)

	j := TotalRecordsInPage(2, 10, 36)
	assert.Equal(t, 10, j)

	m := TotalRecordsInPage(4, 10, 36)
	assert.Equal(t, 6, m)

	o := TotalRecordsInPage(5, 10, 46)
	assert.Equal(t, 6, o)

	l := TotalRecordsInPage(3, 2, 4)
	assert.Equal(t, 0, l)

}
