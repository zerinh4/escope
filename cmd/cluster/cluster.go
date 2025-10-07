package cluster

import (
	"context"
	"errors"
	"escope/cmd/core"
	"escope/internal/connection"
	"escope/internal/constants"
	"escope/internal/elastic"
	"escope/internal/models"
	"escope/internal/services"
	"escope/internal/ui"
	"escope/internal/util"
	"fmt"
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:           "cluster",
	Short:         "Show detailed cluster health information",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := elastic.NewClientWrapper(connection.GetClient())
		clusterService := services.NewClusterService(client)
		formatter := ui.NewClusterFormatter()

		clusterStats, err := util.ExecuteWithTimeout(func() (*models.ClusterStats, error) {
			return clusterService.GetClusterStats(context.Background())
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Failed to get cluster health: %s\n", constants.MsgTimeoutGeneric)
			} else {
				fmt.Printf("Failed to get cluster health: %v\n", err)
			}
			return
		}

		output := formatter.FormatClusterStats(clusterStats)
		fmt.Print(output)
	},
}

func init() {
	core.RootCmd.AddCommand(clusterCmd)
}
