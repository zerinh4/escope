package services

import (
	"context"
	"fmt"
	"github.com/mertbahardogan/escope/internal/constants"
	"github.com/mertbahardogan/escope/internal/interfaces"
	"github.com/mertbahardogan/escope/internal/models"
	"github.com/mertbahardogan/escope/internal/util"
	"time"
)

type GCService interface {
	GetGCInfo(ctx context.Context) ([]models.GCInfo, error)
	GetGCInfoForNode(ctx context.Context, nodeName string) (*models.GCInfo, error)
}

type gcService struct {
	client interfaces.ElasticClient
}

func NewGCService(client interfaces.ElasticClient) GCService {
	return &gcService{
		client: client,
	}
}

func (s *gcService) GetGCInfo(ctx context.Context) ([]models.GCInfo, error) {
	statsData, err := s.client.GetNodesStats(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToGetNodeStats, err)
	}

	var gcInfos []models.GCInfo

	if nodes, ok := statsData[constants.NodesField].(map[string]interface{}); ok {
		for nodeID, nodeData := range nodes {
			if node, ok := nodeData.(map[string]interface{}); ok {
				gcInfo, err := s.parseNodeGCInfo(nodeID, node)
				if err != nil {
					continue
				}
				gcInfos = append(gcInfos, *gcInfo)
			}
		}
	}

	return gcInfos, nil
}

func (s *gcService) GetGCInfoForNode(ctx context.Context, nodeName string) (*models.GCInfo, error) {
	statsData, err := s.client.GetNodesStats(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToGetNodeStats, err)
	}

	if nodes, ok := statsData[constants.NodesField].(map[string]interface{}); ok {
		for nodeID, nodeData := range nodes {
			if node, ok := nodeData.(map[string]interface{}); ok {
				if name, ok := node[constants.NameField].(string); ok && name == nodeName {
					return s.parseNodeGCInfo(nodeID, node)
				}
			}
		}
	}

	return nil, fmt.Errorf(constants.ErrNodeNotFound, nodeName)
}

func (s *gcService) parseNodeGCInfo(nodeID string, nodeData map[string]interface{}) (*models.GCInfo, error) {
	gcInfo := &models.GCInfo{
		NodeName: nodeID,
	}

	if name, ok := nodeData[constants.NameField].(string); ok {
		gcInfo.NodeName = name
	}

	if jvm, ok := nodeData[constants.JVMField].(map[string]interface{}); ok {
		if mem, ok := jvm[constants.MemField].(map[string]interface{}); ok {
			s.parseMemorySpaces(gcInfo, mem)
		}
		if collectors, ok := jvm[constants.GCField].(map[string]interface{}); ok {
			s.parseGCMetrics(gcInfo, collectors)
		}
	}
	s.calculatePerformance(gcInfo)

	return gcInfo, nil
}

func (s *gcService) parseMemorySpaces(gcInfo *models.GCInfo, mem map[string]interface{}) {
	if heap, ok := mem[constants.HeapUsedInBytesField].(float64); ok {
		gcInfo.TotalHeap.Used = int64(heap)
		gcInfo.TotalHeap.UsedStr = util.FormatBytes(int64(heap))
	} else {
		gcInfo.TotalHeap.UsedStr = constants.DashString
	}

	if heapMax, ok := mem[constants.HeapMaxInBytesField].(float64); ok {
		gcInfo.TotalHeap.Max = int64(heapMax)
		gcInfo.TotalHeap.MaxStr = util.FormatBytes(int64(heapMax))
	} else {
		gcInfo.TotalHeap.MaxStr = constants.DashString
	}
	if gcInfo.TotalHeap.Max > 0 {
		gcInfo.TotalHeap.Percent = util.CalculatePercentage(gcInfo.TotalHeap.Used, gcInfo.TotalHeap.Max)
	}

	if pools, ok := mem[constants.PoolsField].(map[string]interface{}); ok {
		s.parseHeapPools(gcInfo, pools)
	} else {
		gcInfo.EdenSpace.UsedStr = constants.DashString
		gcInfo.EdenSpace.MaxStr = constants.DashString
		gcInfo.SurvivorSpace.UsedStr = constants.DashString
		gcInfo.SurvivorSpace.MaxStr = constants.DashString
		gcInfo.OldGeneration.UsedStr = constants.DashString
		gcInfo.OldGeneration.MaxStr = constants.DashString
	}
}

func (s *gcService) parseHeapPools(gcInfo *models.GCInfo, pools map[string]interface{}) {
	for poolName, poolData := range pools {
		if pool, ok := poolData.(map[string]interface{}); ok {
			if isEdenPool(poolName) {
				s.parsePoolMetrics(&gcInfo.EdenSpace, pool)
			} else if isSurvivorPool(poolName) {
				s.parsePoolMetrics(&gcInfo.SurvivorSpace, pool)
			} else if isOldPool(poolName) {
				s.parsePoolMetrics(&gcInfo.OldGeneration, pool)
			}
		}
	}
}

func (s *gcService) parsePoolMetrics(space *models.MemorySpace, pool map[string]interface{}) {
	if used, ok := pool[constants.UsedInBytesField].(float64); ok {
		space.Used = int64(used)
		space.UsedStr = util.FormatBytes(int64(used))
	}

	if max, ok := pool[constants.HeapMaxInBytesField].(float64); ok {
		space.Max = int64(max)
		space.MaxStr = util.FormatBytes(int64(max))
	}

	if space.Max > 0 {
		space.Percent = util.CalculatePercentage(space.Used, space.Max)
	}
}

func isEdenPool(poolName string) bool {
	return poolName == constants.GCYoung
}

func isSurvivorPool(poolName string) bool {
	return poolName == constants.GCSurvivor
}

func isOldPool(poolName string) bool {
	return poolName == constants.GCOld
}

func (s *gcService) parseGCMetrics(gcInfo *models.GCInfo, collectors map[string]interface{}) {
	if collectorsData, ok := collectors[constants.CollectorsField].(map[string]interface{}); ok {
		for collectorName, collector := range collectorsData {
			if collectorData, ok := collector.(map[string]interface{}); ok {
				if isYoungGC(collectorName) {
					s.parseCollectorMetrics(&gcInfo.YoungGC, collectorData)
				} else if isOldGC(collectorName) {
					s.parseCollectorMetrics(&gcInfo.OldGC, collectorData)
				} else if isFullGC(collectorName) {
					s.parseCollectorMetrics(&gcInfo.FullGC, collectorData)
				}
			}
		}
	}

	if gcInfo.YoungGC.CountStr == constants.EmptyString {
		gcInfo.YoungGC.CountStr = constants.DashString
		gcInfo.YoungGC.TotalTimeStr = constants.DashString
	}
	if gcInfo.OldGC.CountStr == constants.EmptyString {
		gcInfo.OldGC.CountStr = constants.DashString
		gcInfo.OldGC.TotalTimeStr = constants.DashString
	}
	if gcInfo.FullGC.CountStr == constants.EmptyString {
		gcInfo.FullGC.CountStr = constants.DashString
		gcInfo.FullGC.TotalTimeStr = constants.DashString
	}
}

func (s *gcService) parseCollectorMetrics(metrics *models.GCMetrics, collectorData map[string]interface{}) {
	if count, ok := collectorData[constants.CollectionCountField].(float64); ok {
		metrics.Count = int64(count)
		metrics.CountStr = util.FormatDocsCount(int64(count))
	}

	if time, ok := collectorData[constants.CollectionTimeInMillisField].(float64); ok {
		metrics.TotalTime = int64(time)
		metrics.TotalTimeStr = formatDuration(int64(time))
	}

	if metrics.Count > 0 {
		metrics.AvgTime = float64(metrics.TotalTime) / float64(metrics.Count)
		metrics.AvgTimeStr = formatDuration(int64(metrics.AvgTime))
	}
}

func (s *gcService) calculatePerformance(gcInfo *models.GCInfo) {
	uptime := time.Hour.Seconds()
	totalGC := gcInfo.YoungGC.Count + gcInfo.OldGC.Count + gcInfo.FullGC.Count

	if uptime > 0 {
		gcInfo.Performance.Frequency = float64(totalGC) / uptime
		gcInfo.Performance.FrequencyStr = fmt.Sprintf(constants.GCFreqFormat, gcInfo.Performance.Frequency)
	}

	totalGCTime := gcInfo.YoungGC.TotalTime + gcInfo.OldGC.TotalTime + gcInfo.FullGC.TotalTime
	if uptime > 0 {
		throughput := (uptime - float64(totalGCTime)/constants.MillisecondsToSeconds) / uptime * constants.HundredMultiplier
		if throughput < 0 {
			throughput = 0
		}
		gcInfo.Performance.Throughput = throughput
		gcInfo.Performance.ThroughputStr = fmt.Sprintf(constants.ThroughputFormat, gcInfo.Performance.Throughput)
	}

	gcInfo.Performance.MemoryPressure = calculateMemoryPressure(gcInfo.TotalHeap.Percent)
}

func isYoungGC(collectorName string) bool {
	return collectorName == constants.GCYoung
}

func isOldGC(collectorName string) bool {
	return collectorName == constants.GCOld
}

func isFullGC(collectorName string) bool {
	return collectorName == constants.GCG1Concurrent
}

func calculateMemoryPressure(heapPercent float64) string {
	if heapPercent < constants.LowMemoryPressure {
		return constants.MemoryPressureLow
	} else if heapPercent < constants.MediumMemoryPressure {
		return constants.MemoryPressureMedium
	} else {
		return constants.MemoryPressureHigh
	}
}

func formatDuration(milliseconds int64) string {
	if milliseconds < constants.MillisecondsToSeconds {
		return fmt.Sprintf(constants.MSFormat, milliseconds)
	} else {
		return fmt.Sprintf(constants.TimeFormatS, float64(milliseconds)/constants.MillisecondsToSeconds)
	}
}
