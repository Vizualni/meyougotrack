package usecases

import (
	"time"

	"github.com/vizualni/meyougotrack/domain"
)

type TrelloYouTrackLink struct {
	Card  *domain.TrelloCardDoingDuration `json:"card"`
	Issue *domain.YouTrackIssue           `json:"issue"`
}

type TrelloRepository interface {
	GetAllCards(boardId string, doingListName string) []domain.TrelloCardDoingDuration
}

type YouTrackRepository interface {
	FindIssueByIssueId(issueId string) (domain.YouTrackIssue, error)
	SaveWorkLog(workLog domain.IssueWorkLog) error
}

type IssueIdExtractor interface {
	Extract(input string) (string, error)
}

type TimeLoggerInteractor struct {
	TrelloRepository   TrelloRepository
	IssueIdExtractor   IssueIdExtractor
	YouTrackRepository YouTrackRepository
}

type SaveTimeLog struct {
	Title       string
	Duration    int
	Date        time.Time
	WorkType    string
	Description string
}

type TimeLogger interface {
	GetTrelloYouTrackCards(boardId string, listDoingName string) []TrelloYouTrackLink
	SaveWorklogs(logs []SaveTimeLog) error
}

func (t *TimeLoggerInteractor) GetTrelloYouTrackCards(boardId string, listDoingName string) []TrelloYouTrackLink {
	cards := t.TrelloRepository.GetAllCards(boardId, listDoingName)

	var linked []TrelloYouTrackLink

	for index := range cards {
		cardTitle := cards[index].Title

		// doesnt really matter if we cannot find exact id from trello title
		youtrackIssueId, err := t.IssueIdExtractor.Extract(cardTitle)

		var youtrackIssue *domain.YouTrackIssue = nil

		if err == nil {
			yt, err := t.YouTrackRepository.FindIssueByIssueId(youtrackIssueId)
			if err != nil {
				panic("issue with that id not found")
			}
			youtrackIssue = &yt
		}

		link := TrelloYouTrackLink{
			Card:  &cards[index],
			Issue: youtrackIssue,
		}

		linked = append(linked, link)
	}

	return linked

}

type NoIssueIdFound struct {
	logs []SaveTimeLog
}

func (NoIssueIdFound) Error() string {
	return "no issue id found"
}

func (t *TimeLoggerInteractor) SaveWorklogs(logs []SaveTimeLog) error {

	// error storage for issues with no id
	noIssueIdsFound := NoIssueIdFound{
		logs: []SaveTimeLog{},
	}

	var issueWorkLogs []domain.IssueWorkLog

	// extract issue id and save it to an array
	// otherwise add error to array
	for _, log := range logs {
		issueId, err := t.IssueIdExtractor.Extract(log.Title)

		if err != nil {
			noIssueIdsFound.logs = append(noIssueIdsFound.logs, log)
			continue
		}

		workLog := domain.IssueWorkLog{
			Date:        log.Date,
			Duration:    log.Duration,
			Type:        log.WorkType,
			Description: log.Description,
			IssueId:     issueId,
		}

		issueWorkLogs = append(issueWorkLogs, workLog)
	}

	// todo: testovi
	// if len(noIssueIdsFound.logs) > 0 {
	// 	return noIssueIdsFound
	// }

	// saving logs to repository
	for _, log := range issueWorkLogs {

		// skip log with zero duration
		if log.Duration <= 0 {
			continue
		}

		err := t.YouTrackRepository.SaveWorkLog(log)

		if err != nil {
			return err
		}
	}

	return nil
}
