package entity

type TrxLog struct {
	Id              int    `json:"id"`
	Username        string `json:"username"`
	Account         string `json:"account"`
	ReferenceNum    string `json:"reference_num"`
	Logged          int    `json:"logged"`
	TrxType         string `json:"trx_type"`
	TrxStatus       string `json:"trx_status"`
	TrxObject       string `json:"trx_object"`
	IpAddressSource string `json:"ip_address_source"`
	Agent           string `json:"agent"`
	TrxDate         string `json:"trx_date"`
	TypeUser        string `json:"type_user"`
}
