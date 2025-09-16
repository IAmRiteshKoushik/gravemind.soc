package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	"github.com/redis/go-redis/v9"
)

type GitHubPRInfo struct {
	RepoOwner         string
	RepoName          string
	PullRequestNumber int
}

type GitHubFile struct {
	Filename string `json:"filename"`
}

func urlParser(prUrl string) (*GitHubPRInfo, error) {
	parsedURL, err := url.Parse(prUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	if parsedURL.Host != "github.com" {
		return nil, fmt.Errorf("not a github.com URL: %w", err)
	}

	pathParts := strings.Split(parsedURL.Path, "/")

	if len(pathParts) < 4 {
		return nil, fmt.Errorf("invalid GitHub PR URL format: not enough path components in %s", prUrl)
	}
	if pathParts[2] != "pull" {
		return nil, fmt.Errorf("invalid GitHub PR URL format: expected 'pull' ")
	}

	owner := pathParts[0]
	repo := pathParts[1]
	prNumStr := pathParts[3]

	prNum, err := strconv.Atoi(prNumStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert PR number '%s' to integer: %w", prNumStr, err)
	}

	return &GitHubPRInfo{
		RepoOwner:         owner,
		RepoName:          repo,
		PullRequestNumber: prNum,
	}, nil
}

func DiscoverFiles(prURL string) ([]string, error) {
	prInfo, err := urlParser(prURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PR URL: %w", err)
	}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/files",
		prInfo.RepoOwner, prInfo.RepoName, prInfo.PullRequestNumber)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// GitHub API recommends sending a User-Agent header
	req.Header.Set("User-Agent", "Go-GitHub-PR-File-Discoverer")
	req.Header.Set("Authorization", "token "+cmd.App.GitHubToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned non-200 status code: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var files []GitHubFile
	err = json.Unmarshal(body, &files)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GitHub API response: %w", err)
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Filename)
	}

	return fileNames, nil
}

func CheckNewFiles() ([]string, bool) {

}

// Function which can check whether the Pirates of Issuebian badge is to be
// given out or not
func CheckIssuebian(username string) (bool, error) {
	sets := []string{
		cmd.CppRank,
		cmd.JavaRank,
		cmd.PyRank,
		cmd.JsRank,
		cmd.GoRank,
		cmd.RustRank,
		cmd.ZigRank,
		cmd.FlutterRank,
		cmd.KotlinRank,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, set := range sets {
		prMergeCount, err := cmd.Valkey.ZScore(ctx, set, username).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return false, err
		}
		if prMergeCount >= 10.0 {
			return true, nil
		}
	}
	return false, nil
}

func FindPrCount(username string, board string) (int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prMergeCount, err := cmd.Valkey.ZScore(ctx, board, username).Result()
	if err == redis.Nil {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return int(prMergeCount), true, nil
}

func IncrPrCount(username string, board string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First increment in language board
	if err := cmd.Valkey.ZIncrBy(ctx, board, 1, username).Err(); err != nil {
		return err
	}
	// Then increment in leaderboard
	if err := cmd.Valkey.ZIncrBy(ctx, cmd.Leaderboard, 0.001, username).Err(); err != nil {
		return err
	}
	return nil
}

func FindRank(username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rank, err := cmd.Valkey.ZRank(ctx, cmd.Leaderboard, username).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int(rank), nil
}
