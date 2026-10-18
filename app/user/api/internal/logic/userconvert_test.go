package logic

import (
	"testing"

	"go-zero-demo/user-rpc/pb/user"
)

func TestToUserInfoMapsRole(t *testing.T) {
	info := toUserInfo(&user.UserInfo{Id: 7, Username: "alice", Role: 1, CreatedAt: 100, UpdatedAt: 200})
	if info.Id != 7 || info.Username != "alice" || info.Role != 1 {
		t.Fatalf("unexpected mapping: %+v", info)
	}
	if info.CreatedAt != 100 || info.UpdatedAt != 200 {
		t.Fatalf("unexpected timestamps: %+v", info)
	}
}

func TestToAddressInfo(t *testing.T) {
	info := toAddressInfo(&user.AddressInfo{Id: 3, UserId: 7, ReceiverName: "张三", IsDefault: 1})
	if info.Id != 3 || info.UserId != 7 || info.ReceiverName != "张三" || info.IsDefault != 1 {
		t.Fatalf("unexpected mapping: %+v", info)
	}
}
