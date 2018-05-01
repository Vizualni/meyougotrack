package usecases_test

import (
	"testing"

	"time"

	"fmt"

	"errors"

	"github.com/vizualni/meyougotrack/domain"
	"github.com/vizualni/meyougotrack/interfaces"
	"github.com/vizualni/meyougotrack/usecases"
)

type youtrackRepositoryMock struct {
	findIssue   func(issueId string) (domain.YouTrackIssue, error)
	saveWorkLog func(workLog domain.IssueWorkLog) error
}

func (y *youtrackRepositoryMock) FindIssueByIssueId(issueId string) (domain.YouTrackIssue, error) {
	return y.findIssue(issueId)
}

func (y *youtrackRepositoryMock) SaveWorkLog(workLog domain.IssueWorkLog) error {
	return y.saveWorkLog(workLog)
}

type trelloRepositoryMock struct {
	getAllCards func(boardId string, doingListName string) []domain.TrelloCardDoingDuration
}

func (t *trelloRepositoryMock) GetAllCards(boardId string, doingListName string) []domain.TrelloCardDoingDuration {
	return t.getAllCards(boardId, doingListName)
}

func TestGetWithNoCardsReturned(t *testing.T) {
	interactor := usecases.TimeLoggerInteractor{
		IssueIdExtractor:   interfaces.SimpleRegexIssueIdExtractor{},
		YouTrackRepository: &youtrackRepositoryMock{},
		TrelloRepository: &trelloRepositoryMock{
			getAllCards: func(boardId string, doingListName string) []domain.TrelloCardDoingDuration {
				return []domain.TrelloCardDoingDuration{}
			},
		},
	}

	linkedCards := interactor.GetTrelloYouTrackCards("something", "another")

	if len(linkedCards) > 0 {
		t.Fatal("Should have returned empty set")
	}

}

func TestGetWithSingleCardReturned(t *testing.T) {
	interactor := usecases.TimeLoggerInteractor{
		IssueIdExtractor: interfaces.SimpleRegexIssueIdExtractor{},
		YouTrackRepository: &youtrackRepositoryMock{
			findIssue: func(issueId string) (domain.YouTrackIssue, error) {
				t.Fatal("if I get called then something is wrong here")
				return domain.YouTrackIssue{}, nil
			},
		},
		TrelloRepository: &trelloRepositoryMock{
			getAllCards: func(boardId string, doingListName string) []domain.TrelloCardDoingDuration {
				return []domain.TrelloCardDoingDuration{
					{
						Title:           "title",
						Description:     "description",
						Duration:        123,
						Date:            time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
						YoutrackSummary: "something",
					},
				}
			},
		},
	}

	linkedCards := interactor.GetTrelloYouTrackCards("something", "another")

	if len(linkedCards) != 1 {
		t.Fatal("Should have returned single item")
	}

	if linkedCards[0].Issue != nil {
		t.Fatal("How did it found issue when no id was provided")
	}

}
func TestGetWithSingleCardWithCorrectIssueIdReturned(t *testing.T) {
	interactor := usecases.TimeLoggerInteractor{
		IssueIdExtractor: interfaces.SimpleRegexIssueIdExtractor{},
		YouTrackRepository: &youtrackRepositoryMock{
			findIssue: func(issueId string) (domain.YouTrackIssue, error) {
				if issueId != "MAT-123" {
					t.Fatal("Incorrect issue id")
				}
				return domain.YouTrackIssue{
					Fields: []domain.Field{
						{
							Name: "summary", Value: "this title is from youtrack"},
					},
				}, nil
			},
		},
		TrelloRepository: &trelloRepositoryMock{
			getAllCards: func(boardId string, doingListName string) []domain.TrelloCardDoingDuration {
				return []domain.TrelloCardDoingDuration{
					{
						Title:           "http://example.com/issue/MAT-123",
						Description:     "description",
						Duration:        123,
						Date:            time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
						YoutrackSummary: "something",
					},
				}
			},
		},
	}

	linkedCards := interactor.GetTrelloYouTrackCards("something", "another")

	if len(linkedCards) != 1 {
		t.Fatal("Should have returned single item")
	}

	if linkedCards[0].Issue == nil {
		t.Fatal("It should have found issue")
	}

	if len(linkedCards[0].Issue.Fields) < 1 {
		t.Fatal("It should have fields")
	}

}

func TestSaveWithZeroCards(t *testing.T) {
	interactor := usecases.TimeLoggerInteractor{
		IssueIdExtractor: interfaces.SimpleRegexIssueIdExtractor{},
		YouTrackRepository: &youtrackRepositoryMock{
			saveWorkLog: func(workLog domain.IssueWorkLog) error {
				t.Fatal("This should not have been called since nothing is being saved", workLog)
				return nil
			},
		},
		TrelloRepository: &trelloRepositoryMock{},
	}
	var logs []usecases.SaveTimeLog

	err := interactor.SaveWorklogs(logs)

	if err != nil {
		t.Fatal("No error expected", err)
	}

}

func TestSaveWithMultipleCards(t *testing.T) {
	var numberOfExpectedCalls = 2
	interactor := usecases.TimeLoggerInteractor{
		IssueIdExtractor: interfaces.SimpleRegexIssueIdExtractor{},
		YouTrackRepository: &youtrackRepositoryMock{
			saveWorkLog: func(workLog domain.IssueWorkLog) error {
				numberOfExpectedCalls--
				return nil
			},
		},
		TrelloRepository: &trelloRepositoryMock{},
	}

	logs := []usecases.SaveTimeLog{
		{
			Date:        time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
			Title:       "http://example.com/issue/MAT-123",
			Description: "lalalla",
			WorkType:    "Work",
			Duration:    123,
		},
		{
			Date:        time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
			Title:       "http://example.com/issue/MAT-456",
			Description: "lalalla",
			WorkType:    "Work",
			Duration:    123,
		},
	}

	err := interactor.SaveWorklogs(logs)

	if numberOfExpectedCalls != 0 {
		t.Fatal("Expected to call save 2 times")
	}

	fmt.Println(err)

}
func TestSaveWithInvalidCardTitle(t *testing.T) {
	interactor := usecases.TimeLoggerInteractor{
		IssueIdExtractor: interfaces.SimpleRegexIssueIdExtractor{},
		YouTrackRepository: &youtrackRepositoryMock{
			saveWorkLog: func(workLog domain.IssueWorkLog) error {
				return nil
			},
		},
		TrelloRepository: &trelloRepositoryMock{},
	}
	var logs = []usecases.SaveTimeLog{
		{
			Date:        time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
			Title:       "no issue id here",
			Description: "lalalla",
			WorkType:    "Work",
			Duration:    123,
		},
		{
			Date:        time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
			Title:       "no issue id here either",
			Description: "lalalla",
			WorkType:    "Work",
			Duration:    123,
		},
		{
			Date:        time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
			Title:       "http://example.com/issue/MAT-456",
			Description: "lalalla",
			WorkType:    "Work",
			Duration:    123,
		},
	}

	err := interactor.SaveWorklogs(logs)

	if _, ok := err.(usecases.NoIssueIdFound); !ok {
		t.Fatal("Expected to have card with no issue id found")
	}

	fmt.Println(err)

}

func TestSaveWhenRepositoryReturnsError(t *testing.T) {
	interactor := usecases.TimeLoggerInteractor{
		IssueIdExtractor: interfaces.SimpleRegexIssueIdExtractor{},
		YouTrackRepository: &youtrackRepositoryMock{
			saveWorkLog: func(workLog domain.IssueWorkLog) error {
				return errors.New("you didn't expect this. didnt you?")
			},
		},
		TrelloRepository: &trelloRepositoryMock{},
	}
	var logs = []usecases.SaveTimeLog{
		{
			Date:        time.Date(2018, 1, 1, 1, 1, 1, 1, time.UTC),
			Title:       "http://example.com/issue/MAT-456",
			Description: "lalalla",
			WorkType:    "Work",
			Duration:    123,
		},
	}

	interactor.SaveWorklogs(logs)

	// todo

}
