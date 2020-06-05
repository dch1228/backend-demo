package user

import (
	"context"

	"github.com/duchenhao/backend-demo/internal/bus"
	"github.com/duchenhao/backend-demo/internal/model"
)

func init() {
	bus.AddHandler(getSignedInUser)
}

func getSignedInUser(ctx context.Context, query *model.GetSignedInUserQuery) error {
	userQuery := &model.GetUserByIdQuery{UserId: query.UserId}
	if err := bus.Dispatch(ctx, userQuery); err != nil {
		ctx.(*model.ReqContext).Logger.Error(err.Error())
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
