package dataloader

import (
	"fmt"
	"github.com/margostino/anfield/common"
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
	"log"
	"time"
)

func (a App) Consume() error {
	for {
		var match domain.Match
		err := a.kafka.Consume(&match)

		if err != nil || match.Metadata == nil {
			break
		}

		if match.Metadata.Finished {
			a.updateCompletion(&match)
		} else {
			a.updateCommentary(&match)
			a.updateAssets(&match)
		}
	}
	return nil // TODO: tbd
}

func (a App) updateAssets(match *domain.Match) {
	scores := a.scorer.CalculateScoring(match.Lineups, match.Data)

	// TODO: normalize key entity
	for name, score := range scores {
		var document domain.AssetDocument
		asset := &domain.Asset{
			Name:        name,
			Score:       score,
			LastUpdated: time.Now().UTC(),
		}
		filter, update := db.UpsertAssets(asset)
		a.db.Assets.Upsert(filter, update, &document)
	}
}

func (a App) updateCompletion(match *domain.Match) {
	var document domain.MatchDocument
	filter, update := db.UpsertMatchCompletion(match)
	a.db.Matches.Upsert(filter, update, &document)
	logging(&document)
}

func (a App) updateCommentary(match *domain.Match) {
	var document domain.MatchDocument
	filter, update := db.UpsertMatch(match)
	a.db.Matches.Upsert(filter, update, &document)
	logging(&document)
}

func logging(document *domain.MatchDocument) {
	id := common.GenerateEventID(document.Metadata.Url)
	dataLength := len(document.Data.Comments)
	message := fmt.Sprintf("New Message from %s with data length %d", id, dataLength)
	log.Println(message)
}
