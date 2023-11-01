package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var lock = &sync.Mutex{}

type config struct {
	RunCollectingData     bool
	RunBenchmarks         bool
	ReposCloningPath      string
	ResultsPath           string
	ReposCloningCount     int
	ReposInfoFileName     string
	MaxWorkersCount       int
	GitHubToken           string
	ClearCacheJobInterval int
}

var configInstance *config

func GetConfig() *config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil {
			configInstance = loadConfig()
		}
	}

	return configInstance
}

func loadConfig() *config {
	configInstance = &config{}
	godotenv.Load()

	configInstance.RunCollectingData, _ = strconv.ParseBool(getEnv("RUN_COLLECTING_DATA", "true"))

	configInstance.RunBenchmarks, _ = strconv.ParseBool(getEnv("RUN_BENCHMARKS", "true"))

	configInstance.ReposCloningPath = getEnv("REPOS_CLONING_PATH", "./repositories")

	configInstance.ResultsPath = getEnv("RESULTS_PATH", "./results")

	clearCacheJobInterval, _ := strconv.ParseInt(getEnv("CLEAR_CACHE_JOB_INTERVAL", "1000"), 10, 32)
	configInstance.ClearCacheJobInterval = int(clearCacheJobInterval)

	configInstance.ReposInfoFileName = getEnv("REPOS_INFO_FILE_NAME", "repos.csv")

	reposCloningCount, _ := strconv.ParseInt(getEnv("REPOS_CLONING_COUNT", "-1"), 10, 32)
	configInstance.ReposCloningCount = int(reposCloningCount)

	maxWorkersCount, _ := strconv.ParseInt(getEnv("MAX_WORKERS_COUNT", "10"), 10, 32)
	configInstance.MaxWorkersCount = int(maxWorkersCount)

	configInstance.GitHubToken = getEnv("GITHUB_TOKEN", "")

	return configInstance
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
