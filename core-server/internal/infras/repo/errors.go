package repo

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")
	ErrArticleNotFound  = errors.New("article not found")
	ErrCommentNotFound  = errors.New("comment not found")
)
