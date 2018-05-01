package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"

	"io/ioutil"

	"time"

	"github.com/vizualni/meyougotrack/interfaces"
	_ "github.com/vizualni/meyougotrack/statik"
	"github.com/vizualni/meyougotrack/usecases"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Cannot read config file: %s", err))
	}

}

func main() {

	trelloApikey := viper.GetString("trello-api-key")
	trelloApiToken := viper.GetString("trello-api-token")
	trelloBoardId := viper.GetString("trello-board-id")
	trelloDoingListName := viper.GetString("trello-doing-list-name")

	youtrackApiKey := viper.GetString("youtrack-api-key")
	youtrackBaseUrl := viper.GetString("youtrack-base-url")

	fmt.Println(youtrackBaseUrl)
	trelloClient := interfaces.NewTrelloAdlioClient(
		trelloApikey,
		trelloApiToken,
		trelloBoardId,
		trelloDoingListName,
	)

	youtrackClient := interfaces.NewYouTrackClient(
		youtrackBaseUrl,
		youtrackApiKey,
	)

	timeLoggerInteractor := &usecases.TimeLoggerInteractor{
		YouTrackRepository: youtrackClient,
		TrelloRepository:   trelloClient,
		IssueIdExtractor:   interfaces.SimpleRegexIssueIdExtractor{},
	}

	web := interfaces.NewWeb(timeLoggerInteractor)

	//statikFS, err := fs.New()
	statikFS := http.Dir("./static")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.FileServer(statikFS))

	mux.HandleFunc("/", serveIndex(statikFS))
	mux.HandleFunc("/get-time", web.GetLoggableItems(trelloBoardId, trelloDoingListName))
	mux.HandleFunc("/save-time", web.Save())

	http.ListenAndServe(":8787", mux)
}

func serveIndex(system http.FileSystem) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		defer func() {
			fmt.Println(time.Now().Sub(now))
		}()

		indexFile, _ := system.Open("/index.html")
		defer indexFile.Close()

		indexFileContent, _ := ioutil.ReadAll(indexFile)

		w.Write(indexFileContent)
	})
}
