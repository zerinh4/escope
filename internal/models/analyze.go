package models

type AnalyzeToken struct {
	Token       string `json:"token"`
	Type        string `json:"type"`
	Position    int    `json:"position"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
}

type AnalyzeResult struct {
	Tokens []AnalyzeToken `json:"tokens"`
}

