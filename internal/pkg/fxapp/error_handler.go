package fxapp

import "github.com/DavidReque/go-food-delivery/internal/pkg/logger"

type FxErrorHandler struct {
	logger.Logger
}

func NewFxErrorHandler(logger logger.Logger) *FxErrorHandler {
	return &FxErrorHandler{Logger: logger}
}

func (h *FxErrorHandler) HandleError(e error) {
	h.Logger.Error(e)
}