package util

import (
	"net/http"
	"reflect"

	"github.com/techlabs/swabhav/tsam/errors"

	uuid "github.com/satori/go.uuid"
)

// MakeMapID Test
func MakeMapID(cntx interface{}) (map[uuid.UUID]bool, error) {
	elems, err := interfaceSlice(cntx)
	if err != nil {
		return nil, err
	}
	mp := make(map[uuid.UUID]bool)
	for _, e := range elems {
		out := reflect.ValueOf(e).Elem()
		id := out.FieldByName("ID").Interface().(uuid.UUID)
		mp[id] = true
	}
	return mp, nil
}

func interfaceSlice(param interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(param)
	if s.Kind() != reflect.Slice {
		return nil, errors.NewHTTPError("given parameter is non-slice type", http.StatusInternalServerError)
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}
