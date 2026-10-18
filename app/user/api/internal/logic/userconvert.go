package logic

import (
	"go-zero-demo/user-api/internal/types"
	"go-zero-demo/user-rpc/pb/user"
)

func toUserInfo(info *user.UserInfo) types.UserInfo {
	return types.UserInfo{
		Id:        info.Id,
		Username:  info.Username,
		Nickname:  info.Nickname,
		Email:     info.Email,
		Phone:     info.Phone,
		Avatar:    info.Avatar,
		Status:    info.Status,
		Remark:    info.Remark,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
		Role:      info.Role,
	}
}

func toAddressInfo(info *user.AddressInfo) types.AddressInfo {
	return types.AddressInfo{
		Id:            info.Id,
		UserId:        info.UserId,
		ReceiverName:  info.ReceiverName,
		ReceiverPhone: info.ReceiverPhone,
		Province:      info.Province,
		City:          info.City,
		District:      info.District,
		DetailAddress: info.DetailAddress,
		PostalCode:    info.PostalCode,
		Tag:           info.Tag,
		IsDefault:     info.IsDefault,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}
}
