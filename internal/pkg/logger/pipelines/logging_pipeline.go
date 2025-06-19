package pipelines

import (
	"context"
	"fmt"
	"time"

	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/reflection/typemapper"
	"github.com/mehdihadeli/go-mediatr"
)

type requestLoggerPipeline struct {
	logger logger.Logger
}

func NewMediatorLoggingPipeline(l logger.Logger) mediatr.PipelineBehavior {
	return &requestLoggerPipeline{logger: l}
}

func (r *requestLoggerPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		r.logger.Infof("Request took %s", elapsed)
	}()

	requestName := typemapper.GetNonePointerTypeName(request)

	r.logger.Infow(
		fmt.Sprintf("Handling request: '%s'", requestName),
		logger.Fields{"Request": request},
	)

	response, err := next(ctx)
	if err != nil {
		r.logger.Infof("Request failed with error: %v", err)

		return nil, err
	}

	responseName := typemapper.GetNonePointerTypeName(response)

	r.logger.Infow(
		fmt.Sprintf(
			"Request handled successfully with response: '%s'",
			responseName,
		),
		logger.Fields{"Response": response},
	)

	return response, nil
}
