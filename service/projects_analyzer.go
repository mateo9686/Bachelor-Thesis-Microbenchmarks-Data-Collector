package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"repos-fetcher/config"
	"repos-fetcher/model"
	"strings"
	"sync"
)

func StartAnalysis(projectsInfo []model.GoProjectInfo, shouldClone bool) {
	config := config.GetConfig()
	functions := prepareJobsToRun(projectsInfo, config.ResultsPath, config.ClearCacheJobInterval)
	runAnalysisInParallel(functions, config.MaxWorkersCount)
	// clear cache at the end
	clearCacheJob := prepareClearCacheJob()
	clearCacheJob()
}

func prepareJobsToRun(projectsInfo []model.GoProjectInfo, resultsPath string, clearCacheJobsInterval int) []func() {
	var functions []func()
	for index, pInfo := range projectsInfo {
		functions = append(functions, prepareCollectProjectDataJob(pInfo, resultsPath))
		if index > 0 && index%clearCacheJobsInterval == 0 {
			functions = append(functions, prepareClearCacheJob())
		}
	}
	return functions
}

func prepareClearCacheJob() func() {
	return func() {
		log.Println("INFO: Cleaning cache...")
		cleanCacheCmd := exec.Command("go", "clean", "-cache")
		cleanModCacheCmd := exec.Command("go", "clean", "-modcache")
		cleanCacheCmd.Output()
		cleanModCacheCmd.Output()
		log.Println("INFO: Cache cleaned...")
	}
}

func prepareCollectProjectDataJob(pInfo model.GoProjectInfo, resultsPath string) func() {
	return func() {
		// put here cloning, running benchmarks and deleting
		cloningErr := cloneRepo(pInfo.GitRepo.SshUrl, pInfo.Name, pInfo.Path)
		if cloningErr != nil {
			pInfo.GitRepo.CloningError = cloningErr.Error()
		} else {
			collectProjectInformation(&pInfo)
			removeFolder(pInfo.Path)
		}
		saveProjectInfo(pInfo, resultsPath)
	}
}

func saveProjectInfo(projectInfo model.GoProjectInfo, path string) {
	file, err := json.MarshalIndent(projectInfo, "", "  ")
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", path, strings.ReplaceAll(projectInfo.GitRepo.FullName, "/", "__")), file, 0666)
}

func worker(id int, jobs <-chan func(), wg *sync.WaitGroup) {
	defer wg.Done()

	for fn := range jobs {
		fn()
		log.Printf("INFO: Worker %d completed a job", id)
	}
}

func runAnalysisInParallel(functions []func(), maxWorkers int) {
	var wg sync.WaitGroup
	jobs := make(chan func())

	// Launch workers
	for i := 1; i <= maxWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Send jobs to workers
	for _, fn := range functions {
		jobs <- fn
	}

	// Close the jobs channel to signal workers that no more jobs are coming
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()
}
