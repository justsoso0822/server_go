package test

import (
	"context"

	"server_go/internal/autodb"
	"server_go/internal/service"
)

type sTest struct{}

func init() {
	service.RegisterTest(&sTest{})
}

func (s *sTest) Index(ctx context.Context) (any, error) {
	ret, err := autodb.Model(ctx, "user u").
		LeftJoin("log_login log", "u.uid=log.uid").
		Fields("u.uid, u.openid, log.time").
		Where("u.uid", 13081).
		Order("log.time desc").
		All()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *sTest) TestDb(ctx context.Context) (any, error) {
	ret, err := autodb.DB(ctx).GetOne(ctx, "select * from user where uid = ?", 13081)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
