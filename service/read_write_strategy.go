package service

import "repos-fetcher/model"

type GitReposSavingStrategy interface {
	saveToFile(gitRepos []model.GitRepo, path string)
}

type GitReposReadingStrategy interface {
	readFromFile(path string, target interface{})
}
