package service

import (
	"bisabelajar/domain"
	"bisabelajar/dto"
	"bisabelajar/repository"
)

type SeriesService struct {
	seriesRepo repository.SeriesRepository
}

func NewSeriesService(seriesRepo repository.SeriesRepository) SeriesService {
	return SeriesService{seriesRepo: seriesRepo}
}

func (b *SeriesService) InsertSeries(seriesDto dto.SeriesDto) (int, error) {
	id, err := b.seriesRepo.InsertSeries(seriesDto)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (b *SeriesService) GetSeriesByID(id int) (*domain.Series, error) {
	series, err := b.seriesRepo.GetSeriesByID(id)
	if err != nil {
		return nil, err
	}
	return series, nil
}
