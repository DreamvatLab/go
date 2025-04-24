package xsecurity

import "github.com/DreamvatLab/go/xdto"

type IRouteProvider interface {
	CreateRoute(*xdto.Route) error
	GetRoute(string) (*xdto.Route, error)
	UpdateRoute(*xdto.Route) error
	RemoveRoute(string) error
	GetRoutes() (map[string]*xdto.Route, error)
}
