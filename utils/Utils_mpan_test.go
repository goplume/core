package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidMpan_EmptyFailed(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMpan("")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMpan_InvalidLenght(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMpan("12")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMpan_Invalid_format(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMpan("+77071234567")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMpan_Invalid_char_format(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMpan("sdfgs")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMpan_IsValid(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMpan("4598959951692349")
	assert.Nil(err)
	assert.True(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMpan_Invalid_zerro(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMpan("087008706081524")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}
