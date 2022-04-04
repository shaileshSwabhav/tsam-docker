package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type date struct {
	Date string
}

// ParseDateFormat to update existing field with tag date from RFC3339 to Normal Date Format
func ParseDateFormat(element interface{}) {
	structType := reflect.TypeOf(element).Elem()
	if reflect.TypeOf(element).Elem().Kind() == reflect.Slice {
		slice := reflect.ValueOf(element).Elem()
		for index := 0; index < slice.Len(); index++ {
			ParseDateFormat(slice.Index(index).Addr().Interface())
		}
		return
	}
	if structType.Kind() == reflect.Struct {
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			fieldVal := reflect.ValueOf(element).Elem().Field(i)
			if field.Type.Kind() == reflect.Slice {
				for ind := 0; ind < fieldVal.Len(); ind++ {
					ParseDateFormat(fieldVal.Index(ind).Addr().Interface())
				}
				continue
			}

			if field.Type.Kind() == reflect.Struct {
				ParseDateFormat(fieldVal.Addr().Interface())
				continue
			}
			fmt.Println(field.Name + " " + field.Tag.Get("gorm"))
			if field.Tag.Get("gorm") == "type:date" {
				out := convertor(fieldVal)
				fieldVal.Set(out)
			}
		}
	}
}

func convertor(el reflect.Value) reflect.Value {
	out := el.Interface()
	fmt.Println(reflect.ValueOf(el))
	str, ok := out.(string)
	if !ok {
		fmt.Println("Error occured")
	}
	st := strings.Split(str, "T")
	inter := date{Date: st[0]}
	return reflect.ValueOf(&inter).Elem().Field(0)
}

// GetCurrentDateString Return string of Current date in Given Format
func GetCurrentDateString(format string) *string {
	date := time.Now().Format(format)
	strdate := string(date)
	return &strdate
}

// AddDateToMonth will add the date 01 to the month type.
func AddDateToMonth(month *string) {
	if month != nil {
		*month = *month + "-01"
	}
}

func GetBeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func GetEndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}

func GetBeginningOfWeek(date time.Time) (prevMonday time.Time) {
	prevMonday = time.Now()

	offset := int(time.Monday - date.Weekday())
	if offset > 0 {
		offset = -6
	}

	return prevMonday.AddDate(0, 0, offset)
}

func GetEndOfWeek(date time.Time) (prevMonday time.Time) {
	prevMonday = time.Now()

	offset := int(time.Monday - date.Weekday())
	if offset > 0 {
		offset = -6
	}

	return prevMonday.AddDate(0, 0, offset+7)
}
