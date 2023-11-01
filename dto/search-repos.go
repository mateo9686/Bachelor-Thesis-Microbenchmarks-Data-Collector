package dto

import "repos-fetcher/model"

type SearchReposDto struct {
	TotalCount int `json:"total_count"`
	Items      []model.GitRepo
}
