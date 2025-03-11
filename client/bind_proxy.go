package client

import (
	"context"
	"encoding/json"
	"reflect"
)

func BindProxy(s service, proxy proxy) error {
	typ := reflect.TypeOf(s)
	val := reflect.ValueOf(s)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fTyp := typ.Field(i)
		fVal := val.Field(i)

		if fVal.CanSet() {
			fn := reflect.MakeFunc(fTyp.Type, func(args []reflect.Value) (results []reflect.Value) {
				ctx := args[0].Interface().(context.Context)
				reqData, err := json.Marshal(args[1].Interface())
				resVal := reflect.New(fTyp.Type.Out(0).Elem())
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}

				req := &Request{
					ServiceName: s.Name(),
					MethodName:  fTyp.Name,
					Args:        reqData,
				}

				var res *Response
				res, err = proxy.invoke(ctx, req)
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}

				err = json.Unmarshal(res.Data, resVal.Interface())
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}

				return []reflect.Value{
					resVal,
					reflect.Zero(reflect.TypeOf(new(error)).Elem()),
				}
			})

			fVal.Set(fn)
		}
	}

	return nil
}
