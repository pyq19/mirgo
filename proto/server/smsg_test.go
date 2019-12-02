package server

import (
	"reflect"
	"testing"
)

func TestConst(t *testing.T) {
	t.Log(CONNECTED)
	t.Log(CLIENT_VERSION)
	t.Log(reflect.TypeOf(CONNECTED))
}
