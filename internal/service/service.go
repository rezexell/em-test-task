package service

import (
	"github.com/rezexell/em-test-task/internal/repository"
)

type Subscription interface {
	GetAllSubs() string
}
type Service struct {
	Subscription
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Subscription: NewSubService(repo.Subscription)}
}
