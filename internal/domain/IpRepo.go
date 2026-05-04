package domain

import "context"

type IpRepo interface {
	GetAll() (map[string]bool, error)
	Add(ctx context.Context, ip string) error
	Remove(ctx context.Context, ip string) error
	IsAllowed(ctx context.Context, ip string) (bool, error)
}
