package account

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type AccountService struct {
	Repository    AccountRepository
	RedisClient   *redis.Client
	CacheDuration int
}

func NewAccountService(repository AccountRepository, redisAddr string, cacheDuration int) *AccountService {
	redis := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &AccountService{Repository: repository, RedisClient: redis, CacheDuration: cacheDuration}
}

func (s *AccountService) CreateAccount(account Account) (Account, error) {
	return s.Repository.CreateAccount(account)
}

func (s *AccountService) GetAccountBalance(id int) (int, error) {
	ctx := context.Background()
	cachedBalance, err := s.RedisClient.Get(ctx, "account_balance_"+strconv.Itoa(id)).Int()
	if err == nil {
		return cachedBalance, nil
	} else if err != redis.Nil {
		return 0, err
	}

	balance, err := s.Repository.GetAccountBalance(id)
	if err != nil {
		return 0, err
	}

	s.RedisClient.Set(ctx, "account_balance_"+strconv.Itoa(id), balance, time.Duration(s.CacheDuration)*time.Minute)

	return balance, nil
}
