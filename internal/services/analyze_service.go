package services

import (
	"context"
	"fmt"
	"github.com/mertbahardogan/escope/internal/interfaces"
	"github.com/mertbahardogan/escope/internal/models"
)

type AnalyzeService interface {
	AnalyzeText(ctx context.Context, analyzerName, text string, analyzeType string) (models.AnalyzeResult, error)
}

type analyzeService struct {
	client interfaces.ElasticClient
}

func NewAnalyzeService(client interfaces.ElasticClient) AnalyzeService {
	return &analyzeService{
		client: client,
	}
}

func (a *analyzeService) AnalyzeText(ctx context.Context, analyzerName, text string, analyzeType string) (models.AnalyzeResult, error) {
	result, err := a.client.GetAnalyze(ctx, analyzerName, text, analyzeType)
	if err != nil {
		return models.AnalyzeResult{}, fmt.Errorf("analyze request failed: %w", err)
	}

	if errorData, ok := result["error"].(map[string]interface{}); ok {
		if reason, ok := errorData["reason"].(string); ok {
			return models.AnalyzeResult{}, fmt.Errorf("elasticsearch error: %s", reason)
		}
		return models.AnalyzeResult{}, fmt.Errorf("elasticsearch error: %v", errorData)
	}

	var analyzeResult models.AnalyzeResult

	if tokens, ok := result["tokens"].([]interface{}); ok {
		for _, tokenData := range tokens {
			if token, ok := tokenData.(map[string]interface{}); ok {
				analyzeToken := models.AnalyzeToken{}

				if tokenStr, ok := token["token"].(string); ok {
					analyzeToken.Token = tokenStr
				}
				if tokenType, ok := token["type"].(string); ok {
					analyzeToken.Type = tokenType
				}
				if position, ok := token["position"].(float64); ok {
					analyzeToken.Position = int(position)
				}
				if startOffset, ok := token["start_offset"].(float64); ok {
					analyzeToken.StartOffset = int(startOffset)
				}
				if endOffset, ok := token["end_offset"].(float64); ok {
					analyzeToken.EndOffset = int(endOffset)
				}

				analyzeResult.Tokens = append(analyzeResult.Tokens, analyzeToken)
			}
		}
	}

	return analyzeResult, nil
}
