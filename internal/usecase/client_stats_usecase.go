package usecase

import "proxy/internal/domain"

type ClientStatsUseCase struct {
	repo domain.ClientStatsRepo
}

func NewClientStatsUseCase(repo domain.ClientStatsRepo) *ClientStatsUseCase {
	return &ClientStatsUseCase{repo: repo}
}

func (uc *ClientStatsUseCase) Update(stat domain.ClientStat) {
	uc.repo.Update(stat)
}

func (uc *ClientStatsUseCase) GetTop(limit int, sortBy string) []domain.ClientStat {
	return uc.repo.GetTop(limit, sortBy)
}
