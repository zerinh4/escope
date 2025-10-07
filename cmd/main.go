package main

import (
	_ "escope/cmd/check"
	_ "escope/cmd/cluster"
	_ "escope/cmd/config"
	"escope/cmd/core"
	_ "escope/cmd/index"
	_ "escope/cmd/lucene"
	_ "escope/cmd/node"
	_ "escope/cmd/segments"
	_ "escope/cmd/shard"
	_ "escope/cmd/sort"
	_ "escope/cmd/system"
	_ "escope/cmd/termvectors"
	_ "escope/cmd/version"
)

func main() {
	core.Execute()
}
