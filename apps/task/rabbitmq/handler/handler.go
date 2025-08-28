package handler

import (
	"PaiPai/apps/task/rabbitmq/internal/svc"
	"context"
	"log"
)

type UserHandler struct {
	svcCtx *svc.ServiceContext
}

func (h *UserHandler) Error() string {
	//TODO implement me
	panic("implement me")
}

func NewUserHandler(svcCtx *svc.ServiceContext) *UserHandler {
	return &UserHandler{svcCtx: svcCtx}
}

func HandleCreated(ctx context.Context, message map[string]interface{}) error {
	log.Printf("Processing user created event: %+v", message)

	// 可添加业务逻辑处理
	//userID, ok := message["user_id"].(float64)
	//if !ok {
	//	return fmt.Errorf("invalid user_id in message")
	//}
	//
	//log.Printf("Successfully processed user creation for ID: %v", userID)

	return nil
}
