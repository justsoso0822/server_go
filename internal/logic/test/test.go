package test

import (
	"context"
	"server_go/internal/service"
)

type sTest struct{}

func init() {
	service.RegisterTest(&sTest{})
}

func (s *sTest) Index(ctx context.Context) (any, error) {
	return "test", nil
}
