package xredis

import (
	"crypto/tls"
	"strconv"
	"strings"

	"github.com/DreamvatLab/go/xerr"
	"github.com/DreamvatLab/go/xlog"
	"github.com/redis/go-redis/v9"
)

// RedisConfig represents the configuration for Redis connection
// It supports both single node and cluster mode
type RedisConfig struct {
	Addrs    []string    // List of Redis node addresses, multiple addresses for cluster mode
	Username string      // Redis username for authentication (optional)
	Password string      // Redis password for authentication (optional)
	DB       int         // Redis database number (0-15)
	TLS      *tls.Config // TLS configuration for Redis connection (optional)
}

// ParseRedisConfig parses a Redis connection string into a RedisConfig struct
// Supported formats:
// Single node: redis://[username:password@]host:port[/db]
// Cluster mode: redis://[username:password@]host1:port1,host2:port2[/db]
func ParseRedisConfig(connStr string) (*RedisConfig, error) {
	if len(connStr) == 0 {
		return nil, xerr.New("empty connection string")
	}

	config := &RedisConfig{
		Addrs: make([]string, 0),
		DB:    0, // Default to database 0
	}

	// Parse connection string in format: redis://username:password@host:port/db
	// Example: redis://:Famous901@192.168.188.166:6379
	// Example cluster: redis://user:pass@192.168.188.166:6379,192.168.188.167:6379/1

	// Remove redis:// prefix if exists
	if len(connStr) > 8 && connStr[:8] == "redis://" {
		connStr = connStr[8:]
	}

	// Split authentication info and host info by '@'
	parts := strings.Split(connStr, "@")
	if len(parts) == 2 {
		// Process authentication info (username:password format)
		auth := parts[0]
		if len(auth) > 0 {
			// Parse username:password format
			if strings.Contains(auth, ":") {
				authParts := strings.Split(auth, ":")
				if len(authParts) != 2 {
					return nil, xerr.Errorf("invalid authentication format: %s", auth)
				}
				config.Username = authParts[0]
				config.Password = authParts[1]
			} else {
				// Invalid authentication format - must contain ':'
				return nil, xerr.Errorf("invalid authentication format: %s", auth)
			}
		}
		connStr = parts[1]
	}

	// Split host:port and database number by '/'
	parts = strings.Split(connStr, "/")
	hostPorts := parts[0]
	if len(hostPorts) == 0 {
		return nil, xerr.New("missing host and port")
	}

	// Parse database number if exists
	if len(parts) > 1 {
		dbStr := parts[1]
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			return nil, xerr.Errorf("invalid database number: %s", dbStr)
		}
		config.DB = db
	}

	// Split multiple addresses for cluster mode
	addresses := strings.Split(hostPorts, ",")
	for _, addr := range addresses {
		if addr == "" {
			return nil, xerr.New("empty address in cluster configuration")
		}
		config.Addrs = append(config.Addrs, addr)
	}

	if len(config.Addrs) == 0 {
		return nil, xerr.New("no valid addresses found")
	}

	return config, nil
}

func NewClient(config *RedisConfig) redis.UniversalClient {
	addrCount := len(config.Addrs)
	if addrCount == 0 {
		xlog.Fatal("addrs cannot be empty")
		return nil
	} else if addrCount == 1 {
		c := &redis.Options{
			Addr: config.Addrs[0],
			DB:   config.DB,
		}
		if config.Username != "" {
			c.Username = config.Username
		}
		if config.Password != "" {
			c.Password = config.Password
		}
		if config.TLS != nil {
			c.TLSConfig = config.TLS
		}
		return redis.NewClient(c)
	} else {
		c := &redis.ClusterOptions{
			Addrs: config.Addrs,
		}
		if config.Username != "" {
			c.Username = config.Username
		}
		if config.Password != "" {
			c.Password = config.Password
		}
		if config.TLS != nil {
			c.TLSConfig = config.TLS
		}
		return redis.NewClusterClient(c)
	}
}
