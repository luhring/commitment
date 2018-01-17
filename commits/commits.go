package commits

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

// CommitReport for display to user
type CommitReport struct {
	Message, URL string
}

// GetCommitReport at index n from Repository
func (r *Repository) GetCommitReport(n int) CommitReport {
	const baseURL = "https://api.github.com"

	initialResponse, errGettingHTTPResourceInitially := http.Get(baseURL + "/repos/" + r.User + "/" + r.RepositoryName + "/commits")

	if errGettingHTTPResourceInitially != nil {
		log.Fatal(errGettingHTTPResourceInitially)
	}

	defer initialResponse.Body.Close()

	linkHeader := initialResponse.Header.Get("Link")
	urlForLastPage := getURLForLastPageFromLinkHeader(linkHeader)

	response, errGettingHTTPResource := http.Get(urlForLastPage)

	if errGettingHTTPResource != nil {
		log.Fatal(errGettingHTTPResource)
	}

	defer response.Body.Close()

	body, errReadingResponseBody := ioutil.ReadAll(response.Body)

	if errReadingResponseBody != nil {
		log.Fatal(errReadingResponseBody)
	}

	var commitItems []CommitItem
	unmarshallingErr := json.Unmarshal(body, &commitItems)
	if unmarshallingErr != nil {
		log.Fatal(unmarshallingErr)
	}

	indexOfFirstCommitItem := len(commitItems) - 1

	selectedCommitItem := commitItems[indexOfFirstCommitItem-n]

	return CommitReport{
		Message: selectedCommitItem.Commit.Message,
		URL:     selectedCommitItem.Html_url,
	}
}

type CommitItem struct {
	Sha          string
	Commit       Commit
	Url          string
	Html_url     string
	Comments_url string
}

type Commit struct {
	Committer Committer
	Message   string
}

type Committer struct {
	Name, Email, Date string
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

type PageLink struct {
	URL, rel string
}
