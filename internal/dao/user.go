package dao

import (
	"context"
	"time"

	"github.com/duchenhao/backend-demo/internal/bus"
	"github.com/duchenhao/backend-demo/internal/model"
	"github.com/duchenhao/backend-demo/internal/util"
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

func createUser(ctx context.Context, cmd *model.CreateUserCommand) error {
	user := &model.User{
		Name:       cmd.Name,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		LastSeenAt: time.Now().AddDate(-10, 0, 0),
	}

	salt, err := util.GetRandomString(8)
	if err != nil {
		return err
	}
	user.Salt = salt

	encodedPassword, err := util.EncodePassword(cmd.Password, user.Salt)
	if err != nil {
		return err
	}
	user.Password = encodedPassword

	_, err = db.Context(ctx).Insert(user)
	if err != nil {
		return err
	}

	cmd.User = user
	return nil
}

func updateUserLastSeenAt(ctx context.Context, cmd *model.UpdateUserLastSeenAtCommand) error {
	_, err := db.Context(ctx).
		Table(&model.User{}).ID(cmd.UserId).
		Update(model.User{LastSeenAt: time.Now()})
	return err
}

func getUserById(ctx context.Context, query *model.GetUserByIdQuery) error {
	user := &model.User{}
	if has, err := db.Context(ctx).ID(query.UserId).Get(user); err != nil {
		return err
	} else if !has {
		return model.ErrUserNotFound
	}
	query.User = user
	return nil
}
