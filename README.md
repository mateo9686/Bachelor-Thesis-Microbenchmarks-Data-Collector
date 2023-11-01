# Repos Fetcher

## Dependencies

- [scc](https://github.com/boyter/scc) (Source Code Counter)
- grep
- git
- find

## Installation

To install **scc**, you can use the following command:

```sh
go install github.com/boyter/scc/v3@latest
```

## Configuration

Before running the Repos Fetcher, make sure to set the following environment variables in your configuration:

- **RUN_COLLECTING_DATA**: Set to `true` to enable data collection. (default `true`)
- **RUN_BENCHMARKS**: Set to `true` to enable benchmarking. (default `true`)
- **REPOS_CLONING_PATH**: The path where repositories will be cloned. (default `./repositories`)
- **RESULTS_PATH**: The directory for storing results. (default `./results`)
- **REPOS_CLONING_COUNT**: The number of repositories to clone (values less than 0 will fetch all repositories). (default `-1`)
- **REPOS_INFO_FILE_NAME**: The name of the CSV file for storing repository information. (default `repos.csv`)
- **MAX_WORKERS_COUNT**: The maximum number of concurrent workers. (default `10`)
- **CLEAR_CACHE_JOB_INTERVAL**: The interval (in milliseconds) for clearing the cache. (default `1000`)
- **GITHUB_TOKEN**: Your GitHub access token (optional, but recommended to avoid rate limits).

Please ensure that you've set these environment variables correctly to customize the behavior of the Repos Fetcher.

Make sure to review and adapt the configuration according to your requirements before running the tool.

## Execution

The `start.sh` script should be used to start the data collector. This script is responsible for the following operations:

1. deleting temporary Go builds
2. checking if the results folder exists and, if necessary, creating it
3. adding the logs related to the previous launch of the collector to the main log file.
4. run the program with a soft memory limit

If necessary, the soft memory limit can be adjusted to the environment in which the program is run (value of the `GOMEMLIMIT` variable).
More information - https://pkg.go.dev/runtime

## Analysis

the analysis directory contains a python jupyter notebook with the functions used for data analysis.
