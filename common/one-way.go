package common

import "context"

var OneWay = struct {
}{}

func SetOneWay(ctx context.Context) context.Context {
	return context.WithValue(ctx, OneWay, true)
}

func HasOneWay(ctx context.Context) bool {
	val := ctx.Value(OneWay)
	return val == true
}
