package account

type AccountService struct {
	Repository AccountRepository
}

func NewAccountService(repository AccountRepository) *AccountService {
	return &AccountService{Repository: repository}
}

func (s *AccountService) CreateAccount(account Account) (Account, error) {
	return s.Repository.CreateAccount(account)
}

func (s *AccountService) GetAccountBalance(id int) (int, error) {
	return s.Repository.GetAccountBalance(id)
}
