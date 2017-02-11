package gconfig

import "testing"

func TestNewGconfig(t *testing.T) {
	conf := NewGconfig()
	t.Log(conf.Username)
}
