package model

import "time"

type BenchmarkData struct {
	Timestamp                      time.Time        `json:"timestamp"`
	BenchmarkCount                 int              `json:"benchmarkCount"`
	BenchmarkSuitesCount           int              `json:"benchmarkSuiteCount"`
	SuccessfulBenchmarkCount       int              `json:"successfulBenchmarkCount"`
	SuccessfulBenchmarkSuitesCount int              `json:"successfulBenchmarkSuitesCount"`
	Os                             string           `json:"goos"`
	Architecture                   string           `json:"goarch"`
	Suites                         []BenchmarkSuite `json:"suites"`
}
