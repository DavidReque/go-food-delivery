package mediator

import "github.com/DavidReque/go-food-delivery/internal/pkg/core/cqrs"

func RegisterMediatorHandlers(handlers []cqrs.HandlerRegisterer) error {
	for _, handler := range handlers {
		err := handler.RegisterHandler()
		if err != nil {
			return err
		}
	}

	return nil
}
