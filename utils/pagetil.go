package utils

func GeneratePageNumbers(totalPages, page int) []int {
	var pages []int

	// 如果总页数小于等于3，返回所有页码
	if totalPages <= 3 {
		pages = make([]int, totalPages)
		for i := 0; i < totalPages; i++ {
			pages[i] = i + 1
		}
	} else {
		// 计算起始页码
		start := page
		if page+3 > totalPages {
			start = totalPages - 2
		}

		// 构建页码切片
		pages = make([]int, 3)
		for i := 0; i < 3; i++ {
			pages[i] = start + i
		}
	}

	if len(pages) == 0 {
		pages = []int{1}
	}

	return pages
}
