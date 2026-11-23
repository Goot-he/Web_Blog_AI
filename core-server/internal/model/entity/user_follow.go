package entity

import "time"

type UserFollow struct {
	ID          uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FollowerID  uint64
	FollowingID uint64
	Status      string
	Version     int64
}

const (
	FollowStatusFollowing = "following"
	FollowStatusNothing   = "nothing"
)
