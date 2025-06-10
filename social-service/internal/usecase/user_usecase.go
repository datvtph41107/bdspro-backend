package usecase

import (
	"context"
	"fmt"
	iusecase "social/internal/interface"
	"strings"
)

type UserUsecase struct {
	userClient iusecase.UserClient
}

func NewUserUsecase(userClient iusecase.UserClient) *UserUsecase {
	return &UserUsecase{
		userClient: userClient,
	}
}

func (u *UserUsecase) getFullNames(ctx context.Context, userIds []uint64, sum int) string {
	users, _ := u.userClient.GetUsers(ctx, userIds)
	var names []string
	for _, user := range users {
		names = append(names, user.FullName)
	}

	// todo: tổng số người
	var fullNames string
	if len(users) == sum || sum < 4 {
		fullNames = strings.Join(names, ", ")
	} else {
		fullNames = fmt.Sprintf("%s và %d người khác", strings.Join(names, ", "), sum-len(users))
	}
	return fullNames
}
