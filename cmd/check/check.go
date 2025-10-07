package check

import (
	"context"
	"errors"
	"fmt"
	"github.com/mertbahardogan/escope/cmd/core"
	"github.com/mertbahardogan/escope/internal/connection"
	"github.com/mertbahardogan/escope/internal/constants"
	"github.com/mertbahardogan/escope/internal/elastic"
	"github.com/mertbahardogan/escope/internal/interfaces"
	"github.com/mertbahardogan/escope/internal/models"
	"github.com/mertbahardogan/escope/internal/services"
	"github.com/mertbahardogan/escope/internal/ui"
	"github.com/mertbahardogan/escope/internal/util"
	"github.com/spf13/cobra"
	"time"
)

var (
	duration string
	interval string
)

var checkCmd = &cobra.Command{
	Use:           "check",
	Short:         "Check cluster health metrics",
	Long:          `Check various aspects of your Elasticsearch cluster health and performance`,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := elastic.NewClientWrapper(connection.GetClient())
		checkService := services.NewCheckService(client)
		formatter := ui.NewCheckFormatter()

		if duration != "" {
			runContinuousCheck(context.Background(), client, checkService, formatter)
			return
		}

		runSingleCheck(context.Background(), checkService, formatter)
	},
}

func runSingleCheck(ctx context.Context, checkService services.CheckService, formatter *ui.CheckFormatter) {
	clusterHealth, err := util.ExecuteWithTimeout(func() (*models.ClusterInfo, error) {
		return checkService.GetClusterHealthCheck(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get cluster health: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get cluster health: %v\n", err)
		}
		return
	}

	nodeHealths, err := util.ExecuteWithTimeout(func() ([]models.CheckNodeHealth, error) {
		return checkService.GetNodeHealthCheck(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get node health: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get node health: %v\n", err)
		}
		return
	}

	shardHealth, err := util.ExecuteWithTimeout(func() (*models.ShardHealth, error) {
		return checkService.GetShardHealthCheck(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get shard health: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get shard health: %v\n", err)
		}
		return
	}

	shardWarnings, err := util.ExecuteWithTimeout(func() (*models.ShardWarnings, error) {
		return checkService.GetShardWarningsCheck(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get shard warnings: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get shard warnings: %v\n", err)
		}
		return
	}

	indexHealths, err := util.ExecuteWithTimeout(func() ([]models.IndexHealth, error) {
		return checkService.GetIndexHealthCheck(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get index health: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get index health: %v\n", err)
		}
		return
	}

	resourceUsage, err := util.ExecuteWithTimeout(func() (*models.ResourceUsage, error) {
		return checkService.GetResourceUsageCheck(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get resource usage: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get resource usage: %v\n", err)
		}
		return
	}

	performance, err := util.ExecuteWithTimeout(func() (*models.Performance, error) {
		return checkService.GetPerformanceCheck(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get performance stats: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get performance stats: %v\n", err)
		}
		return
	}

	nodeBreakdown, err := util.ExecuteWithTimeout(func() (*models.NodeBreakdown, error) {
		return checkService.GetNodeBreakdown(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get node breakdown: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get node breakdown: %v\n", err)
		}
		return
	}

	segmentWarnings, err := util.ExecuteWithTimeout(func() (*models.SegmentWarnings, error) {
		return checkService.GetSegmentWarningsCheck(ctx)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Failed to get segment warnings: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Failed to get segment warnings: %v\n", err)
		}
		return
	}

	output := formatter.FormatCheckReport(
		clusterHealth,
		nodeHealths,
		shardHealth,
		shardWarnings,
		indexHealths,
		resourceUsage,
		performance,
		nodeBreakdown,
		segmentWarnings,
	)
	fmt.Print(output)
}

func runContinuousCheck(ctx context.Context, client interfaces.ElasticClient, checkService services.CheckService, formatter *ui.CheckFormatter) {
	durationTime, err := time.ParseDuration(duration)
	if err != nil {
		fmt.Printf("Invalid duration format: %v\n", err)
		fmt.Println("Valid formats: 1m, 5m, 1h, etc.")
		return
	}

	intervalTime := time.Duration(constants.DefaultInterval) * time.Second
	if interval != "" {
		intervalTime, err = time.ParseDuration(interval)
		if err != nil {
			fmt.Printf("Invalid interval format: %v\n", err)
			fmt.Println("Valid formats: 5s, 10s, 1m, etc.")
			return
		}
	}

	monitoringService := services.NewMonitoringService(client)

	result, err := monitoringService.MonitorCluster(ctx, durationTime, intervalTime)
	if err != nil {
		fmt.Printf("Monitoring failed: %v\n", err)
		return
	}

	if result.SampleCount > 0 {
		runSingleCheck(ctx, checkService, formatter)
	} else {
		fmt.Println("No samples collected during monitoring period.")
	}
}

func init() {
	checkCmd.Flags().StringVarP(&duration, "duration", "d", "", "Duration for continuous monitoring (e.g., 1m, 5m, 1h)")
	checkCmd.Flags().StringVarP(&interval, "interval", "i", "",
		"Sampling interval for continuous monitoring (e.g., 5s, 10s, 1m, default: 2s)")

	core.RootCmd.AddCommand(checkCmd)
}
