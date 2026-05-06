package repository

import (
	"proxy/internal/domain"
	"sort"
	"sync"
)

type ClientStatsRepo struct {
	mu    sync.RWMutex
	stats map[string]*domain.ClientStat
}

func NewClientStatsRepo() *ClientStatsRepo {
	return &ClientStatsRepo{
		stats: make(map[string]*domain.ClientStat),
	}
}

func (r *ClientStatsRepo) Update(stat domain.ClientStat) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.stats[stat.IP]
	if !ok {
		copy := stat
		r.stats[stat.IP] = &copy
		return
	}
	existing.Requests += stat.Requests
	existing.TrafficMB += stat.TrafficMB
	totalReq := float64(existing.Requests)
	existing.AvgLatencyMs = (existing.AvgLatencyMs*(totalReq-float64(stat.Requests)) + stat.AvgLatencyMs*float64(stat.Requests)) / totalReq
	existing.ErrorRate = (existing.ErrorRate*(totalReq-float64(stat.Requests)) + stat.ErrorRate*float64(stat.Requests)) / totalReq
	existing.LastSeen = stat.LastSeen
}

func (r *ClientStatsRepo) GetTop(limit int, sortBy string) []domain.ClientStat {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.ClientStat, 0, len(r.stats))
	for _, v := range r.stats {
		list = append(list, *v)
	}
	switch sortBy {
	case "traffic":
		sort.Slice(list, func(i, j int) bool { return list[i].TrafficMB > list[j].TrafficMB })
	case "errors":
		sort.Slice(list, func(i, j int) bool { return list[i].ErrorRate > list[j].ErrorRate })
	default: // "requests"
		sort.Slice(list, func(i, j int) bool { return list[i].Requests > list[j].Requests })
	}
	if limit > 0 && limit < len(list) {
		list = list[:limit]
	}
	return list
}
