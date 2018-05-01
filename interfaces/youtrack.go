package interfaces

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"regexp"

	"github.com/vizualni/meyougotrack/domain"
)

type YouTrackClient struct {
	baseUrl    string
	token      string
	httpClient HttpClient
}

type workLogXml struct {
	XMLName     xml.Name `xml:"workItem"`
	Date        int64    `xml:"date"`
	Duration    int      `xml:"duration"`
	Description string   `xml:"description"`
	Worktype    worktype `xml:"worktype"`
}
type worktype struct {
	Name string `xml:"name"`
}

func (youtrackClient *YouTrackClient) SaveWorkLog(workLog domain.IssueWorkLog) error {
	request, _ := youtrackClient.buildRequest(fmt.Sprintf("%s/rest/issue/%s/timetracking/workitem", youtrackClient.baseUrl, workLog.IssueId))
	request.Method = "POST"
	request.Header.Add("Content-Type", "application/xml")

	var workLogXml workLogXml

	workLogXml.Date = workLog.Date.Unix() * 1000 // because milliseconds
	workLogXml.Duration = workLog.Duration
	workLogXml.Description = workLog.Description
	workLogXml.Worktype = worktype{workLog.Type}

	byteBody, e := xml.Marshal(workLogXml)

	request.Body = ioutil.NopCloser(bytes.NewReader(byteBody))
	request.ContentLength = int64(len(byteBody))
	response, e := youtrackClient.httpClient.Do(&request)

	byteBody, e = ioutil.ReadAll(response.Body)

	if e != nil {
		return e
	}

	if response.StatusCode != http.StatusCreated {
		return errors.New("no new work item created")
	}

	return nil
}

func (youtrackClient *YouTrackClient) buildRequest(requestUrl string) (http.Request, error) {
	r := http.Request{}

	r.Header = http.Header{}

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", youtrackClient.token))

	u, err := url.Parse(requestUrl)

	if err != nil {
		return http.Request{}, err
	}

	r.URL = u

	return r, nil
}

func (youtrackClient *YouTrackClient) FindIssueByIssueId(issueId string) (domain.YouTrackIssue, error) {

	request, _ := youtrackClient.buildRequest(fmt.Sprintf("%s/rest/issue/%s", youtrackClient.baseUrl, issueId))

	response, e := youtrackClient.httpClient.Do(&request)

	fmt.Println(e)

	b, _ := ioutil.ReadAll(response.Body)

	var issue domain.YouTrackIssue

	e = xml.Unmarshal(b, &issue)

	return issue, nil
}

func NewYouTrackClient(baseUrl, apikey string) *YouTrackClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &YouTrackClient{
		httpClient: &OfficialHttpClientAdapter{&http.Client{Transport: tr}},
		token:      apikey,
		baseUrl:    baseUrl,
	}
}

type SimpleRegexIssueIdExtractor struct{}

func (SimpleRegexIssueIdExtractor) Extract(input string) (string, error) {
	regex := "https?://(.+?)/issue/(\\w+-\\d+\\b+)"

	r, err := regexp.Compile(regex)

	if err != nil {
		panic(err)
	}

	solutions := r.FindStringSubmatch(input)

	if len(solutions) != 3 {
		return "", errors.New("issue id not found")
	}

	return solutions[2], nil
}
