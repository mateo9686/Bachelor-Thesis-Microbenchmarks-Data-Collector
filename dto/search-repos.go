package dto

import "microbenchmarks-data-collector/model"

type SearchReposDto struct {
	TotalCount int `json:"total_count"`
	Items      []model.GitRepo
}
