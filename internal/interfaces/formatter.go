package interfaces

import (
	"escope/internal/models"
)

type GCFormatter interface {
	FormatGCTable(gcInfos []models.GCInfo) string
	FormatGCDetails(gcInfo models.GCInfo) string
}

type TermvectorsFormatter interface {
	FormatTermvectorsTable(termInfos []models.TermInfo) string
	FormatTermSearchResult(termInfos []models.TermInfo, searchTerm string) string
	FormatSummary(termInfos []models.TermInfo) string
}
