package config

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	data, err := ioutil.ReadFile("apiservice_conf.json")
	assert.Nil(err)
	assert.NotEmpty(data)

	var expected Configuration
	err = json.Unmarshal(data, &expected)
	assert.Nil(err)

	got := Parse("apiservice_conf.json")
	assert.IsType(&Configuration{}, got)
	assert.NotNil(got)
	assert.Equal(&expected, got)
}
