package xredis

import (
	"crypto/tls"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisaside"
	"time"
)

func (c *Config) NewCache() (rueidisaside.CacheAsideClient, error) {
	var tlsConfig *tls.Config
	if c.TLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client, err := rueidisaside.NewClient(rueidisaside.ClientOption{
		ClientBuilder: nil,
		ClientOption: rueidis.ClientOption{
			TLSConfig:      tlsConfig,
			SendToReplicas: nil,
			Username:       c.Username,
			Password:       c.Password,
			ClientName:     c.ClientName,
			InitAddress:    c.Addrs,
			SelectDB:       c.DB,
		},
		ClientTTL: time.Duration(c.TTL) * time.Millisecond,
	})
	return client, err
}
