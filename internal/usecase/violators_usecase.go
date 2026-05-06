package usecase

import "proxy/internal/domain"

type ViolatorsUseCase struct {
	repo domain.ViolatorsRepo
}

func NewViolatorsUseCase(repo domain.ViolatorsRepo) *ViolatorsUseCase {
	return &ViolatorsUseCase{repo: repo}
}

func (uc *ViolatorsUseCase) Add(v domain.Violator) {
	uc.repo.Add(v)
}

func (uc *ViolatorsUseCase) GetAll() []domain.Violator {
	return uc.repo.GetAll()
}

func (uc *ViolatorsUseCase) Unban(ip string) {
	uc.repo.Remove(ip)
}
