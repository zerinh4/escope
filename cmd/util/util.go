package util

import "github.com/mertbahardogan/escope/internal/util"

func FormatBytes(bytes int64) string {
	return util.FormatBytes(bytes)
}

func IsSystemIndex(name string) bool {
	return util.IsSystemIndex(name)
}
