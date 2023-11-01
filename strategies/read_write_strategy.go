package strategies

import "microbenchmarks-data-collector/model"

type GitReposSavingStrategy interface {
	SaveToFile(gitRepos []model.GitRepo, path string)
}

type GitReposReadingStrategy interface {
	ReadFromFile(path string, target interface{})
}
