package user

import (
	"time"

	"github.com/rs/xid"

	"github.com/duchenhao/backend-demo/internal/bus"
	"github.com/duchenhao/backend-demo/internal/events"
	"github.com/duchenhao/backend-demo/internal/model"
	"github.com/duchenhao/backend-demo/internal/util"
)

func init() {
	bus.AddHandler(getSignedInUser)
	bus.AddHandler(signUp)
	bus.AddHandler(login)
}

func signUp(cmd *model.SignUpCommand) error {
	user := &model.User{
		Id:         xid.New().String(),
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

	createCmd := &model.CreateUserCommand{Ctx: cmd.Ctx, User: user}
	if err := bus.Dispatch(createCmd); err != nil {
		return err
	}

	cmd.User = createCmd.User

	event := &events.OnCreateUser{}
	event.UserId = cmd.User.Id
	bus.Publish(event)
	return nil
}

func login(query *model.LoginQuery) error {
	userQuery := &model.GetUserByNameQuery{Ctx: query.Ctx, Name: query.Name}
	if err := bus.Dispatch(userQuery); err != nil {
		query.Ctx.Logger.Error(err.Error())
		return err
	}

	user := userQuery.User
	if err := util.ValidatePassword(query.Password, user.Password, user.Salt); err != nil {
		return err
	}

	query.User = user
	return nil
}

func getSignedInUser(query *model.GetSignedInUserQuery) error {
	userQuery := &model.GetUserByIdQuery{Ctx: query.Ctx, UserId: query.UserId}
	if err := bus.Dispatch(userQuery); err != nil {
		query.Ctx.Logger.Error(err.Error())
		return err
	}

	user := userQuery.User

	query.User = &model.SignedInUser{
		UserId:     user.Id,
		Name:       user.Name,
		LastSeenAt: user.LastSeenAt,
	}
	return nil
}
