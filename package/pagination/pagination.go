package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultLimit  = 10
	DefaultOffset = 0
	DefaultPage   = 1
	DefaultSize   = 10
	MaxLimit      = 100
	MaxSize       = 100
)

type Pagination struct {
	Limit  int
	Offset int
	Page   int
	Size   int
}

func ParsePagination(r *http.Request) Pagination {
	p := Pagination{
		Limit:  DefaultLimit,
		Offset: DefaultOffset,
		Page:   DefaultPage,
		Size:   DefaultSize,
	}

	hasPage := false
	hasSize := false

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if pg, err := strconv.Atoi(pageStr); err == nil && pg > 0 {
			p.Page = pg
			hasPage = true
		}
	}

	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			if s > MaxSize {
				s = MaxSize
			}
			p.Size = s
			hasSize = true
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > MaxLimit {
				l = MaxLimit
			}
			p.Limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			p.Offset = o
		}
	}

	if hasPage || hasSize {
		if hasPage && hasSize {
			p.Offset = (p.Page - 1) * p.Size
			p.Limit = p.Size
		} else if hasPage {
			p.Offset = (p.Page - 1) * p.Size
			p.Limit = p.Size
		} else if hasSize {
			p.Limit = p.Size
		}
	}

	return p
}

func GetLimitOffset(r *http.Request) (limit, offset int) {
	p := ParsePagination(r)
	return p.Limit, p.Offset
}
