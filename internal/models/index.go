package models

import (
	"sync"
	"time"
)

type IndexInfo struct {
	Alias     string
	Name      string
	Health    string
	Status    string
	DocsCount string
	StoreSize string
	Primary   string
	Replica   string
}

type LuceneStats struct {
	IndexName                string
	SegmentCount             int
	SegmentMemory            string
	SegmentMemoryBytes       int64
	IndexMemory              string
	IndexMemoryBytes         int64
	TermsMemory              string
	TermsMemoryBytes         int64
	StoredMemory             string
	StoredMemoryBytes        int64
	DocValuesMemory          string
	DocValuesMemoryBytes     int64
	PointsMemory             string
	PointsMemoryBytes        int64
	NormsMemory              string
	NormsMemoryBytes         int64
	FixedBitSetMemory        string
	FixedBitSetMemoryBytes   int64
	VersionMapMemory         string
	VersionMapMemoryBytes    int64
	MaxUnsafeAutoIDTimestamp int64
}

type IndexDetailInfo struct {
	Name         string
	SearchRate   string
	IndexRate    string
	AvgQueryTime string
	AvgIndexTime string
}

type IndexStatsSnapshot struct {
	IndexName  string
	QueryTotal int64
	QueryTime  int64
	IndexTotal int64
	IndexTime  int64
	Timestamp  time.Time
}

type IndexStatsCache struct {
	mu        sync.RWMutex
	snapshots map[string]*IndexStatsSnapshot
}

func NewIndexStatsCache() *IndexStatsCache {
	return &IndexStatsCache{
		snapshots: make(map[string]*IndexStatsSnapshot),
	}
}

func (c *IndexStatsCache) GetSnapshot(indexName string) (*IndexStatsSnapshot, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snapshot, exists := c.snapshots[indexName]
	return snapshot, exists
}

func (c *IndexStatsCache) SetSnapshot(snapshot *IndexStatsSnapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snapshots[snapshot.IndexName] = snapshot
}
