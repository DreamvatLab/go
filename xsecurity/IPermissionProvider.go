package xsecurity

import "github.com/DreamvatLab/go/xdto"

type IPermissionProvider interface {
	CreatePermission(*xdto.Permission) error
	GetPermission(string) (*xdto.Permission, error)
	UpdatePermission(*xdto.Permission) error
	RemovePermission(string) error
	GetPermissions() (map[string]*xdto.Permission, error)
}
