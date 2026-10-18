package logic

import (
	"context"

	"go-zero-demo/pkg/auth"
)

func injectAdmin(ctx context.Context) (context.Context, int64) {
	adminID, ok := auth.UserID(ctx)
	if !ok || adminID <= 0 {
		return ctx, 0
	}
	return auth.InjectUserID(ctx, adminID), adminID
}
