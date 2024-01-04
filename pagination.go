package pageable

import "math"

type Response struct {
	Content          interface{}      `json:"content"`
	Pageable         PaginationDetail `json:"pageable"`
	TotalElements    int              `json:"totalElements"`
	Last             bool             `json:"last"`
	TotalPages       int              `json:"totalPages"`
	First            bool             `json:"first"`
	Size             int              `json:"size"`
	Number           int              `json:"number"`
	Sort             SortDetail       `json:"sort"`
	NumberOfElements int              `json:"numberOfElements"`
	Empty            bool             `json:"empty"`
}

type PaginationDetail struct {
	Offset     int        `json:"offset"`
	PageNumber int        `json:"pageNumber"`
	PageSize   int        `json:"pageSize"`
	Paged      bool       `json:"paged"`
	Unpaged    bool       `json:"unpaged"`
	Sort       SortDetail `json:"sort"`
}

type SortDetail struct {
	Sorted   bool `json:"sorted"`
	Empty    bool `json:"empty"`
	Unsorted bool `json:"unsorted"`
}

func NewSortDetail(sort string) SortDetail {
	if sort == "" {
		return SortDetail{Sorted: false, Unsorted: true}
	}

	return SortDetail{
		Unsorted: false,
		Sorted:   true,
	}
}
func NewPaginatedResponse(data interface{}, dataCount int, totalCount int64, page int, limit int, sort string) *Response {
	pr := Response{
		Content:          data,
		Empty:            dataCount == 0,
		Number:           page,
		NumberOfElements: dataCount,
		Size:             limit,
		TotalElements:    int(totalCount),
		TotalPages:       int(math.Ceil(float64(totalCount) / float64(limit))),
		First:            page == 1,
		Last:             (int(totalCount) - (page * limit)) <= limit,
	}
	//Pagination detail
	pr.Pageable.PageNumber = page
	pr.Pageable.PageSize = limit
	pr.Pageable.Offset = (page - 1) * limit
	pr.Pageable.Paged = true
	pr.Pageable.Unpaged = false
	//Sort Detail
	pr.Sort = NewSortDetail(sort)
	pr.Pageable.Sort = NewSortDetail(sort)
	return &pr
}
