package helpers

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
)

func Pagination(curPage int, totalPage int, showPages ...int) []int {
	pages := []int{}
	size := 9

	if len(showPages) > 0 {
		size = showPages[0]
	}

	if totalPage <= size {
		for i := 1; i <= totalPage; i++ {
			pages = append(pages, i)
		}

		return pages
	}

	half := float64(size) / 2.0
	mid := int(math.Ceil(half))
	step := int(math.Floor(half))

	if curPage <= mid {
		for i := 1; i <= mid+step; i++ {
			pages = append(pages, i)
		}

		return pages
	}

	if curPage+step >= totalPage {
		for i := totalPage - size + 1; i <= totalPage; i++ {
			pages = append(pages, i)
		}

		return pages
	}

	for i := curPage - mid + 1; i <= curPage+step; i++ {
		pages = append(pages, i)
	}

	return pages
}

func HttpBuildQueryUrl(uri string, query url.Values, page ...int) string {
	if len(page) > 0 {
		query.Add("page", strconv.Itoa(page[0]))
	}

	if len(query) == 0 {
		return uri
	}

	url := fmt.Sprintf("%s?%s", uri, query.Encode())

	return url
}
