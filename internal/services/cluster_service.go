package services

import (
	"context"
	"fmt"
	"github.com/mertbahardogan/escope/internal/constants"
	"github.com/mertbahardogan/escope/internal/interfaces"
	"github.com/mertbahardogan/escope/internal/models"
	"github.com/mertbahardogan/escope/internal/util"
)

type ClusterService interface {
	GetClusterHealth(ctx context.Context) (*models.ClusterInfo, error)
	GetClusterStats(ctx context.Context) (*models.ClusterStats, error)
}

type clusterService struct {
	client       interfaces.ElasticClient
	nodeService  NodeService
	indexService IndexService
}

func NewClusterService(client interfaces.ElasticClient) ClusterService {
	return &clusterService{
		client:       client,
		nodeService:  NewNodeService(client),
		indexService: NewIndexService(client),
	}
}

func (s *clusterService) GetClusterHealth(ctx context.Context) (*models.ClusterInfo, error) {
	healthData, err := s.client.GetClusterHealth(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrClusterHealthRequestFailed2, err)
	}

	health := &models.ClusterInfo{}

	if clusterName, ok := healthData[constants.ClusterNameField].(string); ok {
		health.ClusterName = clusterName
	}

	if status, ok := healthData[constants.StatusField].(string); ok {
		health.Status = status
	}

	if numberOfNodes, ok := healthData[constants.NumberOfNodesField].(float64); ok {
		health.NumberOfNodes = int(numberOfNodes)
	}

	if activePrimaryShards, ok := healthData[constants.ActivePrimaryShardsField].(float64); ok {
		health.ActivePrimaryShards = int(activePrimaryShards)
	}

	if activeShards, ok := healthData[constants.ActiveShardsField].(float64); ok {
		health.ActiveShards = int(activeShards)
	}

	if unassignedShards, ok := healthData[constants.UnassignedShardsField].(float64); ok {
		health.UnassignedShards = int(unassignedShards)
	}

	if relocatingShards, ok := healthData[constants.RelocatingShardsField].(float64); ok {
		health.RelocatingShards = int(relocatingShards)
	}

	if initializingShards, ok := healthData[constants.InitializingShardsField].(float64); ok {
		health.InitializingShards = int(initializingShards)
	}

	if delayedUnassignedShards, ok := healthData[constants.DelayedUnassignedShardsField].(float64); ok {
		health.DelayedUnassignedShards = int(delayedUnassignedShards)
	}

	if numberOfPendingTasks, ok := healthData[constants.NumberOfPendingTasksField].(float64); ok {
		health.NumberOfPendingTasks = int(numberOfPendingTasks)
	}

	if numberOfInFlightFetch, ok := healthData[constants.NumberOfInFlightFetchField].(float64); ok {
		health.NumberOfInFlightFetch = int(numberOfInFlightFetch)
	}

	if taskMaxWaitingInQueueMillis, ok := healthData[constants.TaskMaxWaitingInQueueField].(float64); ok {
		health.TaskMaxWaitingInQueueMillis = int(taskMaxWaitingInQueueMillis)
	}

	if activeShardsPercentAsNumber, ok := healthData[constants.ActiveShardsPercentField].(float64); ok {
		health.ActiveShardsPercentAsNumber = activeShardsPercentAsNumber
	}

	if timedOut, ok := healthData[constants.TimedOutField].(bool); ok {
		health.TimedOut = timedOut
	}

	return health, nil
}

func (s *clusterService) GetClusterStats(ctx context.Context) (*models.ClusterStats, error) {
	clusterStatsData, err := s.client.GetClusterStats(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrClusterStatsRequestFailed2, err)
	}

	healthData, err := s.client.GetClusterHealth(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrClusterHealthRequestFailed2, err)
	}

	nodesStatsData, err := s.client.GetNodesStats(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrNodeStatsRequestFailed2, err)
	}

	stats := &models.ClusterStats{}

	s.parseBasicInfo(healthData, stats)
	s.parseNodeInfo(clusterStatsData, stats)
	s.parseIndicesData(clusterStatsData, stats)
	s.parseVersionInfo(clusterStatsData, stats)
	s.parseResourceUsage(nodesStatsData, stats)
	s.calculatePercentages(stats)

	return stats, nil
}

func (s *clusterService) parseBasicInfo(healthData map[string]interface{}, stats *models.ClusterStats) {
	if clusterName, ok := healthData[constants.ClusterNameField].(string); ok {
		stats.ClusterName = clusterName
	}
	if status, ok := healthData[constants.StatusField].(string); ok {
		stats.Status = status
	}
}

func (s *clusterService) parseNodeInfo(clusterStatsData map[string]interface{}, stats *models.ClusterStats) {
	nodes, ok := clusterStatsData[constants.NodesField].(map[string]interface{})
	if !ok {
		return
	}

	if count, ok := nodes[constants.CountField].(map[string]interface{}); ok {
		if total, ok := count[constants.TotalField].(float64); ok {
			stats.TotalNodes = int(total)
		}
		if data, ok := count[constants.NodeRoleData].(float64); ok {
			stats.DataNodes = int(data)
		}
		if master, ok := count[constants.NodeRoleMaster].(float64); ok {
			stats.MasterNodes = int(master)
		}
		if ingest, ok := count[constants.NodeRoleIngest].(float64); ok {
			stats.IngestNodes = int(ingest)
		}
		if coordinating, ok := count["coordinating_only"].(float64); ok {
			stats.CoordinatingNodes = int(coordinating)
		}
	}

	if jvm, ok := nodes[constants.JVMField].(map[string]interface{}); ok {
		if versions, ok := jvm["versions"].([]interface{}); ok {
			for _, v := range versions {
				if versionMap, ok := v.(map[string]interface{}); ok {
					if version, ok := versionMap["version"].(string); ok {
						stats.JVMVersions = append(stats.JVMVersions, version)
					}
				}
			}
		}
	}
}

func (s *clusterService) parseIndicesData(clusterStatsData map[string]interface{}, stats *models.ClusterStats) {
	indices, ok := clusterStatsData["indices"].(map[string]interface{})
	if !ok {
		return
	}

	if count, ok := indices["count"].(float64); ok {
		stats.TotalIndices = int(count)
	}

	if docs, ok := indices["docs"].(map[string]interface{}); ok {
		if count, ok := docs["count"].(float64); ok {
			stats.TotalDocuments = int64(count)
		}
	}

	if shards, ok := indices["shards"].(map[string]interface{}); ok {
		if total, ok := shards["total"].(float64); ok {
			stats.TotalShards = int(total)
		}
		if primaries, ok := shards["primaries"].(float64); ok {
			stats.PrimaryShards = int(primaries)
		}
	}

	if store, ok := indices["store"].(map[string]interface{}); ok {
		if sizeInBytes, ok := store["size_in_bytes"].(float64); ok {
			stats.UsedDiskBytes = int64(sizeInBytes)
			if stats.TotalShards > 0 {
				stats.AvgShardSizeGB = float64(sizeInBytes) / float64(stats.TotalShards) / (1024 * 1024 * 1024)
			}
		}
	}
}

func (s *clusterService) parseVersionInfo(clusterStatsData map[string]interface{}, stats *models.ClusterStats) {
	nodes, ok := clusterStatsData["nodes"].(map[string]interface{})
	if !ok {
		return
	}

	if versions, ok := nodes["versions"].([]interface{}); ok {
		if len(versions) > 0 {
			if version, ok := versions[0].(string); ok {
				stats.ESVersion = version
			}
		}
	}
}

func (s *clusterService) parseResourceUsage(nodesStatsData map[string]interface{}, stats *models.ClusterStats) {
	nodes, ok := nodesStatsData["nodes"].(map[string]interface{})
	if !ok {
		return
	}

	for _, nodeData := range nodes {
		node, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}

		s.parseJVMMemory(node, stats)
		s.parseSystemMemory(node, stats)
		s.parseDiskUsage(node, stats)
	}
}

func (s *clusterService) parseJVMMemory(node map[string]interface{}, stats *models.ClusterStats) {
	jvm, ok := node["jvm"].(map[string]interface{})
	if !ok {
		return
	}

	mem, ok := jvm["mem"].(map[string]interface{})
	if !ok {
		return
	}

	if heapUsed, ok := mem["heap_used_in_bytes"].(float64); ok {
		stats.UsedHeapBytes += int64(heapUsed)
	}
	if heapMax, ok := mem["heap_max_in_bytes"].(float64); ok {
		stats.TotalHeapBytes += int64(heapMax)
	}
}

func (s *clusterService) parseSystemMemory(node map[string]interface{}, stats *models.ClusterStats) {
	os, ok := node["os"].(map[string]interface{})
	if !ok {
		return
	}

	mem, ok := os["mem"].(map[string]interface{})
	if !ok {
		return
	}

	if totalBytes, ok := mem["total_in_bytes"].(float64); ok {
		stats.TotalMemoryBytes += int64(totalBytes)
	}
	if usedBytes, ok := mem["used_in_bytes"].(float64); ok {
		stats.UsedMemoryBytes += int64(usedBytes)
	}
}

func (s *clusterService) parseDiskUsage(node map[string]interface{}, stats *models.ClusterStats) {
	fs, ok := node["fs"].(map[string]interface{})
	if !ok {
		return
	}

	total, ok := fs["total"].(map[string]interface{})
	if !ok {
		return
	}

	if totalBytes, ok := total["total_in_bytes"].(float64); ok {
		stats.TotalDiskBytes += int64(totalBytes)
	}
	if availBytes, ok := total["available_in_bytes"].(float64); ok {
		stats.AvailableDiskBytes += int64(availBytes)
	}
}

func (s *clusterService) calculatePercentages(stats *models.ClusterStats) {
	stats.HeapUsagePercent = util.CalculatePercentage(stats.UsedHeapBytes, stats.TotalHeapBytes)
	stats.MemoryUsagePercent = util.CalculatePercentage(stats.UsedMemoryBytes, stats.TotalMemoryBytes)

	usedDisk := stats.TotalDiskBytes - stats.AvailableDiskBytes
	stats.DiskUsagePercent = util.CalculatePercentage(usedDisk, stats.TotalDiskBytes)
}
