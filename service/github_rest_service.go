package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"microbenchmarks-data-collector/config"
	"microbenchmarks-data-collector/dto"
	"microbenchmarks-data-collector/model"
	"microbenchmarks-data-collector/strategies"
	"net/http"
	"net/url"
	"time"
)

func FetchGitHubReposInfo(csvFilePath string, reposCloningCount int) {
	if FileExists(csvFilePath) {
		log.Printf("INFO: File '%s' already exists. Skipping data fetching.", csvFilePath)
		return
	}
	log.Println("INFO: Fetching GitHub repos data...")
	repos := []model.GitRepo{}
	maxResultsCount := 1000
	maxPageSize := 100

	var remainingToClone int
	if reposCloningCount < 0 {
		remainingToClone = maxResultsCount
	} else {
		remainingToClone = reposCloningCount
	}
	baseQuery := "language:go"
	query := fmt.Sprintf("%s+stars:>=25", baseQuery)

MainLoop:
	for remainingToClone > 0 {
		remainingToCloneWithCurrentQuery := calculateRemainingToClone(remainingToClone, maxResultsCount)
		remainingToClone -= remainingToCloneWithCurrentQuery
		pageNumber := 1
		for remainingToCloneWithCurrentQuery > 0 {
			pageSize := calculateRemainingToClone(remainingToCloneWithCurrentQuery, maxPageSize)
			reposToAppend := FetchGitReposByQuery(query, pageSize, pageNumber)
			if len(reposToAppend) == 0 {
				break MainLoop
			}
			repos = append(repos, reposToAppend...)
			remainingToCloneWithCurrentQuery -= pageSize
			pageNumber++
		}
		query = updateQuery(baseQuery, repos)
		if reposCloningCount < 0 {
			remainingToClone = maxResultsCount
		}
	}
	repos = removeDuplicates(repos)
	Save(repos, csvFilePath, strategies.CsvGitReposReadWriteStrategy{})
	log.Printf("INFO: Data fetched and saved in %s", csvFilePath)
}

func removeDuplicates(repos []model.GitRepo) []model.GitRepo {
	isAddedMap := make(map[string]bool)
	var uniqueRepos []model.GitRepo
	for _, repo := range repos {
		if _, alreadyAdded := isAddedMap[repo.FullName]; !alreadyAdded {
			uniqueRepos = append(uniqueRepos, repo)
			isAddedMap[repo.FullName] = true
		}
	}
	return uniqueRepos
}

func calculateRemainingToClone(availableCount, maxAllowedCount int) int {
	if availableCount > maxAllowedCount {
		return maxAllowedCount
	}
	return availableCount
}

func updateQuery(baseQuery string, repos []model.GitRepo) string {
	lastRepo := repos[len(repos)-1]
	starsCount := lastRepo.StarsCount
	return fmt.Sprintf("%s+stars:>=%d", baseQuery, starsCount)
}

func FetchGitReposByQuery(query string, count int, page int) []model.GitRepo {
	const BASE_URL = "https://api.github.com/search/repositories"
	const GITHUB_API_VERSION = "2022-11-28"
	TOKEN := config.GetConfig().GitHubToken

	httpClient := http.DefaultClient
	params := getRequestParameters(count, page)
	apiURL := BASE_URL + "?" + params.Encode() + fmt.Sprintf("&q=%s", query)

	req, _ := http.NewRequest("GET", apiURL, nil)
	if len(TOKEN) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TOKEN))
	}
	req.Header.Add("X-GitHub-Api-Version", GITHUB_API_VERSION)

	log.Printf("INFO: Sending GET request to: %s", apiURL)
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("ERROR: GET request to %s failed. Error: %v", apiURL, err)
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 403 {
			maxRetries := 5
			secondsToWait := 10
			numberOfRetries := 1
			for numberOfRetries <= maxRetries && res.StatusCode != 200 {
				log.Println("INFO: rate limit has probably been exceeded... retrying")
				time.Sleep(time.Duration(secondsToWait) * time.Second)
				res, _ = httpClient.Do(req)
				secondsToWait += 5
				numberOfRetries++
			}
		} else {
			log.Printf("ERROR: Get request to %s returned status: %s", apiURL, res.Status)
		}
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var searchReposResponse dto.SearchReposDto

	if err := json.Unmarshal(resBytes, &searchReposResponse); err != nil {
		log.Fatalf("ERROR: Decoding response failed: %v. Response body: %s", err, res.Body)
	}

	return searchReposResponse.Items
}

func getRequestParameters(count, page int) url.Values {
	params := url.Values{}
	params.Add("sort", "stars")
	params.Add("order", "asc")
	params.Add("per_page", fmt.Sprint(count))
	params.Add("page", fmt.Sprint(page))
	return params
}

func Save(repos []model.GitRepo, path string, saver strategies.GitReposSavingStrategy) {
	saver.SaveToFile(repos, path)
}

func Read(path string, target interface{}, reader strategies.GitReposReadingStrategy) {
	reader.ReadFromFile(path, target)
}
