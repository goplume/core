package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidMsisdn_EmptyFailed(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMsisdn("")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMsisdn_InvalidLenght(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMsisdn("12")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMsisdn_Invalid_format(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMsisdn("+77071234567")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMsisdn_Invalid_char_format(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMsisdn("sdfgs")
	assert.NotNil(err)
	assert.False(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}

func Test_ValidMsisdn_IsValid(t *testing.T) {
	assert := assert.New(t)
	valid, err := IsValidMsisdn("77071234567")
	assert.Nil(err)
	assert.True(valid)
	log.Info(fmt.Sprintf("valid %t, err %+v", valid, err))
}
