package model

type CodeStats struct {
	FilesCount        int    `json:"Count"`
	LinesCount        int    `json:"Lines"`
	BlankLinesCount   int    `json:"Blank"`
	CodeLinesCount    int    `json:"Code"`
	CommentLinesCount int    `json:"Comment"`
	Bytes             int    `json:"Bytes"`
	Complexity        int    `json:"Complexity"`
	Language          string `json:"Name,omitempty"`
}
