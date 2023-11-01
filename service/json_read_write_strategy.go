package service

import (
	"encoding/json"
	"io/ioutil"
	"microbenchmarks-data-collector/model"
)

type JsonGitReposReadWriteStrategy struct{}

func (strategy JsonGitReposReadWriteStrategy) saveToFile(gitRepos []model.GitRepo, path string) {
	file, err := json.MarshalIndent(gitRepos, "", " ")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile("test.json", file, 0666)
}

func (strategy JsonGitReposReadWriteStrategy) readFromFile(path string, target interface{}) {}
