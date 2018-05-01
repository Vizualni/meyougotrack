package interfaces

import (
	"time"

	"github.com/adlio/trello"
	"github.com/vizualni/meyougotrack/domain"
)

type TrelloAdlioClient struct {
	client *trello.Client
}

func (t *TrelloAdlioClient) GetAllCards(boardId string, doingListName string) []domain.TrelloCardDoingDuration {

	board, _ := t.client.GetBoard(boardId, trello.Defaults())

	cards, _ := board.GetCards(trello.Defaults())

	var doingCards []domain.TrelloCardDoingDuration

	for _, card := range cards {

		var cardDuration time.Duration = 0

		actions, _ := card.GetActions(trello.Arguments{
			"filter": "all",
		})

		durations, _ := actions.GetListDurations()

		for _, d := range durations {
			if d.ListName == doingListName {
				cardDuration += d.Duration
			}
		}

		doingCard := domain.TrelloCardDoingDuration{
			Title:       card.Name,
			Duration:    int64(cardDuration.Minutes()),
			Date:        card.CreatedAt(),
			Description: card.Desc,
		}

		doingCards = append(doingCards, doingCard)
	}

	return doingCards
}

func NewTrelloAdlioClient(apiKey string, apiToken string, boardId string, doingListName string) *TrelloAdlioClient {
	return &TrelloAdlioClient{
		client: trello.NewClient(apiKey, apiToken),
	}
}
