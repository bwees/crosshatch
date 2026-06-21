package dtos

type Go2RTCStream struct {
	Producers []Go2RTCProducer `json:"producers"`
}

type Go2RTCProducer struct {
	URL string `json:"url"`
}

type Go2RTCInfo struct {
	Version    string `json:"version"`
	ConfigPath string `json:"config_path"`
}
