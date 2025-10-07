package node

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
	"strings"
)

var nodeCmd = &cobra.Command{
	Use:           "node",
	Short:         "Show node information with health summary and JVM heap details",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := elastic.NewClientWrapper(connection.GetClient())
		nodeService := services.NewNodeService(client)

		nodes, err := util.ExecuteWithTimeout(func() ([]models.NodeInfo, error) {
			return nodeService.GetNodesInfo(context.Background())
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Node info check failed: %s\n", constants.MsgTimeoutGeneric)
			} else {
				fmt.Printf("Node info check failed: %v\n", err)
			}
			return
		}

		headers := []string{"Roles", "CPU%", "Mem%", "Heap%", "Disk%", "Free Disk", "Total Disk", "Docs", "Heap Used", "Heap Max", "IP", "Name"}
		rows := make([][]string, 0, len(nodes))

		for _, node := range nodes {
			var filteredRoles []string
			for _, role := range node.Roles {
				if role == "data" || role == "master" {
					filteredRoles = append(filteredRoles, role)
				}
			}

			roles := "-"
			if len(filteredRoles) > 0 {
				roles = strings.Join(filteredRoles, ",")
			}

			name := "-"
			if node.Name != "" {
				name = node.Name
			}

			memPercent := "-"
			if node.MemPercent != "" {
				memPercent = node.MemPercent
			}

			diskPercent := "-"
			if node.DiskPercent != "" {
				diskPercent = node.DiskPercent
			}

			diskTotal := "-"
			if node.DiskTotal != "" {
				diskTotal = node.DiskTotal
			}

			heapUsed := "-"
			if node.HeapUsed != "" {
				heapUsed = node.HeapUsed
			}

			heapMax := "-"
			if node.HeapMax != "" {
				heapMax = node.HeapMax
			}

			docsStr := util.FormatDocsCount(node.Documents)

			if len(roles) > 13 {
				roles = roles[:10] + "..."
			}
			if len(heapUsed) > 10 {
				heapUsed = heapUsed[:7] + "..."
			}
			if len(heapMax) > 10 {
				heapMax = heapMax[:7] + "..."
			}

			row := []string{
				roles,
				node.CPUPercent,
				memPercent,
				node.HeapPercent,
				diskPercent,
				node.DiskAvail,
				diskTotal,
				docsStr,
				heapUsed,
				heapMax,
				node.IP,
				name,
			}
			rows = append(rows, row)
		}

		formatter := ui.NewGenericTableFormatter()
		fmt.Print(formatter.FormatTable(headers, rows))
		fmt.Printf("Total: %d nodes\n", len(nodes))
	},
}

func init() {
	core.RootCmd.AddCommand(nodeCmd)
	nodeCmd.AddCommand(nodeDistCmd)
}
