package service

import (
	"bytes"
	"encoding/json"
	"log"
	"microbenchmarks-data-collector/config"
	"microbenchmarks-data-collector/model"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func collectProjectInformation(projectInfo *model.GoProjectInfo) {
	log.Printf("INFO: Started collecting data about %s", projectInfo.Name)
	(*projectInfo).GoVersion, _ = getProjectGoVersion(projectInfo.Path)
	(*projectInfo).BenchmarkData = collectBenchmarkData(projectInfo.Path, projectInfo.Name)
	(*projectInfo).CodeStatistics = collectCodeStatistics(projectInfo.Path, projectInfo.Name)
}

func getProjectGoVersion(path string) (string, error) {
	cmd := exec.Command("grep", "-oP", `^go \K\d+\.\d+`, path+"/go.mod")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("WARN: Failed to check GO version under %s: %v", path, err)
		return "n/a", err
	}
	version := string(output)
	return strings.Trim(version, "\n"), nil
}

func collectCodeStatistics(path string, projectName string) *model.CodeStats {
	log.Printf("INFO: Running 'scc' on %s to collect code statistics...", projectName)
	sccCmd := exec.Command("scc", "--format", "json", path)
	output, err := sccCmd.Output()
	if err != nil {
		log.Printf("ERROR: Failed to run 'scc' on %s: %v", projectName, err)
	}
	var codeStatistics []model.CodeStats
	err = json.Unmarshal(output, &codeStatistics)
	if err != nil {
		log.Printf("ERROR: Failed to parse 'scc' JSON result for %s: %v", projectName, err)
	}
	if len(codeStatistics) == 0 {
		log.Printf("Warning: The 'scc' command did not return any results for project '%s'.", projectName)
		return nil
	}
	for _, stats := range codeStatistics {
		if stats.Language == "Go" {
			log.Printf("INFO: Finished collecting code statistics for %s", projectName)
			stats.Language = ""
			return &stats
		}
	}
	log.Printf("Warning: The 'scc' command did not collect any Go statistics for project '%s'.", projectName)
	return nil
}

func collectBenchmarkData(path string, projectName string) *model.BenchmarkData {
	benchmarkData := model.BenchmarkData{}
	if !config.GetConfig().RunBenchmarks {
		log.Println("INFO: Benchmarks are disabled")
		return nil
	}
	log.Printf("INFO: Running benchmarks for %s...", projectName)
	benchmarkData.Timestamp = time.Now().UTC().Truncate(time.Second)
	benchmarkCmdOutput := runBenchmarks(path, projectName)
	if benchmarkCmdOutput == nil {
		return nil
	}
	log.Printf("INFO: Running benchmarks for %s finished. Collecting benchmark data...", projectName)
	benchmarkData.Architecture, benchmarkData.Os = getBenchmarkRunEnv(benchmarkCmdOutput)
	benchmarkData.Suites = getBenchmarksResults(benchmarkCmdOutput, path)
	benchmarkData.BenchmarkSuitesCount, benchmarkData.SuccessfulBenchmarkSuitesCount, benchmarkData.BenchmarkCount, benchmarkData.SuccessfulBenchmarkCount = getBenchmarkCounts(benchmarkData.Suites)
	return &benchmarkData
}

func runBenchmarks(path, projectName string) []byte {
	benchmarkCmd := exec.Command("find", path, "-name", "go.mod", "-execdir", "go", "test", "-run=^$", "-bench=.", "-benchtime=1s", "-timeout=10m", "./...", ";")
	benchmarkCmdOutput, err := benchmarkCmd.Output()
	if err != nil {
		log.Printf("ERROR: Failed to run benchmarks for project '%s': %v", projectName, err)

		return nil
	}
	return benchmarkCmdOutput
}

func getBenchmarkCounts(benchmarkSuites []model.BenchmarkSuite) (int, int, int, int) {
	var benchmarkSuitesCount, successfulBenchmarkSuitesCount, benchmarkCount, successfulBenchmarkCount int
	for _, suite := range benchmarkSuites {
		benchmarkSuitesCount++
		benchmarkCount += len(suite.Benchmarks)
		successfulBenchmarkCountInSuite := 0
		for _, benchmark := range suite.Benchmarks {
			if benchmark.Succeded {
				successfulBenchmarkCountInSuite++
			}
		}
		successfulBenchmarkCount += successfulBenchmarkCountInSuite
		if successfulBenchmarkCountInSuite > 0 {
			successfulBenchmarkSuitesCount++
		}
	}
	return benchmarkSuitesCount, successfulBenchmarkSuitesCount, benchmarkCount, successfulBenchmarkCount
}

func getBenchmarkRunEnv(benchmarkCmdOutput []byte) (string, string) {
	arch, os := "n/a", "n/a"
	grepCmd := exec.Command("grep", "-e", "goos", "-e", "goarch")
	grepCmd.Stdin = bytes.NewReader(benchmarkCmdOutput)
	grepCmdOutput, err := grepCmd.Output()
	if err != nil {
		return arch, os
	}
	outputLines := strings.Split(string(grepCmdOutput), "\n")
	if len(outputLines) < 2 {
		return arch, os
	}
	outputLines = outputLines[:2]
	for _, value := range outputLines {
		if strings.HasPrefix(value, "goos: ") {
			os = strings.TrimPrefix(value, "goos: ")
		}
		if strings.HasPrefix(value, "goarch: ") {
			arch = strings.TrimPrefix(value, "goarch: ")
		}
	}
	return arch, os
}

func getBenchmarksResults(benchmarkCmdOutput []byte, projectPath string) []model.BenchmarkSuite {
	outputLines := strings.Split(string(benchmarkCmdOutput), "\n")
	benchmarkRegex := `(\S+)\s+(\d+)\s+([\d.]+)\s+\S+`
	failedBenchmarkRegex := `(\S+)\s+(\S+)`
	suites := []model.BenchmarkSuite{}
	suite := model.BenchmarkSuite{}

	for index, outputLine := range outputLines {
		prefix := "pkg: "
		if strings.HasPrefix(outputLine, prefix) {
			suite.Package = strings.TrimPrefix(outputLine, prefix)
		}
		if strings.HasPrefix(outputLine, "Benchmark") {
			benchmark := model.Benchmark{}
			re := regexp.MustCompile(benchmarkRegex)
			matches := re.FindStringSubmatch(outputLine)
			if len(matches) > 0 {
				benchmark.Name = matches[1]
				benchmark.Runs, _ = strconv.Atoi(matches[2])
				benchmark.Succeded = true
				benchmark.NsPerOp, _ = strconv.ParseFloat(matches[3], 64)
				suite.Benchmarks = append(suite.Benchmarks, benchmark)
			} else {
				re = regexp.MustCompile(failedBenchmarkRegex)
				matches = re.FindStringSubmatch(outputLine)
				if len(matches) < 1 {
					break
				}
				printFailedBenchmarkLog(matches[1])
				benchmark.Name = matches[1]
				benchmark.FailureDesc = matches[2]
				benchmark.Succeded = false
				suite.Benchmarks = append(suite.Benchmarks, benchmark)
			}
		}
		if strings.HasPrefix(outputLine, "--- FAIL") && strings.HasPrefix(outputLines[index+1], "    ") { // second condition checks if the next line reports an error
			benchmarkFailureRegex := regexp.MustCompile(`--- FAIL: (Benchmark[^\s]+)( \(\S+\))?`)
			matches := benchmarkFailureRegex.FindStringSubmatch(outputLine)
			if len(matches) < 1 {
				break
			}
			printFailedBenchmarkLog(matches[1])
			benchmark := model.Benchmark{}
			benchmark.Name = matches[1]
			benchmark.Succeded = false
			benchmark.FailureDesc = strings.TrimSpace(outputLines[index+1]) // next line after failure log contains error message
			suite.Benchmarks = append(suite.Benchmarks, benchmark)
		}
		if (strings.HasPrefix(outputLine, "FAIL") || strings.HasPrefix(outputLine, "ok")) && suite.Benchmarks != nil { //
			pattern := `^(FAIL|ok)\s+([^ ]+)\s+([\d.]+[a-z]+)$`
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(outputLine)

			if len(matches) > 0 {
				// Extract the relevant information from submatches
				suite.Succeded = matches[1] == "ok"
				suite.Package = matches[2]
				suite.Duration = matches[3]
				suites = append(suites, suite)
				suite = model.BenchmarkSuite{}
			}
		}
	}
	for index, suite := range suites {
		suites[index].BenchmarkModificationInfo = getBenchmarkModificationInfo(suite, projectPath)
	}
	return suites
}

func getBenchmarkModificationInfo(suite model.BenchmarkSuite, projectPath string) []model.BenchmarkModificationInfo {
	benchmarkFunctionNames := []string{}
	for _, benchmark := range suite.Benchmarks {
		name := removeTrailingNumber(strings.Split(benchmark.Name, "/")[0])
		if !contains(benchmarkFunctionNames, benchmark.Name) {
			benchmarkFunctionNames = append(benchmarkFunctionNames, name)
		}
	}
	benchmarkFunctionInfos := []model.BenchmarkModificationInfo{}
	for _, benchmarkFuncName := range benchmarkFunctionNames {
		cmd := exec.Command("grep", "-rn", "func "+benchmarkFuncName, projectPath)
		cmdOutput, err := cmd.Output()
		if err != nil {
			log.Printf("WARN: Failed to search for file containing benchmark: %s", benchmarkFuncName)
			return nil
		}
		filePath := "." + strings.TrimPrefix(strings.Split(string(cmdOutput), ":")[0], projectPath)
		gitLogCmd := exec.Command("git", "log", "-L", ":"+benchmarkFuncName+":"+filePath)
		gitLogCmd.Dir = projectPath
		cmdOutput, err = gitLogCmd.Output()
		if err != nil {
			log.Printf("WARN: Failed to execute 'git log' for: %s", benchmarkFuncName)
			continue
		}
		datePrefix := "Date:   "
		grepCmd := exec.Command("grep", datePrefix)
		grepCmd.Stdin = bytes.NewReader(cmdOutput)
		cmdOutput, err = grepCmd.Output()
		if err != nil {
			log.Printf("WARN: Failed to execute 'grep' for 'git log' for function: %s. Error: %v", benchmarkFuncName, err)
			continue
		}
		funcModificationDates := strings.Split(string(cmdOutput), "\n")
		funcModificationDates = funcModificationDates[:len(funcModificationDates)-1]

		lastModificationDate := parseDate(strings.TrimPrefix(funcModificationDates[0], datePrefix))
		creationDate := parseDate(strings.TrimPrefix(funcModificationDates[len(funcModificationDates)-1], datePrefix))
		benchmarkFunctionInfos = append(benchmarkFunctionInfos, model.BenchmarkModificationInfo{
			Name:              benchmarkFuncName,
			CreatedOn:         creationDate,
			LastUpdateOn:      lastModificationDate,
			ModificationCount: len(funcModificationDates),
		})
	}
	return benchmarkFunctionInfos
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func parseDate(gitLogDateString string) string {
	layout := "Mon Jan 2 15:04:05 2006 -0700"

	parsedDate, err := time.Parse(layout, gitLogDateString)
	if err != nil {
		log.Printf("ERROR: Failed to parse date: %v", err)
		return "n/a"
	}

	desiredLayout := "2006-01-02"
	return parsedDate.Format(desiredLayout)
}

func removeTrailingNumber(input string) string {
	re := regexp.MustCompile(`-\d+$`)
	return re.ReplaceAllString(input, "")
}

func printFailedBenchmarkLog(name string) {
	log.Printf("WARN: Failed benchmark: %s", name)
}
