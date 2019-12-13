/*
 * 模仿orm.QuerySet;实现Filter, OrderBy, Limit方法
 * Usage
 * qs := new(QueryString)
 * qs.OrderBy("func_package", "-profile", "age")
 * qs.Filter("", "profile__icontains", "abc")
 * qs.Filter("call_dialplans", "age__icontains", 123)
 * var interfaceSlice []interface{} = make([]interface{}, len(params.FuncIds))
 * for i, d := range params.FuncIds {
 * 	   interfaceSlice[i] = d
 * }
 * queryString.Filter("c", "id__in", interfaceSlice...)
 * qs.OrderBy("cf", "-profile")
 * qs.Limit(5, 10)
 * qs.Limit(-1)
 * res := qs.String()
 * fmt.Printf("%#v\n", qs)
 * fmt.Println(res)
 */
package utils

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type QueryString struct {
	FilterStr  string
	OrderByStr string
	LimitStr   string
}

func (qs *QueryString) Filter(prefix string, expr string, args ...interface{}) {
	if len(qs.FilterStr) > 0 {
		qs.FilterStr += " and "
	}

	arr := strings.Split(expr, "__")
	if len(prefix) > 0 {
		arr[0] = prefix + "." + arr[0]
	}
	if len(arr) == 2 {
		switch arr[1] {
		case "exact":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " = '" + ConvertArgToString(arg) + "'"
			}
		case "iexact":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " ilike '" + ConvertArgToString(arg) + "'"
			}
		case "contains":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " like '%" + ConvertArgToString(arg) + "%'"
			}
		case "icontains":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " ilike '%" + ConvertArgToString(arg) + "%'"
			}
		case "startswith":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " like '" + ConvertArgToString(arg) + "%'"
			}
		case "istartswith":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " ilike '" + ConvertArgToString(arg) + "%'"
			}
		case "endswith":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " like '%" + ConvertArgToString(arg) + "'"
			}
		case "iendswith":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " ilike '%" + ConvertArgToString(arg) + "'"
			}
		case "in":
			for i, arg := range args {
				if i == 0 {
					qs.FilterStr += arr[0] + " in ('" + ConvertArgToString(arg)
					if i == len(args)-1 {
						qs.FilterStr += "')"
					}
				} else if i == len(args)-1 {
					qs.FilterStr += "','" + ConvertArgToString(arg) + "')"
				} else {
					qs.FilterStr += "','" + ConvertArgToString(arg)
				}
			}
		case "regex":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " ~ '" + ConvertArgToString(arg) + "'"
			}
		case "iregex":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " ~* '" + ConvertArgToString(arg) + "'"
			}
		case "nregex":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " !~ '" + ConvertArgToString(arg) + "'"
			}
		case "niregex":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " !~* '" + ConvertArgToString(arg) + "'"
			}
		case "gt":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " > '" + ConvertArgToString(arg) + "'"
			}
		case "gte":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " >= '" + ConvertArgToString(arg) + "'"
			}
		case "lt":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " < '" + ConvertArgToString(arg) + "'"
			}
		case "lte":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + " <= '" + ConvertArgToString(arg) + "'"
			}
		case "origin":
			for i, arg := range args {
				if i > 1 {
					qs.FilterStr += " and "
				}
				qs.FilterStr += arr[0] + ConvertArgToString(arg)
			}
		}
	} else if len(arr) == 1 {
		for i, arg := range args {
			if i > 1 {
				qs.FilterStr += " and "
			}
			qs.FilterStr += arr[0] + " = '" + ConvertArgToString(arg) + "'"
		}
	} else {
		qs.FilterStr += ""
	}
}

func (qs *QueryString) OrderBy(prefix string, exprs ...string) {
	if len(qs.OrderByStr) > 0 {
		qs.OrderByStr += ", "
	} else {
		qs.OrderByStr = " order by "
	}
	for index, expr := range exprs {
		if len(exprs) > 1 && index == len(exprs)-1 {
			qs.OrderByStr += ", "
		}
		if strings.HasPrefix(expr, "-") {
			if len(prefix) > 0 {
				qs.OrderByStr += prefix + "." + strings.TrimPrefix(expr, "-") + " desc"
			} else {
				qs.OrderByStr += strings.TrimPrefix(expr, "-") + " desc"
			}
		} else {
			if len(prefix) > 0 {
				qs.OrderByStr += prefix + "." + strings.TrimPrefix(expr, "-") + " asc"
			} else {
				qs.OrderByStr += strings.TrimPrefix(expr, "-") + " asc"
			}
		}
	}
}

func (qs *QueryString) Limit(limit interface{}, args ...interface{}) {
	aLimit := orm.ToInt64(limit)
	if aLimit < 1 {
		return
	}
	qs.LimitStr = " limit " + strconv.FormatInt(aLimit, 10) + " "
	if len(args) > 0 {
		b := orm.ToInt64(args[0])
		if b > 1 {
			qs.LimitStr += "offset " + strconv.FormatInt(b, 10)
		}
	}
}

func (qs *QueryString) String() string {
	return qs.FilterStr + qs.OrderByStr + qs.LimitStr
}

func ConvertArgToString(arg interface{}) (s string) {
	if a, ok := arg.(int); ok {
		s = strconv.Itoa(a)
	} else if a, ok := arg.(string); ok {
		s = a
	} else {
		s = ""
	}
	return
}
