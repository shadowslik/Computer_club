package usecase

import (
	"context"
	"errors"
	"net"
	"proxy/internal/domain"
	"strings"
)

type ListUseCase struct {
	repo domain.IpRepo
}

func NewListUseCase(repo domain.IpRepo) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) GetAll() (map[string]bool, error) {
	return uc.repo.GetAll()
}

func (uc *ListUseCase) IsAllowed(ctx context.Context, ip string) (bool, error) {
	if !checkIp(ip) {
		return false, errors.New("ip address is not valid")
	}
	return uc.repo.IsAllowed(ctx, ip)
}

func (uc *ListUseCase) Add(ctx context.Context, ip string) (bool, error) {
	if !checkIp(ip) {
		return false, errors.New("ip address is not valid")
	}
	allowed, err := uc.repo.IsAllowed(ctx, ip)
	if err != nil {
		return false, err
	}
	if allowed {
		return false, nil
	}
	if err := uc.repo.Add(ctx, ip); err != nil {
		return false, err
	}
	return true, nil
}

func (uc *ListUseCase) Remove(ctx context.Context, ip string) (bool, error) {
	if !checkIp(ip) {
		return false, errors.New("ip address is not valid")
	}
	allowed, err := uc.repo.IsAllowed(ctx, ip)
	if err != nil {
		return false, err
	}
	if !allowed {
		return false, nil
	}
	if err := uc.repo.Remove(ctx, ip); err != nil {
		return false, err
	}
	return true, nil
}

func checkIp(s string) bool {
	s = strings.TrimSpace(s)
	if strings.Contains(s, "/") {
		_, _, err := net.ParseCIDR(s)
		return err == nil
	}
	if strings.Contains(s, "-") {
		parts := strings.SplitN(s, "-", 2)
		if len(parts) != 2 {
			return false
		}
		startIP := net.ParseIP(strings.TrimSpace(parts[0]))
		endIP := net.ParseIP(strings.TrimSpace(parts[1]))
		if startIP == nil || endIP == nil {
			return false
		}
		return true
	}

	ip := net.ParseIP(s)
	return ip != nil
}
