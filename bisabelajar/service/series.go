package service

import (
	middlewarev1 "bisabelajar/api/v1/middleware"
	"bisabelajar/domain"
	"bisabelajar/dto"
	"bisabelajar/repository"
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type SeriesService struct {
	seriesRepo repository.SeriesRepository
	tracer     trace.Tracer
	log        *logrus.Logger
}

func NewSeriesService(seriesRepo repository.SeriesRepository, log *logrus.Logger, tracer trace.Tracer) SeriesService {
	return SeriesService{seriesRepo: seriesRepo, log: log, tracer: tracer}
}

func (b *SeriesService) InsertSeries(seriesDto dto.SeriesDto, requestDetail middlewarev1.RequestDetails) (int, error) {
	_, span := b.tracer.Start(context.Background(), "InsertSeries")
	defer span.End()

	span.AddEvent("Inserting series")

	id, err := b.seriesRepo.InsertSeries(seriesDto)
	if err != nil {
		span.RecordError(err)
		b.log.Errorf("error insert series because %s with request details %v", err.Error(), requestDetail)
		return 0, err
	}

	b.log.Infof("request %v done", requestDetail)
	return id, nil
}

func (b *SeriesService) GetSeriesByID(id int, requestDetail middlewarev1.RequestDetails) (*domain.Series, error) {
	_, span := b.tracer.Start(context.Background(), "GetSeriesByID")
	defer span.End()

	span.AddEvent("Get series")

	series, err := b.seriesRepo.GetSeriesByID(id)
	if err != nil {
		span.RecordError(err)
		b.log.Errorf("error get series because %s with request details %v", err.Error(), requestDetail)
		return nil, err
	}

	b.log.Infof("request %v done", requestDetail)
	return series, nil
}
