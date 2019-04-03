package fetcher

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	// так как знаем внутреннее устройство, можем игнорировать ограничение
	// на синглтон
	f := newFetcher()

	val, err := f.Get()
	if err != nil {
		t.Error(err)
	}

	const expected = "xxx"
	if val != expected {
		t.Errorf("%q != %q", val, expected)
	}
}

func TestList(t *testing.T) {
	f := newFetcher()

	val, err := f.List()
	if err != nil {
		t.Error(err)
	}

	expected := []string{"xxx", "yyy", "zzz"}
	if !reflect.DeepEqual(val, expected) {
		t.Errorf("%v != %v", val, expected)
	}
}
