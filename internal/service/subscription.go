package service

import (
	"github.com/rezexell/em-test-task/internal/repository"
)

type SubService struct {
	repo repository.Subscription
}

func NewSubService(repo repository.Subscription) *SubService {
	return &SubService{repo: repo}
}

func (s *SubService) GetAllSubs() string {
	return s.repo.GetAllSubs()
}
