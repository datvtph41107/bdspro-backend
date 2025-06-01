package user_service

import (
	"user/repo"
)

type CertificationService struct {
	Repo        *repo.CertificationRepo
	ProfileRepo *repo.ProfileRepo
}

func NewCertificationService(repo *repo.CertificationRepo, profileRepo *repo.ProfileRepo) *CertificationService {
	return &CertificationService{
		Repo:        repo,
		ProfileRepo: profileRepo,
	}
}
