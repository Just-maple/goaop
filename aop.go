package goaop

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
)

func Parse(i interface{}, dir string) (err error) {
	v := reflect.ValueOf(i)
	p := &parser{}
	p.walk(context.Background(), v)

	if len(modPath) > 0 {
		err = p.gen(dir)
	}
	return
}

type parser struct {
	proxyInterfaces []proxyInterface
}

type proxyInterface struct {
	Name      string
	Pkg       string
	Parent    string
	ParentPkg string
	Type      reflect.Type
	FieldName string
}

var (
	modPath = func() string {
		ret, _ := exec.Command("go", "list", "-m").Output()
		return strings.TrimSpace(string(ret))
	}()

	importMyself = func() string {
		return reflect.TypeOf(new(InterceptorFactory)).Elem().PkgPath()
	}()

	funcNameReplacer = strings.NewReplacer(".", "_", "/", "_", "-", "_")
)

const tag = `goaop`

func (p *parser) walk(ctx context.Context, v reflect.Value) {
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Interface:
		v = v.Elem()
	case reflect.Struct:
	default:
		return
	}

	v = reflect.Indirect(v)
	if v.IsZero() || !strings.HasPrefix(v.Type().PkgPath(), modPath) {
		return
	}

	pkg := strings.TrimPrefix(funcNameReplacer.Replace(strings.TrimPrefix(v.Type().PkgPath(), modPath)), "_")
	structName := fmt.Sprintf("%s__%s", pkg, v.Type().Name())

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := v.Type().Field(i)
		tagGoAop := ft.Tag.Get(tag)

		if tagGoAop == "-" {
			continue
		}

		if len(ft.PkgPath) > 0 {
			continue
		}

		if fv.Kind() == reflect.Interface && len(tagGoAop) > 0 {
			depT := ft.Type
			name := fmt.Sprintf("%s__%s__%s", structName, ft.Name, funcNameReplacer.Replace(depT.String()))

			if wrappedProxy, ok := WrapInterceptorProxy(tagGoAop, name, fv.Interface()); ok {
				fv.Set(reflect.ValueOf(wrappedProxy).Convert(ft.Type))
			}

			p.proxyInterfaces = append(p.proxyInterfaces, proxyInterface{
				Pkg:       depT.PkgPath(),
				Parent:    structName,
				ParentPkg: pkg,
				Type:      depT,
				FieldName: ft.Name,
				Name:      name,
			})
		}

		if ctx.Value(fv) != nil {
			continue
		}
		nCtx := context.WithValue(ctx, fv, struct{}{})

		p.walk(nCtx, fv)
	}
}
