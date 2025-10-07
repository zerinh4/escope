package ui

import (
	"fmt"
)

type IndexDetailFormatter struct{}

func NewIndexDetailFormatter() *IndexDetailFormatter {
	return &IndexDetailFormatter{}
}

func (f *IndexDetailFormatter) FormatIndexDetail(info *IndexDetailInfo) string {
	genericFormatter := NewGenericTableFormatter()

	var headerText string
	if info.CheckCount > 0 {
		headerText = fmt.Sprintf("%s | Check %d", info.Name, info.CheckCount)
	} else {
		headerText = info.Name
	}

	headers := []string{headerText}
	rows := [][]string{
		{fmt.Sprintf("Search Rate: %s", info.SearchRate)},
		{fmt.Sprintf("Index Rate: %s", info.IndexRate)},
		{fmt.Sprintf("Query Time: %s", info.AvgQueryTime)},
		{fmt.Sprintf("Index Time: %s", info.AvgIndexTime)},
	}

	return genericFormatter.FormatTable(headers, rows)
}

type IndexDetailInfo struct {
	Name         string
	SearchRate   string
	IndexRate    string
	AvgQueryTime string
	AvgIndexTime string
	CheckCount   int
}
