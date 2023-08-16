package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
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

	example = "gitsearch \"tesla.com boto language:python\""
)

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
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s&sort=%s", searchTerm, sortBy)
	response, err := sendRequest(url, headers)
	if err != nil {
		fmt.Printf("[WRN] An error occurred while making the request! to %s", url)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("[WRN] An error occurred while making the request! {status_code: %d}", response.StatusCode)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("[WRN] An error occurred while reading the response!")
		return
	}

	json := string(data[:])

	fmt.Println(banner)

	items := gjson.Get(json, "items.#.html_url").Array()
	fmt.Printf("[INF] Total found: %d\n", len(items))
	for _, value := range items {
		url := value.String()
		lastSlashIndex := strings.LastIndex(url, "/")
		nameOfFile := url[lastSlashIndex+1:]
		fmt.Printf("[%s] %s\n", nameOfFile, url)
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
