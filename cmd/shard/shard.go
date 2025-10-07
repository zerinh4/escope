package shard

import (
	"context"
	"errors"
	"escope/cmd/core"
	"escope/cmd/sort"
	"escope/cmd/system"
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

var shardCmd = &cobra.Command{
	Use:                "shard",
	Short:              "Show shard usage and unassigned shard info",
	SilenceErrors:      true,
	DisableSuggestions: true,
	Args:               cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		client := elastic.NewClientWrapper(connection.GetClient())
		shardService := services.NewShardService(client)

		shards, err := util.ExecuteWithTimeout(func() ([]models.ShardInfo, error) {
			return shardService.GetAllShardInfos(context.Background())
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Shard info fetch failed: %s\n", constants.MsgTimeoutGeneric)
			} else {
				fmt.Printf("Shard info fetch failed: %v\n", err)
			}
			return
		}

		var filteredShards []models.ShardInfo
		for _, s := range shards {
			if !util.IsSystemIndex(s.Index) {
				filteredShards = append(filteredShards, s)
			}
		}

		headers := []string{"Shard", "Type", "State", "Size", "Node IP", "Index"}
		rows := make([][]string, 0, len(filteredShards))

		for _, s := range filteredShards {
			typeStr := util.ConvertShardName(s.Prirep)
			row := []string{
				s.Shard,
				typeStr,
				s.State,
				s.Store,
				s.IP,
				s.Index,
			}
			rows = append(rows, row)
		}

		formatter := ui.NewGenericTableFormatter()
		fmt.Print(formatter.FormatTable(headers, rows))
		fmt.Printf("Total: %d shards\n", len(filteredShards))
	},
}

func init() {
	core.RootCmd.AddCommand(shardCmd)

	system.NewSystemCommand(shardCmd, "shard")
	sort.NewSortCommand(shardCmd, "shard")
}
