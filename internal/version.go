package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GoVersion struct {
	Version   string   `json:"version"`
	Downloads []string `json:"downloads"`
}

func FetchGoVersions() map[string]GoVersion {
	url := "https://mirrors.aliyun.com/golang/"
	versions, err := fetchGoVersions(url)
	if err != nil {
		log.Printf("Failed to fetch Go versions: %v", err)
	}
	return versions
}

func fetchGoVersions(url string) (map[string]GoVersion, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the webpage: %v", err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to load the webpage: %v", err)
	}

	versions := make(map[string]GoVersion)
	re := regexp.MustCompile(`go(\d+\.\d+(\.\d+)?)`)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || !strings.HasPrefix(href, "go") {
			return
		}

		// 提取 Go 版本号
		matches := re.FindStringSubmatch(href)
		if len(matches) < 2 {
			return
		}
		version := matches[1]
		downloadURL := url + href

		// 如果版本不存在，则初始化
		if _, found := versions[version]; !found {
			versions[version] = GoVersion{
				Version:   version,
				Downloads: []string{},
			}
		}

		// 追加下载链接
		v := versions[version]
		v.Downloads = append(v.Downloads, downloadURL)
		versions[version] = v
	})

	return versions, nil
}

var GVMHome = os.Getenv("HOME") + "/gvm"

func FetchLocalVersions() []string {
	// 列出本地已安装的 Go 版本， 目录为~/gvm/versions
	// 获取 ~/gvm/versions 下的目录名称，并返回
	dis, err := os.ReadDir(GVMHome + "/versions")
	if err != nil {
		log.Fatal(err)
	}
	var versions []string
	for _, di := range dis {
		if di.IsDir() {
			dirName := di.Name()
			if strings.HasPrefix(dirName, "go") {
				versions = append(versions, strings.TrimPrefix(dirName, "go"))
			}
		}
	}
	return versions
}
