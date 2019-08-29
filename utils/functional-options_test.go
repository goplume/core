package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetBool(t *testing.T) {
	var B Bool
	B.Set(EvaluateBool())
	get := B.Get()
	if get == true {
		assert.True(t, get)
	} else {
		assert.False(t, get)
	}
}

func EvaluateBool() bool {
	return 1 == 1
}
