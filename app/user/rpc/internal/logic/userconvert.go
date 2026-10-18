package logic

import (
	"go-zero-demo/user-rpc/internal/model"
	"go-zero-demo/user-rpc/pb/user"
)

func toUserInfo(m *model.User) *user.UserInfo {
	if m == nil {
		return nil
	}
	return &user.UserInfo{
		Id:        m.Id,
		Username:  m.Username,
		Nickname:  m.Nickname,
		Email:     m.Email,
		Phone:     m.Phone,
		Avatar:    m.Avatar,
		Status:    int32(m.Status),
		Remark:    m.Remark,
		CreatedAt: m.CreatedAt.Unix(),
		UpdatedAt: m.UpdatedAt.Unix(),
		Role:      int32(m.Role),
	}
}

func toAddressInfo(m *model.UserAddress) *user.AddressInfo {
	if m == nil {
		return nil
	}
	return &user.AddressInfo{
		Id:            int64(m.Id),
		UserId:        m.UserId,
		ReceiverName:  m.ReceiverName,
		ReceiverPhone: m.ReceiverPhone,
		Province:      m.Province,
		City:          m.City,
		District:      m.District,
		DetailAddress: m.DetailAddress,
		PostalCode:    m.PostalCode,
		Tag:           m.Tag,
		IsDefault:     int32(m.IsDefault),
		CreatedAt:     m.CreatedAt.Unix(),
		UpdatedAt:     m.UpdatedAt.Unix(),
	}
}
