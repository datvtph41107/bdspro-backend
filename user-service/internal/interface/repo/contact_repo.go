package repo

type IContactRepo interface {
	ExistByProfile(profileId uint64) (bool, error)
	Existed(followId uint64) error
}
