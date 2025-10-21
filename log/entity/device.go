package entity

type Device struct {
	DeviceID       string `json:"device_id"`
	DeviceVersion  string `json:"device_version"`
	DeviceName     string `json:"device_name"`
	DeviceType     string `json:"device_type"`
	DeviceSequence string `json:"device_sequence"`
	Version        string `json:"version"`
}
