package handler

import (
	"github.com/Mire0726/unibox/backend/app/usecase"
)

type Handler struct {
	authUC      usecase.AuthUsecase
	workspaceUC usecase.WorkspaceUsecase
	massageUC   usecase.MessageUsecase
	channelUC   usecase.ChannelUsecase
}

func NewHandler(
	authUC usecase.AuthUsecase,
	channelUC usecase.ChannelUsecase,
	massageUC usecase.MessageUsecase,
	workspaceUC usecase.WorkspaceUsecase,
) *Handler {
	return &Handler{
		authUC:      authUC,
		channelUC:   channelUC,
		massageUC:   massageUC,
		workspaceUC: workspaceUC,
	}
}
