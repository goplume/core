package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsBlankString(t *testing.T) {
	assert.True(t, StringIsBlank(""))
	assert.True(t, StringIsBlank(" "))
	assert.True(t, StringIsBlank("\t"))
	assert.True(t, StringIsBlank("\n"))
	assert.True(t, StringIsBlank("\r"))
	assert.True(t, StringIsBlank("\n\r"))
	assert.False(t, StringIsBlank("   s  "))
}

func Test_WhiteSpace(t *testing.T) {
	assert.False(t, IsWhiteSpaceChar('\\'))
	assert.True(t, IsWhiteSpaceChar(' '))
	assert.True(t, IsWhiteSpaceChar('\t'))
	assert.True(t, IsWhiteSpaceChar('\n'))
	assert.True(t, IsWhiteSpaceChar('\r'))
}
