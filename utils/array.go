package utils

import (
	"reflect"
	"time"
)

// ContainsForTimeAndCheckDuration function is use to pass array of time.Time and check time is already exist or not in between greater or lesser one minute
func ContainsForTimeAndCheckDuration(array []time.Time, element time.Time) bool {

	isResult := false

	for _, time := range array {
		if time == element {
			isResult = true
			return isResult
		}

		duration := element.Sub(time)
		if !(duration.Minutes() > 15.0) && !(duration.Minutes() < -15.0) {
			isResult = true
			return isResult
		}
	}
	return isResult
}

// Contains function is use to check the slice given value is contains or not
func Contains(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {

			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}

	return false
}
