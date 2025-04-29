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

type RedisRouteProvider struct {
	redis    redis.Cmdable
	RouteKey string
}

func NewRedisRouteProvider(routeKey string, config *xredis.RedisConfig) IRouteProvider {
	if routeKey == "" {
		xlog.Fatal("routeKey cannot be empty")
	}

	r := new(RedisRouteProvider)

	r.redis = xredis.NewClient(config)

	r.RouteKey = routeKey

	return r
}

// *******************************************************************************************************************************
// Route
func (x *RedisRouteProvider) CreateRoute(in *xdto.Route) error {
	j, err := json.Marshal(in)
	if err != nil {
		return xerr.WithStack(err)
	}

	cmd := x.redis.HSet(context.Background(), x.RouteKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRouteProvider) GetRoute(id string) (*xdto.Route, error) {
	cmd := x.redis.HGet(context.Background(), x.RouteKey, id)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}

	r := new(xdto.Route)
	j, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(xbytes.StrToBytes(j), r)
	return r, xerr.WithStack(err)
}
func (x *RedisRouteProvider) UpdateRoute(in *xdto.Route) error {
	j, err := json.Marshal(in)
	if err != nil {
		return xerr.WithStack(err)
	}

	cmd := x.redis.HSet(context.Background(), x.RouteKey, in.ID, j)
	return cmd.Err()
}
func (x *RedisRouteProvider) RemoveRoute(id string) error {
	cmd := x.redis.HDel(context.Background(), x.RouteKey, id)
	return cmd.Err()
}
func (x *RedisRouteProvider) GetRoutes() (map[string]*xdto.Route, error) {
	cmd := x.redis.HGetAll(context.Background(), x.RouteKey)
	r, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	m := make(map[string]*xdto.Route, len(r))
	for key, value := range r {
		dto := new(xdto.Route)
		err = json.Unmarshal(xbytes.StrToBytes(value), dto)
		if err != nil {
			xlog.Errorf("%v\nKey:%s, Value: %s", err, key, value)
			continue
		}

		m[key] = dto
	}
	return m, err
}
