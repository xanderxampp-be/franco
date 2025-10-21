package entity

import (
	"net/http"
)

type Responselog struct {
	Username          string
	AccountDebet      string
	Amount            int
	AmountFloat       float64
	Fee               int
	TransactionRefnum string
	ThirdParty        string
	ResponseHeader    http.Header
	ResponseBody      interface{}
	TrxType           string
	DeviceId          string
	DeviceType        string
	DeviceName        string
	DeviceVersion     string
	DeviceSequence    string
	Version           string
	ResponseCode      string
	Trace             string
	Timestamp         string
	Elapsed           string
	FastMenuFlag      bool
	TypeUser          string
}
