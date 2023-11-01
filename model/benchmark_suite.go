package model

type BenchmarkSuite struct {
	Package                   string                      `json:"pkg"`
	Succeded                  bool                        `json:"succeded"`
	Duration                  string                      `json:"duration"`
	Benchmarks                []Benchmark                 `json:"benchmarks"`
	BenchmarkModificationInfo []BenchmarkModificationInfo `json:"benchmarkModificationInfo"`
}
