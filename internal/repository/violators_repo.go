package repository

import (
	"proxy/internal/domain"
	"sync"
	"time"
)

type ViolatorsRepo struct {
	mu        sync.RWMutex
	violators map[string]*domain.Violator
}

func NewViolatorsRepo() *ViolatorsRepo {
	r := &ViolatorsRepo{
		violators: make(map[string]*domain.Violator),
	}
	go r.cleanExpired()
	return r
}

func (r *ViolatorsRepo) Add(v domain.Violator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.violators[v.IP] = &v
}

func (r *ViolatorsRepo) GetAll() []domain.Violator {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.Violator, 0, len(r.violators))
	for _, v := range r.violators {
		list = append(list, *v)
	}
	return list
}

func (r *ViolatorsRepo) Remove(ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.violators, ip)
}

func (r *ViolatorsRepo) cleanExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		now := time.Now()
		r.mu.Lock()
		for ip, v := range r.violators {
			if now.After(v.BannedUntil) {
				delete(r.violators, ip)
			}
		}
		r.mu.Unlock()
	}
}
