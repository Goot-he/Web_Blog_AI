package rpc

import (
	"context"
	"core-server/internal/model/aggregate"
	"time"

	"core-server/internal/application"
	"core-server/internal/model/dto"
	"core-server/internal/rpc/commentpb"
)

type CommentRPC struct {
	commentpb.UnimplementedCommentServiceServer
	CommentService *application.CommentService
}

func NewCommentRPC(commentService *application.CommentService) *CommentRPC {
	return &CommentRPC{CommentService: commentService}
}

func (c *CommentRPC) CreateComment(ctx context.Context, req *commentpb.CreateCommentReq) (*commentpb.CreateCommentResp, error) {
	resp, err := c.CommentService.CreateComment(ctx, &dto.CreateCommentRequest{
		ArticleID: req.GetArticleId(),
		Content:   req.GetContent(),
		UserID:    req.GetUserId(),
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.CreateCommentResp{Success: resp}, nil
}

func (c *CommentRPC) CreateReply(ctx context.Context, req *commentpb.CreateReplyReq) (*commentpb.CreateReplyResp, error) {
	resp, err := c.CommentService.CreateReply(ctx, &dto.CreateReplyRequest{
		ArticleID: req.GetArticleId(),
		RootID:    req.GetRootId(),
		UserID:    req.GetUserId(),
		Content:   req.GetContent(),
		ReplyToID: req.GetReplyToId(),
	})
	if err != nil {
		return nil, err
	}
	return &commentpb.CreateReplyResp{Success: resp}, nil
}

func (c *CommentRPC) DeleteComment(ctx context.Context, req *commentpb.DeleteCommentReq) (*commentpb.DeleteCommentResp, error) {
	if err := c.CommentService.DeleteComment(ctx, &dto.DeleteCommentRequest{
		ArticleID: req.GetArticleId(),
		ID:        req.GetId(),
		UserID:    req.GetUserId(),
	}); err != nil {
		return nil, err
	}
	return &commentpb.DeleteCommentResp{Success: true}, nil
}

func (c *CommentRPC) DeleteReply(ctx context.Context, req *commentpb.DeleteReplyReq) (*commentpb.DeleteReplyResp, error) {
	if err := c.CommentService.DeleteReply(ctx, &dto.DeleteReplyRequest{
		ID:        req.GetId(),
		ArticleID: req.GetArticleId(),
		UserID:    req.GetUserId(),
	}); err != nil {
		return nil, err
	}
	return &commentpb.DeleteReplyResp{Success: true}, nil
}

func (c *CommentRPC) GetArticleComments(ctx context.Context, req *commentpb.GetArticleCommentsReq) (*commentpb.GetArticleCommentsResp, error) {
	resp, err := c.CommentService.GetArticleComments(ctx, &dto.GetArticleCommentsRequest{
		ArticleID: req.GetArticleId(),
		Page:      int(req.Page),
		Size:      int(req.Size),
	})
	if err != nil {
		return nil, err
	}

	res := make([]*commentpb.ArticleCommentInfo, 0, len(resp.Comments))
	for _, comment := range resp.Comments {
		res = append(res, ToComment(comment))
	}
	return &commentpb.GetArticleCommentsResp{
		Comments: res,
	}, nil
}

func (c *CommentRPC) GetCommentReplies(ctx context.Context, req *commentpb.GetCommentRepliesReq) (*commentpb.GetCommentRepliesResp, error) {
	resp, err := c.CommentService.GetCommentReplies(ctx, &dto.GetCommentRepliesRequest{
		RootID: req.GetRootId(),
		Page:   int(req.Page),
		Size:   int(req.Size),
	})
	if err != nil {
		return nil, err
	}
	replies := make([]*commentpb.ReplyInfo, 0, len(resp.Replies))
	for _, reply := range resp.Replies {
		replies = append(replies, ToReply(reply))
	}
	return &commentpb.GetCommentRepliesResp{Replies: replies}, nil
}

// =====================================================================================================================

func ToComment(info *aggregate.CommentAggregate) *commentpb.ArticleCommentInfo {
	if info == nil {
		return nil
	}
	pb := &commentpb.ArticleCommentInfo{
		Id:         info.Comment.ID,
		ArticleId:  info.Comment.ArticleID,
		UserId:     info.Comment.UserID,
		RootId:     info.Comment.RootID,
		ReplyToId:  info.Comment.ReplyToID,
		Content:    info.Comment.Content,
		LikeCount:  info.Comment.LikeCount,
		ChildCount: info.Comment.ChildCount,
		CreatedAt:  info.Comment.CreatedAt.Format(time.DateTime),
		UserName:   info.Author.Name,
		UserAvatar: info.Author.Avatar,
	}
	return pb
}

func ToReply(info *dto.CommentInfo) *commentpb.ReplyInfo {
	if info == nil || info.Comment == nil || info.Comment.Comment == nil {
		return nil
	}
	c := info.Comment.Comment
	return &commentpb.ReplyInfo{
		Id:              c.ID,
		ArticleId:       c.ArticleID,
		UserId:          c.UserID,
		RootId:          c.RootID,
		ReplyToId:       c.ReplyToID,
		Content:         c.Content,
		LikeCount:       c.LikeCount,
		CreatedAt:       c.CreatedAt.Format(time.DateTime),
		UserName:        info.Comment.Author.Name,
		UserAvatar:      info.Comment.Author.Avatar,
		ReplyToUserName: info.ReplyToUserName,
	}
}
