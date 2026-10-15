package services

import (
	"userv4/internal/interface/repo"
)

type CertificationUsecase struct {
	Repo        repo.ICertificationRepo
	ProfileRepo repo.IProfileRepo
}

func NewCertificationUsecase(repo repo.ICertificationRepo, profileRepo repo.IProfileRepo) *CertificationUsecase {
	return &CertificationUsecase{
		Repo:        repo,
		ProfileRepo: profileRepo,
	}
}
