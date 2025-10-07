package ui

import (
	"fmt"
	"github.com/mertbahardogan/escope/internal/interfaces"
	"github.com/mertbahardogan/escope/internal/models"
	"sort"
	"strings"
)

type GCFormatter struct{}

func NewGCFormatter() interfaces.GCFormatter {
	return &GCFormatter{}
}

func (f *GCFormatter) FormatGCTable(gcInfos []models.GCInfo) string {
	if len(gcInfos) == 0 {
		return "No GC information found\n"
	}

	sort.Slice(gcInfos, func(i, j int) bool {
		return gcInfos[i].TotalHeap.Percent > gcInfos[j].TotalHeap.Percent
	})

	var totalNodes int
	var highUsage, mediumUsage, lowUsage int

	headers := []string{"Heap Usage %", "Memory Pressure", "Name"}
	rows := make([][]string, 0, len(gcInfos))

	for _, gc := range gcInfos {
		totalNodes++
		if gc.TotalHeap.Percent >= 80 {
			highUsage++
		} else if gc.TotalHeap.Percent >= 60 {
			mediumUsage++
		} else {
			lowUsage++
		}

		heapUsage := fmt.Sprintf("%.1f%%", gc.TotalHeap.Percent)

		row := []string{
			heapUsage,
			gc.Performance.MemoryPressure,
			gc.NodeName,
		}
		rows = append(rows, row)
	}

	formatter := NewGenericTableFormatter()
	var output strings.Builder
	output.WriteString(formatter.FormatTable(headers, rows))

	output.WriteString(fmt.Sprintf("Total Nodes: %d\n", totalNodes))
	output.WriteString(fmt.Sprintf("High Usage (≥80%%): %d (%.1f%%)\n", highUsage, float64(highUsage)/float64(totalNodes)*100))
	output.WriteString(fmt.Sprintf("Medium Usage (60-79%%): %d (%.1f%%)\n", mediumUsage, float64(mediumUsage)/float64(totalNodes)*100))
	output.WriteString(fmt.Sprintf("Low Usage (<60%%): %d (%.1f%%)\n", lowUsage, float64(lowUsage)/float64(totalNodes)*100))

	output.WriteString("\nUse '--name=<node_name>' for detailed information about a specific node.\n")

	return output.String()
}

func (f *GCFormatter) FormatGCDetails(gcInfo models.GCInfo) string {
	var output strings.Builder

	output.WriteString("Heap Memory:\n")
	output.WriteString("  Eden Space:     " + formatMemorySpace(gcInfo.EdenSpace) + "\n")
	output.WriteString("  Survivor Space: " + formatMemorySpace(gcInfo.SurvivorSpace) + "\n")
	output.WriteString("  Old Generation: " + formatMemorySpace(gcInfo.OldGeneration) + "\n")
	output.WriteString("  Total Heap:     " + formatMemorySpace(gcInfo.TotalHeap) + "\n\n")

	output.WriteString("GC Statistics:\n")
	output.WriteString("  Young GC:       " + formatGCMetrics(gcInfo.YoungGC) + "\n")
	output.WriteString("  Old GC:         " + formatGCMetrics(gcInfo.OldGC) + "\n")
	output.WriteString("  Full GC:        " + formatGCMetrics(gcInfo.FullGC) + "\n\n")

	output.WriteString("Performance:\n")
	output.WriteString("  GC Frequency:   " + gcInfo.Performance.FrequencyStr + "\n")
	output.WriteString("  GC Throughput:  " + gcInfo.Performance.ThroughputStr + "\n")
	output.WriteString("  Memory Pressure: " + gcInfo.Performance.MemoryPressure + "\n")

	return output.String()
}

func formatMemorySpace(space models.MemorySpace) string {
	if space.Max == 0 {
		return fmt.Sprintf("%s / %s", space.UsedStr, space.MaxStr)
	}
	return fmt.Sprintf("%s / %s", space.UsedStr, space.MaxStr)
}

func formatGCMetrics(metrics models.GCMetrics) string {
	if metrics.Count == 0 {
		return "0 count / 0ms total (0ms avg)"
	}
	return fmt.Sprintf("%s count / %s total (%s avg)",
		metrics.CountStr, metrics.TotalTimeStr, metrics.AvgTimeStr)
}
