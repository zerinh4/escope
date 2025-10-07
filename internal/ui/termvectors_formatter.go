package ui

import (
	"escope/internal/interfaces"
	"escope/internal/models"
	"fmt"
	"sort"
	"strings"
)

type TermvectorsFormatter struct{}

func NewTermvectorsFormatter() interfaces.TermvectorsFormatter {
	return &TermvectorsFormatter{}
}

func (f *TermvectorsFormatter) FormatTermvectorsTable(termInfos []models.TermInfo) string {
	if len(termInfos) == 0 {
		return "No term vectors found\n"
	}

	fieldGroups := make(map[string][]models.TermInfo)
	for _, termInfo := range termInfos {
		fieldGroups[termInfo.Field] = append(fieldGroups[termInfo.Field], termInfo)
	}

	var output strings.Builder
	var fields []string
	for field := range fieldGroups {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	for _, fieldName := range fields {
		terms := fieldGroups[fieldName]

		sort.Slice(terms, func(i, j int) bool {
			return terms[i].TermFreq > terms[j].TermFreq
		})

		output.WriteString(fmt.Sprintf("\nField: %s (%d terms)\n", fieldName, len(terms)))

		headers := []string{"Term", "Frequency"}
		rows := make([][]string, 0, len(terms))

		for _, term := range terms {
			displayTerm := term.Term
			if len(displayTerm) > 28 {
				displayTerm = term.Term[:25] + "..."
			}

			row := []string{
				displayTerm,
				fmt.Sprintf("%d", term.TermFreq),
			}
			rows = append(rows, row)
		}

		formatter := NewGenericTableFormatter()
		output.WriteString(formatter.FormatTable(headers, rows))
	}

	return output.String()
}

func (f *TermvectorsFormatter) FormatTermSearchResult(termInfos []models.TermInfo, searchTerm string) string {
	var output strings.Builder

	fieldGroups := make(map[string][]models.TermInfo)
	for _, termInfo := range termInfos {
		fieldGroups[termInfo.Field] = append(fieldGroups[termInfo.Field], termInfo)
	}

	found := false
	var foundTerms []models.TermInfo
	for _, termInfo := range termInfos {
		if termInfo.Term == searchTerm {
			found = true
			foundTerms = append(foundTerms, termInfo)
		}
	}

	if found {
		output.WriteString("\nSEARCH TERM FOUND!\n\n")
		output.WriteString(fmt.Sprintf("------- term: %s -------\n\n", searchTerm))
		output.WriteString(" Field            │ Frequency      \n")
		output.WriteString("─────────────────────────────────────\n")

		for _, term := range foundTerms {
			output.WriteString(fmt.Sprintf(" %-16s │ %-15d  \n", term.Field, term.TermFreq))
		}
		output.WriteString("\n")
	} else {
		output.WriteString("\nSEARCH TERM NOT FOUND!\n\n")
	}

	return output.String()
}

func (f *TermvectorsFormatter) FormatSummary(termInfos []models.TermInfo) string {
	if len(termInfos) == 0 {
		return "No term vectors to summarize\n"
	}

	fieldGroups := make(map[string][]models.TermInfo)
	totalTerms := 0
	maxFreq := 0
	var allTerms []string

	for _, termInfo := range termInfos {
		fieldGroups[termInfo.Field] = append(fieldGroups[termInfo.Field], termInfo)
		totalTerms++
		if termInfo.TermFreq > maxFreq {
			maxFreq = termInfo.TermFreq
		}
		allTerms = append(allTerms, termInfo.Term)
	}

	var output strings.Builder

	output.WriteString("\nTERM VECTORS SUMMARY\n")
	output.WriteString("─────────────────────\n")
	output.WriteString(fmt.Sprintf("Total Terms: %d\n", totalTerms))
	output.WriteString(fmt.Sprintf("Fields Analyzed: %d\n", len(fieldGroups)))
	output.WriteString(fmt.Sprintf("Highest Frequency: %d\n\n", maxFreq))

	output.WriteString("FIELD BREAKDOWN\n")
	output.WriteString("────────────────\n")
	for fieldName, terms := range fieldGroups {
		output.WriteString(fmt.Sprintf("   • %s: %d terms\n", fieldName, len(terms)))
	}
	output.WriteString("\n")
	return output.String()
}
