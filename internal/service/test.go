// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ITest interface {
		Index(ctx context.Context) (any, error)
		TestDb(ctx context.Context) (any, error)
	}
)

var (
	localTest ITest
)

func Test() ITest {
	if localTest == nil {
		panic("implement not found for interface ITest, forgot register?")
	}
	return localTest
}

func RegisterTest(i ITest) {
	localTest = i
}
