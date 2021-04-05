package goaop

import (
	"fmt"
	"reflect"
	"runtime"
)

type JoinPoint struct {
	Method  interface{}
	Params  []interface{}
	Returns []interface{}
}

func (jp *JoinPoint) Defer() {
	if f := runtime.FuncForPC(reflect.ValueOf(jp.Method).Pointer()); f != nil {
		fmt.Printf("[GOAOP] END: %v\n", f.Name())
	}
}

func (_ *JoinPoint) InitInterceptor(point *JoinPoint) (Interceptor, bool) {
	if f := runtime.FuncForPC(reflect.ValueOf(point.Method).Pointer()); f != nil {
		fmt.Printf("[GOAOP] START: %v\n", f.Name())
	}
	return point, true
}
