package usecase

type Usecase struct {
	repo UserRepository
}

func New(r UserRepository) *Usecase {
	return &Usecase{
		repo: r,
	}
}
