package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/t0mk/rocketreport/config"

)

const (
	treesMissingMsg = `
	You need to download the reward trees to use the *reward* plugins
	You can download them by running
		$ NETWORK=mainnet rocketreport download-reward-trees
	Specify the network you are using
	`
)


func checkIfAllRewardsTreesDownloaded() error {
	dataDir, err := config.RewardTreesPath()
	if err != nil {
		return fmt.Errorf("can't get reward trees path: %v", err)
	}
	_, err = os.Stat(dataDir)
	if err != nil {
		return fmt.Errorf("can't find %s: %v\n%s", dataDir, err, treesMissingMsg)
	}
	urls, err := RewardTreesRawUrls()
	if err != nil {
		return fmt.Errorf("can't get reward trees urls: %v", err)
	}
	if len(urls) == 0 {
		return fmt.Errorf("no reward trees found\n%s", treesMissingMsg)
	}
	return nil
}

func RewardTreesDownloadedCheck() (bool, error) {
	readyToDownload := false
	dataDir, err := config.RewardTreesPath()
	if err != nil {
		return readyToDownload, fmt.Errorf("can't get reward trees path: %v", err)
	}
	_, err = os.Stat(dataDir)
	if err != nil {
		return readyToDownload, fmt.Errorf("can't find %s: %v\n%s", dataDir, err, treesMissingMsg)
	}
	urls, err := RewardTreesRawUrls()
	if err != nil {
		return readyToDownload, fmt.Errorf("can't get reward trees urls: %v", err)
	}
	if len(urls) == 0 {
		return readyToDownload, fmt.Errorf("no reward trees found\n%s", treesMissingMsg)
	}
	existingFilesList, err := filepath.Glob(filepath.Join(dataDir, "*"))
	if err != nil {
		return readyToDownload, fmt.Errorf("can't get list of files in %s: %v", dataDir, err)
	}
	existingFilesMap := make(map[string]struct{})
	for _, f := range existingFilesList {
		existingFilesMap[f] = struct{}{}
	}
	readyToDownload = true
	if len(existingFilesList) == 0 {
		return readyToDownload, fmt.Errorf("no reward trees downloaded\n%s", treesMissingMsg)
	}
	for _, url := range urls {
		if _, ok := existingFilesMap[filepath.Join(dataDir, url[strings.LastIndex(url, "/")+1:])]; !ok {
			return readyToDownload, fmt.Errorf("not all reward trees downloaded\n%s", treesMissingMsg)
		}
	}
	return false, nil
}

type GitHubFile struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	RawURL string `json:"download_url"`
}

func RewardTreesRawUrls() ([]string, error) {
	network := string(config.Network())
	endpoint := "https://api.github.com/repos/rocket-pool/rewards-trees/contents/" + network
	raw, err := GetHTTPResponseBodyFromUrl(endpoint)
	if err != nil {
		return nil, err
	}
	var files []GitHubFile
	err = json.Unmarshal(raw, &files)
	if err != nil {
		return nil, err
	}
	var urls []string
	for _, file := range files {
		if strings.Contains(file.Name, "rp-rewards-"+network) {
			urls = append(urls, file.RawURL)
		}
	}
	return urls, nil
}

func downloadFile(url string, filename string) error {
	raw, err := GetHTTPResponseBodyFromUrl(url)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, strings.NewReader(string(raw)))
	if err != nil {
		return err
	}
	return nil
}

func donwloadFiles(urls []string, destinationDir string) error {
	sem := make(chan struct{}, 4)
	wg := &sync.WaitGroup{}

	for _, url := range urls {
		sem <- struct{}{}
		wg.Add(1)
		go func(url string) {
			defer func() {
				<-sem
				wg.Done()
			}()
			destination := filepath.Join(destinationDir, url[strings.LastIndex(url, "/")+1:])
			fmt.Printf("Downloading %s to\n  %s\n", url, destination)
			err := downloadFile(url, destination)
			if err != nil {
				fmt.Println(err)
			}
		}(url)
	}

	wg.Wait()
	return nil
}

func DownloadRewardTrees() error {
	destinationDir, err := config.RewardTreesPath()
	if err != nil {
		return err
	}
	err = os.MkdirAll(destinationDir, 0755)
	if err != nil {
		return err
	}

	urls, err := RewardTreesRawUrls()
	if err != nil {
		return err
	}

	err = donwloadFiles(urls, destinationDir)
	if err != nil {
		return err
	}

	return nil
}

/*
	dataDir := filepath.Join(os.UserHomeDir(), ".rocketreport", "reward-trees")
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
*/
