package model

type GoProjectInfo struct {
	Name           string //same gitRepo full name. It is unique and can be used also as an id
	Path           string
	GoVersion      string
	BenchmarkData  *BenchmarkData
	GitRepo        GitRepo
	CodeStatistics *CodeStats
}
