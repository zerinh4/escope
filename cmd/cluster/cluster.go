package cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/mertbahardogan/escope/cmd/core"
	"github.com/mertbahardogan/escope/internal/connection"
	"github.com/mertbahardogan/escope/internal/constants"
	"github.com/mertbahardogan/escope/internal/elastic"
	"github.com/mertbahardogan/escope/internal/models"
	"github.com/mertbahardogan/escope/internal/services"
	"github.com/mertbahardogan/escope/internal/ui"
	"github.com/mertbahardogan/escope/internal/util"
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
