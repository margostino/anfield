package bot

import (
	"github.com/margostino/anfield/db"
	"github.com/margostino/anfield/domain"
)

// TODO: define transaction output
// TODO: evaluate async operation and return a promise
func (a App) buy(userId int64, key string, units int) {
	asset := a.getAsset(key)
	total := -1 * asset.Score * float64(units)
	a.updateWallet(userId, total)
}

// TODO: validate asset not found
func (a App) getAsset(key string) *domain.AssetDocument {
	filter := db.GetAssetsPatternFilter(key)
	return a.db.Assets.FindOneAsset(filter)
}

// TODO: validate user not found
func (a App) updateWallet(userId int64, total float64) {
	userFilter := db.GetUserFilter(userId)
	update := db.GetUpdateUser(total)
	a.db.Users.UpsertUser(userFilter, update)
}
