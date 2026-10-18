package logic

import (
	"context"

	"go-zero-demo/pkg/auth"
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"
)

func injectUser(ctx context.Context) (context.Context, bool) {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return ctx, false
	}
	return auth.InjectUserID(ctx, userID), true
}

func toSaveAddressReq(req *types.SaveAddressReq) *user.SaveAddressReq {
	return &user.SaveAddressReq{
		Id:            req.Id,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		DetailAddress: req.DetailAddress,
		PostalCode:    req.PostalCode,
		Tag:           req.Tag,
		IsDefault:     req.IsDefault,
	}
}
