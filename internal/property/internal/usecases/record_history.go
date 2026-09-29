package usecases

import "bdspro/internal/repo"

type RecordHistoryUsecase struct {
	Repo repo.RecordHistoryRepo
}

func NewRecordHistoryUsecase(repo repo.RecordHistoryRepo) RecordHistoryUsecase {
	return RecordHistoryUsecase{Repo: repo}
}
