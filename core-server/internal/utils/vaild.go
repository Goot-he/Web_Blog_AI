package utils

func Page(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func Size(size int) int {
	if size <= 0 {
		return 12
	}
	return size
}
