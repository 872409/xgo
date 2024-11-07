package xredis

import (
	"crypto/tls"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	IsCluster  bool `json:",default=true"`
	Addrs      []string
	ClientName string
	Username   string
	Password   string
	DB         int
	MaxRetries int
	Prefix     string
	TLS        bool `json:",default=true"`
	TTL        int  `json:",default=10"`
}

// 单机
func (c *Config) client() *redis.Client {
	var tlsConfig *tls.Config
	if c.TLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return redis.NewClient(&redis.Options{
		Addr:       c.Addrs[0],
		ClientName: c.ClientName,
		Username:   c.Username,
		Password:   c.Password,
		MaxRetries: c.MaxRetries,
		TLSConfig:  tlsConfig,
	})
}

// 集群
func (c *Config) newClusterClient() *redis.ClusterClient {
	var tlsConfig *tls.Config
	if c.TLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:      c.Addrs,
		ClientName: c.ClientName,
		Password:   c.Password,
		Username:   c.Username,
		MaxRetries: c.MaxRetries,
		TLSConfig:  tlsConfig,
	})
}

func (c *Config) NewRedis() redis.UniversalClient {
	if c.IsCluster {
		return c.newClusterClient()
	}
	return c.client()
}
