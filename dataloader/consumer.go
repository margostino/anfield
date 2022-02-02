package dataloader

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"log"
)

func (a App) Consume() error {
	for {
		message, err := a.kafka.ReadMessage()

		if err != nil {
			break
		}

		if message.Metadata.Finished {
			a.upsertCompletion(message)
		} else {
			a.upsertCommentary(message)
			a.upsertAssets(message)
		}
	}
	return nil // TODO: tbd
}

func (a App) upsertAssets(message *domain.Message) {
	scores := a.scorer.CalculateScoring(message.Lineups, message.Data)

	// TODO: normalize key entity
	for key, value := range scores {
		filter := db.GetAssetsFilter(key)
		update := db.GetUpdateAssets(key, value)
		a.db.Assets.UpsertAsset(filter, update)
	}

}

func (a App) upsertCompletion(message *domain.Message) {
	filter := db.GetUrlFilter(message.Metadata.Url)
	update := db.GetUpdateCompletion(message)
	document := a.db.Matches.UpsertMatch(filter, update)
	logging(document)
}

func (a App) upsertCommentary(message *domain.Message) {
	filter := db.GetUrlFilter(message.Metadata.Url)
	update := db.GetUpdateCommentary(message)
	document := a.db.Matches.UpsertMatch(filter, update)
	logging(document)
}

func logging(document *domain.MatchDocument) {
	id := common.GenerateEventID(document.Metadata.Url)
	dataLength := len(document.Data.Comments)
	message := fmt.Sprintf("New Message from %s with data length %d", id, dataLength)
	log.Println(message)
}
