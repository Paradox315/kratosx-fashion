package pagination

func Parse(pageNum, pageSize uint32) (offset, limit int) {
	limit = int(pageSize)
	offset = int(pageNum-1) * limit
	return
}
