package models

type NodeInfo struct {
	Name        string
	IP          string
	Roles       []string
	CPUPercent  string
	MemPercent  string
	HeapPercent string
	DiskPercent string
	DiskAvail   string
	DiskTotal   string
	Documents   int64
	HeapUsed    string
	HeapMax     string
}

type NodeStat struct {
	NodeIP        string
	PrimaryShards int
	ReplicaShards int
	TotalShards   int
	TotalSize     int64
	IndexCount    int
}

type BalanceAnalysis struct {
	MostLoadedNode  string
	LeastLoadedNode string
	MaxShards       int
	MinShards       int
	BalanceRatio    float64
	IsBalanced      bool
	Recommendation  string
}

type NodeHealth struct {
	NodeName    string
	Status      string
	CPUUsage    float64
	MemoryUsage float64
	HeapUsage   float64
	DiskUsage   float64
	IsHealthy   bool
	Issues      []string
}
