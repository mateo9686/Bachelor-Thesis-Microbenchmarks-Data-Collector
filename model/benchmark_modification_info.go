package model

type BenchmarkModificationInfo struct {
	Name              string `json:"name"`
	ModificationCount int    `json:"modificationCount"`
	CreatedOn         string `json:"createdOn"`
	LastUpdateOn      string `json:"lastUpdateOn"`
}
