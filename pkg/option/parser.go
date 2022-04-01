package option

import (
	"fmt"
	pb "kratosx-fashion/api/system/v1"
	"strings"
	"time"
)

const timeFormat = `2006-01-02 15:04:05`

func Parse(opts ...*pb.QueryOption) (where string, order string, args []any) {
	var (
		qs []string
		os []string
	)
	for _, opt := range opts {
		q, o, arg := parseOpt(opt)
		if len(q) > 0 {
			qs = append(qs, q)
		}
		if len(o) > 0 {
			os = append(os, o)
		}
		if arg != nil && len(arg) > 0 {
			args = append(args, arg...)
		}
	}
	where = strings.Join(qs, " AND ")
	order = strings.Join(os, " , ")
	return
}

func parseOpt(opt *pb.QueryOption) (query string, order string, args []any) {
	switch opt.Opt {
	case "EQ":
		var arg any
		if opt.Time {
			arg, _ = time.ParseInLocation(timeFormat, opt.Value, time.Local)
		} else {
			arg = opt.Value
		}
		query = fmt.Sprintf("%s = ?", opt.Field)
		args = append(args, arg)
	case "NEQ":
		var arg any
		if opt.Time {
			arg, _ = time.ParseInLocation(timeFormat, opt.Value, time.Local)
		} else {
			arg = opt.Value
		}
		query = fmt.Sprintf("%s = ?", opt.Field)
		args = append(args, arg)
	case "GT":
		var arg any
		if opt.Time {
			arg, _ = time.ParseInLocation(timeFormat, opt.Value, time.Local)
		}
		query = fmt.Sprintf("%s = ?", opt.Field)
		args = append(args, arg)
	case "GTE":
		var arg any
		if opt.Time {
			arg, _ = time.ParseInLocation(timeFormat, opt.Value, time.Local)
		} else {
			arg = opt.Value
		}
		query = fmt.Sprintf("%s = ?", opt.Field)
		args = append(args, arg)
	case "LT":
		var arg any
		if opt.Time {
			arg, _ = time.ParseInLocation(timeFormat, opt.Value, time.Local)
		} else {
			arg = opt.Value
		}
		query = fmt.Sprintf("%s = ?", opt.Field)
		args = append(args, arg)
	case "LTE":
		var arg any
		if opt.Time {
			arg, _ = time.ParseInLocation(timeFormat, opt.Value, time.Local)
		} else {
			arg = opt.Value
		}
		query = fmt.Sprintf("%s = ?", opt.Field)
		args = append(args, arg)
	case "IN":
		var arg any
		if opt.Time {
			arg, _ = time.ParseInLocation(timeFormat, opt.Value, time.Local)
		} else {
			arg = opt.Value
		}
		query = fmt.Sprintf("%s = ?", opt.Field)
		args = append(args, arg)
	case "BETWEEN":

		var (
			arg1, arg2 any
		)
		if opt.Time {
			arg1, _ = time.ParseInLocation(timeFormat, opt.Interval.From, time.Local)
			arg2, _ = time.ParseInLocation(timeFormat, opt.Interval.To, time.Local)
		} else {
			arg1 = opt.Interval.From
			arg2 = opt.Interval.To
		}
		query = fmt.Sprintf("%s BETWEEN ? AND ?", opt.Field)
		args = append(args, arg1)
		args = append(args, arg2)
	case "LIKE":
		if opt.Time {
			break
		}
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
