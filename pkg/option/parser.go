package option

import (
	"fmt"
	pb "kratosx-fashion/api/system/v1"
)

func Parse(opts ...*pb.QueryOption) (orders []string, keywords map[string]interface{}) {
	orders = make([]string, len(opts))
	keywords = make(map[string]interface{})
	for i, opt := range opts {
		if opt.Sort {
			orders[i] = fmt.Sprintf("%s desc", orders[i])
		} else {
			orders[i] = fmt.Sprintf("%s asc", orders[i])
		}
	}
	for _, opt := range opts {
		if opt.Keyword != "" {
			keywords[opt.Column] = opt.Keyword
		}
	}
	return
}
