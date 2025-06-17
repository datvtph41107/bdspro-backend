package user_service

import user_repo "user/infra/repo"

type CertificationService struct {
	Repo        *user_repo.CertificationRepo
	ProfileRepo *user_repo.ProfileRepo
}

func NewCertificationService(repo *user_repo.CertificationRepo, profileRepo *user_repo.ProfileRepo) *CertificationService {
	return &CertificationService{
		Repo:        repo,
		ProfileRepo: profileRepo,
	}
}
