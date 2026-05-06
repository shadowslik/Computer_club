package usecase

import "proxy/internal/domain"

type DenialsUseCase struct {
	repo domain.DenialsRepo
}

func NewDenialsUseCase(repo domain.DenialsRepo) *DenialsUseCase {
	return &DenialsUseCase{repo: repo}
}

func (uc *DenialsUseCase) Record(ip, reason string) {
	uc.repo.Record(ip, reason)
}

func (uc *DenialsUseCase) GetAll() []domain.DenialStat {
	return uc.repo.GetAll()
}
