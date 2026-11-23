package entity

import "time"

type InteractionFavor struct {
	ID        uint64    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
	UserID    uint64    `db:"user_id" json:"user_id"`
	ArticleID uint64    `db:"article_id" json:"article_id"`
	FolderID  uint64    `db:"folder_id" json:"folder_id"`
	Status    string    `db:"status" json:"status"`
	Version   int64     `db:"version" json:"version"`
}

const (
	FavorStatusFavored = "favored"
	FavorStatusNothing = "nothing"
)

type FavoriteFolder struct {
	ID           uint64
	UserID       uint64
	Name         string
	IsDefault    bool
	ArticleCount int64
}
