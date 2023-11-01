package service

import (
	"microbenchmarks-data-collector/model"
	"os"

	"github.com/gocarina/gocsv"
)

type CsvGitReposReadWriteStrategy struct{}

func (strategy CsvGitReposReadWriteStrategy) saveToFile(gitRepos []model.GitRepo, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	gocsv.MarshalFile(&gitRepos, file)
}

func (strategy CsvGitReposReadWriteStrategy) readFromFile(path string, target interface{}) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = gocsv.UnmarshalBytes(bytes, target)
	if err != nil {
		panic(err)
	}
}
