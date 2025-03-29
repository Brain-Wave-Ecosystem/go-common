package clients

import "github.com/hashicorp/consul/api"

func NewConsulClient(addr string) (*api.Client, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr

	return api.NewClient(cfg)
}
