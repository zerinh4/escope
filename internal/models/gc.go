package models

type GCInfo struct {
	NodeName      string
	EdenSpace     MemorySpace
	SurvivorSpace MemorySpace
	OldGeneration MemorySpace
	TotalHeap     MemorySpace
	YoungGC       GCMetrics
	OldGC         GCMetrics
	FullGC        GCMetrics
	Performance   GCPerformance
}

type MemorySpace struct {
	Used    int64
	Max     int64
	UsedStr string
	MaxStr  string
	Percent float64
}

type GCMetrics struct {
	Count        int64
	TotalTime    int64
	AvgTime      float64
	CountStr     string
	TotalTimeStr string
	AvgTimeStr   string
}

type GCPerformance struct {
	Frequency      float64
	Throughput     float64
	MemoryPressure string
	FrequencyStr   string
	ThroughputStr  string
}

type GCCollection struct {
	Nodes []GCInfo
	Total GCPerformance
}
