package commitment

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// A Repository in GitHub
type Repository struct {
	User, RepositoryName string
}

func (r *Repository) getURLForCommitsOnGitHub() string {
	log.Print("https://api.github.com/repos/" + r.User + "/" + r.RepositoryName + "/commits")
	return "https://api.github.com/repos/" + r.User + "/" + r.RepositoryName + "/commits"
}

func (r *Repository) getGitHubResponseForFirstPageOfCommits() *http.Response {
	requestPath := r.getURLForCommitsOnGitHub()
	response := sendRequestForCommitsToGitHub(requestPath)
	linkHeader := getLinkHeaderFromResponse(response)
	urlForLastPage := getURLForLastPageFromLinkHeader(linkHeader)

	if urlForLastPage != "" {
		response.Body.Close()
		response = sendRequestForCommitsToGitHub(urlForLastPage)
	}

	return response
}

// GetCommitReport at index n from Repository
func (r *Repository) GetCommitReport(n int) CommitReport {
	response := r.getGitHubResponseForFirstPageOfCommits()
	commitItems := getCommitItemsFromResponse(response)

	indexOfFirstCommitItem := len(commitItems) - 1
	indexOfSelectedCommitItem := indexOfFirstCommitItem - n
	selectedCommitItem := commitItems[indexOfSelectedCommitItem]

	return CommitReport{
		Message: selectedCommitItem.Commit.Message,
		URL:     selectedCommitItem.Html_url,
		Date:    selectedCommitItem.Commit.Committer.Date,
	}
}

func getCommitItemsFromResponse(response *http.Response) []commitItem {
	body, errReadingResponseBody := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if errReadingResponseBody != nil {
		log.Fatal(errReadingResponseBody)
	}

	var commitItems []commitItem
	unmarshallingErr := json.Unmarshal(body, &commitItems)
	if unmarshallingErr != nil {
		log.Fatal("Unable to retrieve repository information from GitHub.")
	}

	return commitItems
}

func sendRequestForCommitsToGitHub(fullRequestPath string) *http.Response {
	response, errGettingHTTPResource := http.Get(fullRequestPath)

	if errGettingHTTPResource != nil {
		log.Fatal(errGettingHTTPResource)
	}

	return response
}

func getLinkHeaderFromResponse(response *http.Response) string {
	return response.Header.Get("Link")
}

func getSubstringBetweenIndexes(s string, startIndex int, endIndex int) string {
	characters := strings.Split(s, "")
	substringCharacters := characters[startIndex:endIndex]
	substring := strings.Join(substringCharacters, "")
	return substring
}

func getURLFromLinkHeaderItem(linkHeaderItem string) string {
	indexOfOpeningAngleBracket := strings.Index(linkHeaderItem, "<") + 1
	indexOfClosingAngleBracket := strings.Index(linkHeaderItem, ">")
	return getSubstringBetweenIndexes(linkHeaderItem, indexOfOpeningAngleBracket, indexOfClosingAngleBracket)
}

func getURLForLastPageFromLinkHeader(linkHeader string) string {
	links := strings.Split(linkHeader, ",")

	for _, link := range links {
		if strings.Contains(link, "rel=\"last\"") {
			return getURLFromLinkHeaderItem(link)
		}
	}

	return ""
}
