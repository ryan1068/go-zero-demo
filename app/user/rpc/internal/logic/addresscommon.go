package logic

import (
	"strings"
	"unicode/utf8"

	"go-zero-demo/user-rpc/pb/user"
)

const (
	maxAddressCount  = 20
	maxReceiverName  = 64
	maxReceiverPhone = 32
	maxRegionName    = 64
	maxDetailAddress = 255
	maxPostalCode    = 16
	maxAddressTag    = 16
)

type addressParams struct {
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	DetailAddress string
	PostalCode    string
	Tag           string
	IsDefault     int32
}

func normalizeAddress(in *user.SaveAddressReq) (addressParams, bool) {
	p := addressParams{
		ReceiverName:  strings.TrimSpace(in.ReceiverName),
		ReceiverPhone: strings.TrimSpace(in.ReceiverPhone),
		Province:      strings.TrimSpace(in.Province),
		City:          strings.TrimSpace(in.City),
		District:      strings.TrimSpace(in.District),
		DetailAddress: strings.TrimSpace(in.DetailAddress),
		PostalCode:    strings.TrimSpace(in.PostalCode),
		Tag:           strings.TrimSpace(in.Tag),
		IsDefault:     in.IsDefault,
	}
	if p.ReceiverName == "" || utf8.RuneCountInString(p.ReceiverName) > maxReceiverName ||
		p.ReceiverPhone == "" || utf8.RuneCountInString(p.ReceiverPhone) > maxReceiverPhone ||
		utf8.RuneCountInString(p.Province) > maxRegionName ||
		utf8.RuneCountInString(p.City) > maxRegionName ||
		utf8.RuneCountInString(p.District) > maxRegionName ||
		p.DetailAddress == "" || utf8.RuneCountInString(p.DetailAddress) > maxDetailAddress ||
		utf8.RuneCountInString(p.PostalCode) > maxPostalCode ||
		utf8.RuneCountInString(p.Tag) > maxAddressTag ||
		(p.IsDefault != 0 && p.IsDefault != 1) {
		return p, false
	}
	return p, true
}
