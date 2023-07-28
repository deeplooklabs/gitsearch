package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	version = "1.0"
	site    = "deeplooklabs.com"
)

var (
	banner = fmt.Sprintf(`
       _ _                       _    
  __ _(_) |_ ___ ___ __ _ _ _ __| |_  
 / _`+"`"+` | |  _(_-</ -_) _`+"`"+` | '_/ _| ' \ 
 \__, |_|\__/__/\___\__,_|_| \__|_||_|    %s
 |___/                                

        %s
`, version, site)

	example = "go run gitsearch.go \"tesla.com boto language:python\""
)

type FileResponse struct {
	HTMLURL string `json:"html_url"`
	Payload struct {
		Repo struct {
			CreatedAt string `json:"created_at"`
		} `json:"repo"`
	} `json:"payload"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(banner)
		fmt.Println(example)
		return
	}

	searchTerm := os.Args[1]
	accessToken := os.Getenv("GITHUB_TOKEN")
	sortBy := "updated"
	headers := map[string]string{"Authorization": "Token " + accessToken}
	dateNow := time.Now()
	yearNow := dateNow.Year()
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s&sort=%s", searchTerm, sortBy)
	response, err := sendRequest(url, headers)
	if err != nil {
		fmt.Println("[WRN] An error occurred while making the request!")
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("[WRN] An error occurred while making the request!")
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("[WRN] An error occurred while reading the response!")
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("[WRN] An error occurred while parsing the response!")
		return
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		fmt.Println("[WRN] An error occurred while parsing the response!")
		return
	}

	fmt.Println(banner)
	fmt.Printf("[INF] Total found: %d\n", len(items))

	for _, item := range items {
		file, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("[WRN] An error occurred while parsing the response!")
			return
		}

		fileName, ok := file["name"].(string)
		if !ok {
			fmt.Println("[WRN] An error occurred while parsing the response!")
			return
		}

		fileURL, ok := file["html_url"].(string)
		if !ok {
			fmt.Println("[WRN] An error occurred while parsing the response!")
			return
		}

		fileResponse, err := sendRequest(fileURL, headers)
		if err != nil {
			fmt.Println("[WRN] An error occurred while making the request!")
			return
		}

		defer fileResponse.Body.Close()

		if fileResponse.StatusCode != 200 {
			fmt.Println("[WRN] An error occurred while making the request!")
			return
		}

		data, err := ioutil.ReadAll(fileResponse.Body)
		if err != nil {
			fmt.Println("[WRN] An error occurred while reading the response!")
			return
		}

		var fileData FileResponse
		if err := json.Unmarshal(data, &fileData); err != nil {
			fmt.Println("[WRN] An error occurred while parsing the response!")
			return
		}

		createdAt := fileData.Payload.Repo.CreatedAt
		if strings.HasPrefix(createdAt, fmt.Sprintf("%d", yearNow)) {
			fmt.Println("[WRN] Recent result!")
		}
		fmt.Printf("[%s] [%s] %s\n", createdAt, fileName, fileURL)
	}
}

func sendRequest(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
