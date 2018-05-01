package interfaces

import (
	"io/ioutil"
	"net/http"

	"encoding/json"

	"time"

	"fmt"

	"github.com/vizualni/meyougotrack/usecases"
)

type Web struct {
	timeLogger usecases.TimeLogger
}

type saveTimeJson struct {
	Url         string    `json:"title"`
	Duration    int       `json:"duration"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Worktype    string    `json:"worktype"`
}

func (web Web) GetLoggableItems(boardId, listDoingName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		linkedCards := web.timeLogger.GetTrelloYouTrackCards(boardId, listDoingName)
		b, e := json.Marshal(linkedCards)
		fmt.Println(e)
		w.Write(b)
	})
}

func (web Web) Save() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jsonLogs []saveTimeJson

		body, _ := ioutil.ReadAll(r.Body)

		json.Unmarshal(body, &jsonLogs)

		var logs []usecases.SaveTimeLog

		for _, log := range jsonLogs {

			logs = append(logs, usecases.SaveTimeLog{
				WorkType:    log.Worktype,
				Title:       log.Url,
				Duration:    log.Duration,
				Date:        log.Date,
				Description: log.Description,
			})
		}

		err := web.timeLogger.SaveWorklogs(logs)

		if err != nil {
			fmt.Println(err)
		}

	})
}

func NewWeb(t usecases.TimeLogger) *Web {
	return &Web{
		timeLogger: t,
	}
}
