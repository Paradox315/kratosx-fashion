package option

import (
	"fmt"
	pb "kratosx-fashion/api/system/v1"
	"strings"
)

func Parse(opts ...*pb.QueryOption) (where string, order string, args []any) {
	var (
		qs []string
		os []string
	)
	for _, opt := range opts {
		q, o, arg := parseOpt(opt)
		qs = append(qs, q)
		os = append(os, o)
		arg = append(args, arg...)
	}
	where = strings.Join(qs, "AND")
	order = strings.Join(os, ",")
	return
}

func parseOpt(opt *pb.QueryOption) (query string, order string, args []any) {
	switch opt.Opt {
	case "EQ":
		query = fmt.Sprintf("%s = ?", opt.Field)
		args = append(args, opt.Value)
	case "NEQ":
		query = fmt.Sprintf("%s != ?", opt.Field)
		args = append(args, opt.Value)
	case "GT":
		query = fmt.Sprintf("%s > ?", opt.Field)
		args = append(args, opt.Value)
	case "GTE":
		query = fmt.Sprintf("%s >= ?", opt.Field)
		args = append(args, opt.Value)
	case "LT":
		query = fmt.Sprintf("%s < ?", opt.Field)
		args = append(args, opt.Value)
	case "LTE":
		query = fmt.Sprintf("%s <= ?", opt.Field)
		args = append(args, opt.Value)
	case "IN":
		query = fmt.Sprintf("%s IN (?)", opt.Field)
		args = append(args, opt.Value)
	case "BETWEEN":
		query = fmt.Sprintf("%s BETWEEN ? AND ?", opt.Field)
		args = append(args, opt.Interval.From, opt.Interval.To)
	case "LIKE":
		query = fmt.Sprintf("%s LIKE ?", opt.Field)
		args = append(args, opt.Value+"%")
	case "SORT":
		if opt.Desc {
			order = fmt.Sprintf("%s DESC", opt.Field)
		} else {
			order = fmt.Sprintf("%s ASC", opt.Field)
		}
	}
	return
}
