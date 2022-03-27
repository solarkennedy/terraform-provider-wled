package wled_client

type WLEDClient struct {
	host string
}

func NewWLEDClient(host string) *WLEDClient {
	return &WLEDClient{
		host: host,
	}
}
