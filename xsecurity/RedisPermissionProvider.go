package xsecurity

import (
	"context"
	"encoding/json"

	"github.com/DreamvatLab/go/xbytes"
	"github.com/DreamvatLab/go/xdto"
	"github.com/DreamvatLab/go/xerr"
	"github.com/DreamvatLab/go/xlog"
	"github.com/DreamvatLab/go/xredis"
	"github.com/redis/go-redis/v9"
)

type RedisPermissionProvider struct {
	redis         redis.Cmdable
	PermissionKey string
}

func NewRedisPermissionProvider(permissionKey string, config *xredis.RedisConfig) IPermissionProvider {
	if permissionKey == "" {
		xlog.Fatal("permissionKey key cannot be empty")
	}

	r := new(RedisPermissionProvider)

	r.redis = xredis.NewClient(config)

	r.PermissionKey = permissionKey

	return r
}

// *******************************************************************************************************************************
// Permission
func (x *RedisPermissionProvider) CreatePermission(in *xdto.Permission) error {
	j, err := json.Marshal(in)
	if err != nil {
		return xerr.WithStack(err)
	}

	cmd := x.redis.HSet(context.Background(), x.PermissionKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisPermissionProvider) GetPermission(id string) (*xdto.Permission, error) {
	cmd := x.redis.HGet(context.Background(), x.PermissionKey, id)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}

	r := new(xdto.Permission)
	j, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(xbytes.StrToBytes(j), r)
	return r, xerr.WithStack(err)
}
func (x *RedisPermissionProvider) UpdatePermission(in *xdto.Permission) error {
	j, err := json.Marshal(in)
	if err != nil {
		return xerr.WithStack(err)
	}

	cmd := x.redis.HSet(context.Background(), x.PermissionKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisPermissionProvider) RemovePermission(id string) error {
	cmd := x.redis.HDel(context.Background(), x.PermissionKey, id)
	return cmd.Err()
}
func (x *RedisPermissionProvider) GetPermissions() (map[string]*xdto.Permission, error) {
	cmd := x.redis.HGetAll(context.Background(), x.PermissionKey)
	r, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	m := make(map[string]*xdto.Permission, len(r))
	for key, value := range r {
		dto := new(xdto.Permission)
		err = json.Unmarshal(xbytes.StrToBytes(value), dto)
		if err == nil {
			m[key] = dto
		} else {
			xlog.Errorf("%s, %v", value, err)
		}
	}
	return m, err
}
