package goaop

import "sync"

type InterceptorFactory interface {
	InitInterceptor(point *JoinPoint) (Interceptor, bool)
}

type Interceptor interface {
	Defer()
}

type InterceptorProxy interface {
	RegisterInterceptor(origin interface{}, inceptorFactory InterceptorFactory) interface{}
}

var (
	interceptorProxyContainer   = map[string]InterceptorProxy{}
	interceptorFactoryContainer = map[string]InterceptorFactory{}
	lock                        sync.Mutex
)

func RegisterInterceptorProxy(name string, i InterceptorProxy) {
	lock.Lock()
	defer lock.Unlock()
	interceptorProxyContainer[name] = i
}

func RegisterInterceptorFactory(tag string, i InterceptorFactory) {
	lock.Lock()
	defer lock.Unlock()
	interceptorFactoryContainer[tag] = i
}

func WrapInterceptorProxy(tag, name string, origin interface{}) (wrappedProxy interface{}, ok bool) {
	if _, wrapped := origin.(InterceptorProxy); wrapped {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	// get factory and proxy from container
	factory, ok := interceptorFactoryContainer[tag]
	if !ok {
		return
	}
	proxy, ok := interceptorProxyContainer[name]
	if !ok {
		return
	}
	// new proxy by factory
	return proxy.RegisterInterceptor(origin, factory), true
}
