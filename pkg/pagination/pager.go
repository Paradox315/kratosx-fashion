package pagination

func Parse(pageNum, pageSize uint32) (limit, offset int) {
	limit = int(pageSize)
	offset = int(pageNum-1) * limit
	return
}
