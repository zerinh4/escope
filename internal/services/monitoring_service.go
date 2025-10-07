package services

import (
	"context"
	"escope/internal/constants"
	"escope/internal/interfaces"
	"escope/internal/models"
	"escope/internal/util"
	"fmt"
	"time"
)

type MonitoringService struct {
	client interfaces.ElasticClient
}

type MonitoringResult struct {
	Duration           time.Duration
	SampleCount        int
	ClusterHealthTrend []models.ClusterInfo
	NodeHealthTrend    []models.CheckNodeHealth
	ShardHealthTrend   []models.ShardHealth
	IndexHealthTrend   []models.IndexHealth
	ResourceTrend      []models.ResourceUsage
	PerformanceTrend   []models.Performance
	Issues             []MonitoringIssue
	Recommendations    []string
}

type MonitoringIssue struct {
	Type        string
	Severity    string
	Description string
	Occurrences int
	FirstSeen   time.Time
	LastSeen    time.Time
}

func NewMonitoringService(client interfaces.ElasticClient) *MonitoringService {
	return &MonitoringService{
		client: client,
	}
}

func (s *MonitoringService) MonitorCluster(ctx context.Context, duration time.Duration,
	interval time.Duration) (*MonitoringResult, error) {
	if interval > duration {
		interval = duration / constants.DefaultInterval
		if interval < time.Duration(constants.MinInterval)*time.Second {
			interval = time.Duration(constants.MinInterval) * time.Second
		}
		fmt.Printf("Warning: Interval adjusted to %v (duration was too short)\n", interval)
	}

	result := &MonitoringResult{
		Duration:    duration,
		SampleCount: 0,
	}

	checkService := NewCheckService(s.client)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	endTime := time.Now().Add(duration)

	fmt.Printf("Starting cluster monitoring for %v (sampling every %v)\n", duration, interval)
	fmt.Printf("Monitoring will complete at %s\n\n", endTime.Format("15:04:05"))

	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-ticker.C:
			if time.Now().After(endTime) {
				goto monitoringComplete
			}

			clusterHealth, err := checkService.GetClusterHealthCheck(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to get cluster health at %s: %v\n", time.Now().Format("15:04:05"), err)
				continue
			}

			nodeHealths, err := checkService.GetNodeHealthCheck(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to get node health at %s: %v\n", time.Now().Format("15:04:05"), err)
				continue
			}

			shardHealth, err := checkService.GetShardHealthCheck(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to get shard health at %s: %v\n", time.Now().Format("15:04:05"), err)
				continue
			}

			indexHealths, err := checkService.GetIndexHealthCheck(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to get index health at %s: %v\n", time.Now().Format("15:04:05"), err)
				continue
			}

			resourceUsage, err := checkService.GetResourceUsageCheck(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to get resource usage at %s: %v\n", time.Now().Format("15:04:05"), err)
				continue
			}

			performance, err := checkService.GetPerformanceCheck(ctx)
			if err != nil {
				fmt.Printf("Warning: Failed to get performance stats at %s: %v\n", time.Now().Format("15:04:05"), err)
				continue
			}

			result.ClusterHealthTrend = append(result.ClusterHealthTrend, *clusterHealth)
			result.NodeHealthTrend = append(result.NodeHealthTrend, nodeHealths...)
			result.ShardHealthTrend = append(result.ShardHealthTrend, *shardHealth)
			result.IndexHealthTrend = append(result.IndexHealthTrend, indexHealths...)
			if resourceUsage != nil {
				result.ResourceTrend = append(result.ResourceTrend, *resourceUsage)
			}
			if performance != nil {
				result.PerformanceTrend = append(result.PerformanceTrend, *performance)
			}

			result.SampleCount++

			fmt.Printf("Sample %d collected at %s\n",
				result.SampleCount, time.Now().Format("15:04:05"))

			if clusterHealth.Status == constants.HealthRed {
				fmt.Printf("!!! CRITICAL: Cluster status is RED at %s\n", time.Now().Format("15:04:05"))
			}
			if clusterHealth.UnassignedShards > 0 {
				fmt.Printf("!!! WARNING: %d unassigned shards detected at %s\n",
					clusterHealth.UnassignedShards, time.Now().Format("15:04:05"))
			}
		}
	}

monitoringComplete:

	fmt.Printf("\nMonitoring completed. Collected %d samples over %v\n", result.SampleCount, duration)

	s.analyzeTrends(result)

	return result, nil
}

func (s *MonitoringService) analyzeTrends(result *MonitoringResult) {
	s.analyzeClusterHealthTrend(result)

	s.analyzeResourceTrends(result)

	s.analyzePerformanceTrends(result)

	s.generateRecommendations(result)
}

func (s *MonitoringService) analyzeClusterHealthTrend(result *MonitoringResult) {
	if len(result.ClusterHealthTrend) == 0 {
		return
	}

	redCount := 0
	yellowCount := 0
	greenCount := 0
	unassignedShardsTotal := 0

	for _, health := range result.ClusterHealthTrend {
		switch health.Status {
		case constants.HealthRed:
			redCount++
		case constants.HealthYellow:
			yellowCount++
		case constants.HealthGreen:
			greenCount++
		}
		unassignedShardsTotal += health.UnassignedShards
	}

	if redCount > 0 {
		result.Issues = append(result.Issues, MonitoringIssue{
			Type:        "Cluster Status",
			Severity:    constants.SeverityCritical,
			Description: fmt.Sprintf("Cluster was RED in %d/%d samples", redCount, result.SampleCount),
			Occurrences: redCount,
			FirstSeen:   result.ClusterHealthTrend[0].Timestamp,
			LastSeen:    result.ClusterHealthTrend[len(result.ClusterHealthTrend)-1].Timestamp,
		})
	}

	if yellowCount > result.SampleCount/2 {
		result.Issues = append(result.Issues, MonitoringIssue{
			Type:        "Cluster Status",
			Severity:    constants.SeverityWarning,
			Description: fmt.Sprintf("Cluster was YELLOW in %d/%d samples (>50%%)", yellowCount, result.SampleCount),
			Occurrences: yellowCount,
			FirstSeen:   result.ClusterHealthTrend[0].Timestamp,
			LastSeen:    result.ClusterHealthTrend[len(result.ClusterHealthTrend)-1].Timestamp,
		})
	}

	if unassignedShardsTotal > 0 {
		result.Issues = append(result.Issues, MonitoringIssue{
			Type:        "Shard Assignment",
			Severity:    constants.SeverityWarning,
			Description: fmt.Sprintf("Total unassigned shards across samples: %d", unassignedShardsTotal),
			Occurrences: unassignedShardsTotal,
			FirstSeen:   result.ClusterHealthTrend[0].Timestamp,
			LastSeen:    result.ClusterHealthTrend[len(result.ClusterHealthTrend)-1].Timestamp,
		})
	}
}

func (s *MonitoringService) analyzeResourceTrends(result *MonitoringResult) {
	if len(result.ResourceTrend) == 0 {
		return
	}

	highHeapSamples := 0
	maxHeapUsage := 0.0
	for _, resource := range result.ResourceTrend {
		if resource.HeapUsage > 80 {
			highHeapSamples++
		}
		if resource.HeapUsage > maxHeapUsage {
			maxHeapUsage = resource.HeapUsage
		}
	}

	if highHeapSamples > result.SampleCount/3 {
		result.Issues = append(result.Issues, MonitoringIssue{
			Type:     "Resource Usage",
			Severity: constants.SeverityWarning,
			Description: fmt.Sprintf("High heap usage (>80%%) in %d/%d samples, max: %.1f%%",
				highHeapSamples, result.SampleCount, maxHeapUsage),
			Occurrences: highHeapSamples,
			FirstSeen:   result.ResourceTrend[0].Timestamp,
			LastSeen:    result.ResourceTrend[len(result.ResourceTrend)-1].Timestamp,
		})
	}

	highDiskSamples := 0
	for _, resource := range result.ResourceTrend {
		diskUsed := resource.DiskTotal - resource.DiskAvailable
		diskPercent := util.CalculatePercentage(diskUsed, resource.DiskTotal)
		if diskPercent > 85 {
			highDiskSamples++
		}
	}

	if highDiskSamples > 0 {
		result.Issues = append(result.Issues, MonitoringIssue{
			Type:     "Resource Usage",
			Severity: constants.SeverityWarning,
			Description: fmt.Sprintf("High disk usage (>85%%) in %d/%d samples",
				highDiskSamples, result.SampleCount),
			Occurrences: highDiskSamples,
			FirstSeen:   result.ResourceTrend[0].Timestamp,
			LastSeen:    result.ResourceTrend[len(result.ResourceTrend)-1].Timestamp,
		})
	}
}

func (s *MonitoringService) analyzePerformanceTrends(result *MonitoringResult) {
	if len(result.PerformanceTrend) == 0 {
		return
	}
	slowIndexSamples := 0
	totalIndexTime := int64(0)
	totalIndexCount := int64(0)

	for _, perf := range result.PerformanceTrend {
		if perf.IndexTotal > 0 {
			avgIndexTime := float64(perf.IndexTimeInMillis) / float64(perf.IndexTotal)
			if avgIndexTime > 100 {
				slowIndexSamples++
			}
			totalIndexTime += perf.IndexTimeInMillis
			totalIndexCount += perf.IndexTotal
		}
	}

	if slowIndexSamples > 0 && totalIndexCount > 0 {
		overallAvgIndexTime := float64(totalIndexTime) / float64(totalIndexCount)
		result.Issues = append(result.Issues, MonitoringIssue{
			Type:     "Performance",
			Severity: constants.SeverityWarning,
			Description: fmt.Sprintf("Slow indexing (>100ms avg) in %d/%d samples, overall avg: %.1fms",
				slowIndexSamples, result.SampleCount, overallAvgIndexTime),
			Occurrences: slowIndexSamples,
			FirstSeen:   result.PerformanceTrend[0].Timestamp,
			LastSeen:    result.PerformanceTrend[len(result.PerformanceTrend)-1].Timestamp,
		})
	}

	slowSearchSamples := 0
	totalSearchTime := int64(0)
	totalSearchCount := int64(0)

	for _, perf := range result.PerformanceTrend {
		if perf.QueryTotal > 0 {
			avgSearchTime := float64(perf.QueryTimeInMillis) / float64(perf.QueryTotal)
			if avgSearchTime > 50 {
				slowSearchSamples++
			}
			totalSearchTime += perf.QueryTimeInMillis
			totalSearchCount += perf.QueryTotal
		}
	}

	if slowSearchSamples > 0 && totalSearchCount > 0 {
		overallAvgSearchTime := float64(totalSearchTime) / float64(totalSearchCount)
		result.Issues = append(result.Issues, MonitoringIssue{
			Type:     "Performance",
			Severity: "WARNING",
			Description: fmt.Sprintf("Slow searches (>50ms avg) in %d/%d samples, overall avg: %.1fms",
				slowSearchSamples, result.SampleCount, overallAvgSearchTime),
			Occurrences: slowSearchSamples,
			FirstSeen:   result.PerformanceTrend[0].Timestamp,
			LastSeen:    result.PerformanceTrend[len(result.PerformanceTrend)-1].Timestamp,
		})
	}
}

func (s *MonitoringService) generateRecommendations(result *MonitoringResult) {
	for _, issue := range result.Issues {
		switch issue.Type {
		case "Cluster Status":
			if issue.Severity == "CRITICAL" {
				result.Recommendations = append(result.Recommendations,
					"Investigate cluster RED status immediately - check node failures and shard allocation")
			} else if issue.Severity == "WARNING" {
				result.Recommendations = append(result.Recommendations,
					"Monitor cluster status - consider rebalancing if YELLOW persists")
			}
		case "Shard Assignment":
			result.Recommendations = append(result.Recommendations,
				"Review shard allocation and node capacity - unassigned shards indicate resource constraints")
		case "Resource Usage":
			if issue.Description[:4] == "High" {
				result.Recommendations = append(result.Recommendations,
					"Review resource allocation and consider scaling up nodes or optimizing usage")
			}
		case "Performance":
			if issue.Description[:4] == "Slow" {
				result.Recommendations = append(result.Recommendations,
					"Investigate performance bottlenecks - check indexing patterns and query optimization")
			}
		}
	}

	if len(result.Issues) == 0 {
		result.Recommendations = append(result.Recommendations,
			"Cluster is performing well - continue monitoring for trends")
	} else {
		result.Recommendations = append(result.Recommendations,
			"Schedule regular health checks and monitor identified issues")
	}
}
