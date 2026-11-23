package dto

import (
	"core-server/internal/model/aggregate"
)

type CreateCommentRequest struct {
	ArticleID uint64
	UserID    uint64
	Content   string
}

type CreateReplyRequest struct {
	ArticleID uint64
	RootID    uint64
	UserID    uint64
	ReplyToID uint64
	Content   string
}

type DeleteCommentRequest struct {
	ID        uint64
	ArticleID uint64
	UserID    uint64
}

type DeleteReplyRequest struct {
	ID        uint64
	ArticleID uint64
	UserID    uint64
}

// =====================================================================================================================

type GetArticleCommentsRequest struct {
	ArticleID uint64
	Page      int
	Size      int
}

type GetArticleCommentsResponse struct {
	Comments []*aggregate.CommentAggregate
}

type GetCommentRepliesRequest struct {
	RootID uint64
	Page   int
	Size   int
}

type GetCommentRepliesResponse struct {
	Replies []*CommentInfo
}

type CommentInfo struct {
	Comment         *aggregate.CommentAggregate
	ReplyToUserName string
}
