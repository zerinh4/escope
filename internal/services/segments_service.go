package services

import (
	"context"
	"fmt"
	"github.com/mertbahardogan/escope/internal/constants"
	"github.com/mertbahardogan/escope/internal/interfaces"
	"github.com/mertbahardogan/escope/internal/models"
)

type SegmentsService interface {
	GetSegmentsInfo(ctx context.Context) ([]models.SegmentInfo, error)
}

type segmentsService struct {
	client interfaces.ElasticClient
}

func NewSegmentsService(client interfaces.ElasticClient) SegmentsService {
	return &segmentsService{
		client: client,
	}
}

func (s *segmentsService) GetSegmentsInfo(ctx context.Context) ([]models.SegmentInfo, error) {
	statsData, err := s.client.GetIndexStats(ctx, constants.EmptyString)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrIndexStatsRequestFailed, err)
	}

	segments := s.parseSegmentsData(statsData)

	return segments, nil
}

func (s *segmentsService) parseSegmentsData(statsData map[string]interface{}) []models.SegmentInfo {
	var segments []models.SegmentInfo

	indexStats := parseIndexStatsData(statsData)

	for indexName, total := range indexStats {
		if segmentsData, ok := getSegmentsData(total); ok {
			segmentInfo := models.SegmentInfo{
				Index: indexName,
			}
			if count, ok := segmentsData[constants.CountField].(float64); ok {
				segmentInfo.SegmentCount = int(count)
			}
			if sizeBytes, ok := segmentsData[constants.MemoryInBytesField].(float64); ok {
				segmentInfo.SizeBytes = int64(sizeBytes)
			}

			segments = append(segments, segmentInfo)
		}
	}

	return segments
}
