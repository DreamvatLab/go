package xredis

import (
	"crypto/tls"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestParseRedisConfig(t *testing.T) {
	tests := []struct {
		name      string
		connStr   string
		want      *RedisConfig
		wantErr   bool
		errString string
	}{
		{
			name:    "single node without auth",
			connStr: "redis://localhost:6379",
			want: &RedisConfig{
				Addrs: []string{"localhost:6379"},
				DB:    0,
			},
			wantErr: false,
		},
		{
			name:    "single node with auth",
			connStr: "redis://username:password@localhost:6379",
			want: &RedisConfig{
				Addrs:    []string{"localhost:6379"},
				Username: "username",
				Password: "password",
				DB:       0,
			},
			wantErr: false,
		},
		{
			name:    "single node with db",
			connStr: "redis://localhost:6379/1",
			want: &RedisConfig{
				Addrs: []string{"localhost:6379"},
				DB:    1,
			},
			wantErr: false,
		},
		{
			name:    "single node with auth and db",
			connStr: "redis://username:password@localhost:6379/2",
			want: &RedisConfig{
				Addrs:    []string{"localhost:6379"},
				Username: "username",
				Password: "password",
				DB:       2,
			},
			wantErr: false,
		},
		{
			name:    "cluster mode without auth",
			connStr: "redis://localhost:6379,localhost:6380",
			want: &RedisConfig{
				Addrs: []string{"localhost:6379", "localhost:6380"},
				DB:    0,
			},
			wantErr: false,
		},
		{
			name:    "cluster mode with auth",
			connStr: "redis://username:password@localhost:6379,localhost:6380",
			want: &RedisConfig{
				Addrs:    []string{"localhost:6379", "localhost:6380"},
				Username: "username",
				Password: "password",
				DB:       0,
			},
			wantErr: false,
		},
		{
			name:    "cluster mode with db",
			connStr: "redis://localhost:6379,localhost:6380/3",
			want: &RedisConfig{
				Addrs: []string{"localhost:6379", "localhost:6380"},
				DB:    3,
			},
			wantErr: false,
		},
		{
			name:    "cluster mode with auth and db",
			connStr: "redis://username:password@localhost:6379,localhost:6380/4",
			want: &RedisConfig{
				Addrs:    []string{"localhost:6379", "localhost:6380"},
				Username: "username",
				Password: "password",
				DB:       4,
			},
			wantErr: false,
		},
		{
			name:      "empty connection string",
			connStr:   "",
			want:      nil,
			wantErr:   true,
			errString: "empty connection string",
		},
		{
			name:      "invalid auth format",
			connStr:   "redis://username@localhost:6379",
			want:      nil,
			wantErr:   true,
			errString: "invalid authentication format: username",
		},
		{
			name:      "invalid db format",
			connStr:   "redis://localhost:6379/invalid",
			want:      nil,
			wantErr:   true,
			errString: "invalid database number: invalid",
		},
		{
			name:      "empty address in cluster",
			connStr:   "redis://localhost:6379,,localhost:6380",
			want:      nil,
			wantErr:   true,
			errString: "empty address in cluster configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRedisConfig(tt.connStr)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.Contains(t, err.Error(), tt.errString)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Addrs, got.Addrs)
			assert.Equal(t, tt.want.Username, got.Username)
			assert.Equal(t, tt.want.Password, got.Password)
			assert.Equal(t, tt.want.DB, got.DB)
		})
	}
}

func TestNewClient(t *testing.T) {
	// Test single node client
	singleNodeConfig := &RedisConfig{
		Addrs:    []string{"localhost:6379"},
		Username: "username",
		Password: "password",
		DB:       0,
	}
	singleNodeClient := NewClient(singleNodeConfig)
	assert.NotNil(t, singleNodeClient)
	assert.IsType(t, &redis.Client{}, singleNodeClient)

	// Test cluster client
	clusterConfig := &RedisConfig{
		Addrs:    []string{"localhost:6379", "localhost:6380"},
		Username: "username",
		Password: "password",
		DB:       0,
	}
	clusterClient := NewClient(clusterConfig)
	assert.NotNil(t, clusterClient)
	assert.IsType(t, &redis.ClusterClient{}, clusterClient)

	// Test with TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	tlsRedisConfig := &RedisConfig{
		Addrs:    []string{"localhost:6379"},
		Username: "username",
		Password: "password",
		DB:       0,
		TLS:      tlsConfig,
	}
	tlsClient := NewClient(tlsRedisConfig)
	assert.NotNil(t, tlsClient)
	assert.IsType(t, &redis.Client{}, tlsClient)
}
