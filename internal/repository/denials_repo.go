package repository

import (
	"proxy/internal/domain"
	"sync"
	"time"
)

type DenialsRepo struct {
	mu    sync.RWMutex
	stats map[string]*domain.DenialStat
}

func NewDenialsRepo() *DenialsRepo {
	return &DenialsRepo{
		stats: make(map[string]*domain.DenialStat),
	}
}

func (r *DenialsRepo) Record(ip, reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stat, ok := r.stats[ip]
	if !ok {
		r.stats[ip] = &domain.DenialStat{
			IP:         ip,
			Denials:    1,
			LastDenied: time.Now(),
			Reason:     reason,
		}
		return
	}
	stat.Denials++
	stat.LastDenied = time.Now()
	stat.Reason = reason
}

func (r *DenialsRepo) GetAll() []domain.DenialStat {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.DenialStat, 0, len(r.stats))
	for _, v := range r.stats {
		list = append(list, *v)
	}
	return list
}
