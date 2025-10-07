package services

import (
	"context"
	"escope/internal/interfaces"
	"escope/internal/models"
	"escope/internal/util"
	"fmt"
)

type SystemService interface {
	GetSystemIndices(ctx context.Context) ([]models.IndexInfo, error)
	GetSystemShards(ctx context.Context) ([]models.ShardInfo, error)
}

type systemService struct {
	client interfaces.ElasticClient
}

func NewSystemService(client interfaces.ElasticClient) SystemService {
	return &systemService{
		client: client,
	}
}

func (s *systemService) GetSystemIndices(ctx context.Context) ([]models.IndexInfo, error) {
	indicesData, err := s.client.GetIndices(ctx)
	if err != nil {
		return nil, fmt.Errorf("indices request failed: %w", err)
	}

	var indices []models.IndexInfo
	if indicesList, ok := indicesData[""].([]map[string]interface{}); ok {
		for _, idx := range indicesList {
			index := models.IndexInfo{
				Alias:     util.GetStringField(idx, "alias"),
				Name:      util.GetStringField(idx, "index"),
				Health:    util.GetStringField(idx, "health"),
				Status:    util.GetStringField(idx, "status"),
				DocsCount: util.GetStringField(idx, "docs.count"),
				StoreSize: util.GetStringField(idx, "store.size"),
				Primary:   util.GetStringField(idx, "pri"),
				Replica:   util.GetStringField(idx, "rep"),
			}
			indices = append(indices, index)
		}
	}

	return indices, nil
}

func (s *systemService) GetSystemShards(ctx context.Context) ([]models.ShardInfo, error) {
	shardsData, err := s.client.GetShards(ctx)
	if err != nil {
		return nil, fmt.Errorf("shards request failed: %w", err)
	}

	var shards []models.ShardInfo
	if shardsList, ok := shardsData[""].([]map[string]interface{}); ok {
		for _, shard := range shardsList {
			shardInfo := models.ShardInfo{
				Index:  util.GetStringField(shard, "index"),
				Shard:  util.GetStringField(shard, "shard"),
				Prirep: util.GetStringField(shard, "prirep"),
				State:  util.GetStringField(shard, "state"),
				Docs:   util.GetStringField(shard, "docs"),
				Store:  util.GetStringField(shard, "store"),
				IP:     util.GetStringField(shard, "ip"),
				Node:   util.GetStringField(shard, "node"),
			}
			shards = append(shards, shardInfo)
		}
	}

	return shards, nil
}
