package models

import (
	"math"
)

// PaginationOptions comment
type PaginationOptions struct {
	Size   int
	Number int
}

// GetLimit comment
func (ctx *PaginationOptions) GetLimit() int {
	return int(ctx.Size)
}

// GetPageNumber comment
func (ctx *PaginationOptions) GetPageNumber() int {
	return int(Clamp(int64(ctx.Number-1), 0, int64(math.MaxInt16)))
}

// GetOffset comment
func (ctx *PaginationOptions) GetOffset() int {
	return ctx.GetPageNumber() * ctx.GetLimit()
}

// BuildDelta comment
func (ctx *PaginationOptions) BuildDelta(index int, total int) map[string]int {
	limit := ctx.GetLimit()
	offset := ctx.GetOffset()

	return map[string]int{
		"index":  int(index),
		"limit":  int(limit),
		"offset": int(offset),
		"total":  int(total),
	}
}

// NoDelta comment
func (ctx PaginationOptions) NoDelta() map[string]int {
	return map[string]int{
		"index":  0,
		"limit":  1,
		"offset": 0,
		"total":  1,
	}
}
