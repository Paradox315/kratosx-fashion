package pagination

import (
	"errors"
	pb "kratosx-fashion/api/system/v1"
)

func Parse(req *pb.ListRequest) (offset, limit int, err error) {
	if req == nil {
		err = errors.New("invalid page request")
		return
	}
	limit = int(req.PageSize)
	offset = int(req.PageNum-1) * limit
	return
}
