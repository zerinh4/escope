package models

import (
	"fmt"
	"strings"
	"time"
)

type ClusterInfo struct {
	Timestamp                   time.Time
	ClusterName                 string
	Status                      string
	NumberOfNodes               int
	ActivePrimaryShards         int
	ActiveShards                int
	UnassignedShards            int
	RelocatingShards            int
	InitializingShards          int
	DelayedUnassignedShards     int
	NumberOfPendingTasks        int
	NumberOfInFlightFetch       int
	TaskMaxWaitingInQueueMillis int
	ActiveShardsPercentAsNumber float64
	TimedOut                    bool
}

type NodeBreakdown struct {
	DataNodes         int
	MasterNodes       int
	IngestNodes       int
	CoordinatingNodes int
}

func (n *NodeBreakdown) String() string {
	var parts []string
	if n.DataNodes > 0 {
		parts = append(parts, fmt.Sprintf("Data: %d", n.DataNodes))
	}
	if n.MasterNodes > 0 {
		parts = append(parts, fmt.Sprintf("Master: %d", n.MasterNodes))
	}
	if n.IngestNodes > 0 {
		parts = append(parts, fmt.Sprintf("Ingest: %d", n.IngestNodes))
	}
	if n.CoordinatingNodes > 0 {
		parts = append(parts, fmt.Sprintf("Coord: %d", n.CoordinatingNodes))
	}
	if len(parts) == 0 {
		return "unknown"
	}
	return strings.Join(parts, ", ")
}

type ClusterStats struct {
	ClusterName string
	Status      string

	// Node breakdown
	TotalNodes        int
	DataNodes         int
	MasterNodes       int
	IngestNodes       int
	CoordinatingNodes int

	// Storage
	TotalDiskBytes     int64
	UsedDiskBytes      int64
	AvailableDiskBytes int64
	DiskUsagePercent   float64

	// Memory
	TotalHeapBytes     int64
	UsedHeapBytes      int64
	HeapUsagePercent   float64
	TotalMemoryBytes   int64
	UsedMemoryBytes    int64
	MemoryUsagePercent float64

	// Data
	TotalDocuments int64
	TotalIndices   int
	PrimaryShards  int
	TotalShards    int
	AvgShardSizeGB float64

	// Version info
	ESVersion   string
	JVMVersions []string
}

func (c *ClusterStats) GetNodeBreakdown() string {
	breakdown := &NodeBreakdown{
		DataNodes:         c.DataNodes,
		MasterNodes:       c.MasterNodes,
		IngestNodes:       c.IngestNodes,
		CoordinatingNodes: c.CoordinatingNodes,
	}
	return breakdown.String()
}
