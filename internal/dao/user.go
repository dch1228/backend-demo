package dao

import (
	"context"
	"time"

	"github.com/duchenhao/backend-demo/internal/bus"
	"github.com/duchenhao/backend-demo/internal/model"
)

func init() {
	bus.AddHandler(getUserByName)
	bus.AddHandler(createUser)
	bus.AddHandler(updateUserLastSeenAt)
	bus.AddHandler(getUserById)
}

func getUserByName(ctx context.Context, query *model.GetUserByNameQuery) error {
	user := &model.User{}

	if has, err := db.Context(ctx).Where("name=?", query.Name).Get(user); err != nil {
		return err
	} else if !has {
		return model.ErrUserNotFound
	}

	query.User = user
	return nil
}

func createUser(cmd *model.CreateUserCommand) error {
	_, err := db.Context(cmd.Ctx).Insert(cmd.User)
	if err != nil {
		return err
	}
	return nil
}

func updateUserLastSeenAt(cmd *model.UpdateUserLastSeenAtCommand) error {
	_, err := db.Context(cmd.Ctx).
		Table(&model.User{}).ID(cmd.UserId).
		Update(model.User{LastSeenAt: time.Now()})
	return err
}

func getUserById(query *model.GetUserByIdQuery) error {
	user := &model.User{}
	if has, err := db.Context(query.Ctx).ID(query.UserId).Get(user); err != nil {
		return err
	} else if !has {
		return model.ErrUserNotFound
	}
	query.User = user
	return nil
}
