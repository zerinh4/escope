package services

import "github.com/mertbahardogan/escope/internal/constants"

func parseIndexStatsData(statsData map[string]interface{}) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	indices, ok := statsData[constants.IndicesField].(map[string]interface{})
	if !ok {
		return result
	}

	for indexName, indexData := range indices {
		if index, ok := indexData.(map[string]interface{}); ok {
			if total, ok := index[constants.TotalField].(map[string]interface{}); ok {
				result[indexName] = total
			}
		}
	}

	return result
}

func getSegmentsData(total map[string]interface{}) (map[string]interface{}, bool) {
	segments, ok := total[constants.SegmentsField].(map[string]interface{})
	return segments, ok
}

func getIndexingData(total map[string]interface{}) (map[string]interface{}, bool) {
	indexing, ok := total[constants.IndexingField].(map[string]interface{})
	return indexing, ok
}
