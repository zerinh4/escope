package models

type ShardInfo struct {
	Index  string
	Shard  string
	Prirep string
	State  string
	Docs   string
	Store  string
	IP     string
	Node   string
}

type ShardStat struct {
	IndexName     string
	PrimaryShards int
	ReplicaShards int
	TotalShards   int
	TotalSize     int64
	NodeCount     int
	Nodes         map[string]bool
}

type ShardDistribution struct {
	NodeDistribution  map[string]int
	IndexDistribution map[string]*ShardStat
}

type ShardWarnings struct {
	UnassignedShards   int
	RelocatingShards   int
	InitializingShards int
	UnbalancedShards   bool
	UnbalancedRatio    float64
	Recommendations    []string
	CriticalIssues     []string
	WarningIssues      []string
}
