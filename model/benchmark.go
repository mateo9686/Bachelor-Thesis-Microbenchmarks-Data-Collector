package model

type Benchmark struct {
	Name        string  `json:"name"`
	Runs        int     `json:"runs"`
	NsPerOp     float64 `json:"nsPerOp"`
	Succeded    bool    `json:"succeded"`
	FailureDesc string  `json:"failureDesc"`
}
