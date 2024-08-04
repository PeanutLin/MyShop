package common

import (
	"context"
	"productshop/product_shop/middleware/logs"
)

func DefaultDeferFunc(ctx context.Context) {
	if e := recover(); e != nil {
		logs.Error("goroutine panic", logs.String("error", e.(string)))
		// metrics
	}
}

func GoWithRecovery(ctx context.Context, df func(ctx context.Context), f func()) {
	go func() {
		defer df(ctx)
		f()
	}()
}

func GoWithDefaultRecovery(ctx context.Context, f func()) {
	GoWithRecovery(ctx, DefaultDeferFunc, f)
}
