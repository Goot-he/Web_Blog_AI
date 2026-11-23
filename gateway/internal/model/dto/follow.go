package dto

type FollowRequest struct {
	UserID       uint64 `json:"userId"`
	TargetUserID uint64 `json:"targetUserId"`
}
type FollowListRequest struct {
	UserID   uint64
	Page     int
	PageSize int
}
