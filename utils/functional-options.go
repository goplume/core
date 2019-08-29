package utils

import _ "github.com/cweill/gotests"

type Optional interface {
	IsSet() bool
}

type Bool struct {
	value *bool
}

func (this *Bool) IsSet() bool {
	return this.value != nil
}

func (this *Bool) SetTrue() {
	b := true
	this.value = &b
}

func (this *Bool) Set(v bool) {
	this.value = &v
}

func (this *Bool) Get() bool {
	return *this.value
}
