package models

type TermInfo struct {
	Field    string `json:"field"`
	Term     string `json:"term"`
	TermFreq int    `json:"term_freq"`
}
