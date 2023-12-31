package main

import (
	"fmt"
	"log"
	"microbenchmarks-data-collector/config"
	"microbenchmarks-data-collector/model"
	"microbenchmarks-data-collector/service"
	"microbenchmarks-data-collector/strategies"
	"os"
	"strings"
)

func main() {
	config := config.GetConfig()
	csvFilePath := fmt.Sprintf("%s/%s", config.ReposCloningPath, config.ReposInfoFileName)
	reposCloningCount := config.ReposCloningCount

	service.FetchGitHubReposInfo(csvFilePath, reposCloningCount)

	if !config.RunCollectingData {
		log.Println("INFO: Collecting data is disabled.")
		return
	}

	var gitRepos []model.GitRepo
	service.Read(csvFilePath, &gitRepos, strategies.CsvGitReposReadWriteStrategy{})

	var projectsInfo []model.GoProjectInfo
	os.MkdirAll(config.ResultsPath, 0755)
	for _, repo := range gitRepos {
		resultPath := fmt.Sprintf("%s/%s.json", config.ResultsPath, strings.ReplaceAll(repo.FullName, "/", "__"))
		if !service.FileExists(resultPath) {
			projectsInfo = append(projectsInfo, model.GoProjectInfo{
				Name:      repo.FullName,
				Path:      fmt.Sprintf("%s/%s", config.ReposCloningPath, strings.ReplaceAll(repo.FullName, "/", "__")),
				GoVersion: "",
				GitRepo:   repo,
			})
		}
	}
	service.StartAnalysis(projectsInfo, true)
	log.Println("INFO: Collecting data finished")
}
