package dorm

import (
    "reflect"
    "bytes"
    "fmt"
    "strings"
)

type QuerySet struct {
    orm *Orm
    table string
    cols string
    where string
    limit string
    order string
}

func (q *QuerySet) Query(d ...interface{}) *QuerySet {
    if len(d) > 0 {
        //log.Info("filter ：", d, len(d))
        q.where = d[0].(string)
    }
    return q
}

func (q *QuerySet) First(d interface{}) error {
    q.setLimit(0, 1)
    q.marshal(d)
    return q.orm.retrieve(d, q)
}

func (q *QuerySet) Last(d interface{}) error {
    // todo
    return nil
}

func (q *QuerySet) All(d interface{}) error {
    q.marshal(d)
    return q.orm.retrieve(d, q)
}

func (q *QuerySet) marshal(d interface{}) {
    ind := reflect.Indirect(reflect.ValueOf(d))
    typ := ind.Type()
    if ind.Kind() == reflect.Slice {
        typ = ind.Type().Elem()
    }
    mi := modelCache.getByName(typ.Name())
    // 填充下col
    if q.cols == "" {
        cols := mi.fields.cols(&fieldFilter{pk: true})
        q.setCols(cols...)
    }
    if q.table == "" {
        q.setTableName(mi.table)
    }
}



func ptrCheck(d interface{}) {
    if reflect.ValueOf(d).Kind() != reflect.Ptr {
        panic("参数必须为Ptr")
    }
}

func sliceCheck(d interface{}) {
    val := reflect.ValueOf(d)
    ind := reflect.Indirect(val)
    if ind.Kind() != reflect.Slice {
        panic("参数必须为slice")
    }
}

func (q *QuerySet) Asc(cols ...string) *QuerySet {
    return q._order(true, cols...)
}

func (q *QuerySet) Desc(cols ...string) *QuerySet {
    return q._order(false, cols...)
}

func (q *QuerySet)_order(asc bool, cols ...string) *QuerySet {
    var buf bytes.Buffer
    fmt.Fprint(&buf, q.order)
    if asc {
        fmt.Fprintf(&buf, "%v Asc", strings.Join(cols, " Asc, "))
    } else {
        fmt.Fprintf(&buf, "%v Desc", strings.Join(cols, " Desc, "))
    }
    q.order = buf.String()
    return q
}

func (q *QuerySet)setCols(cols ...string) {
    var buf bytes.Buffer
    for i, c := range cols {
        buf.WriteString(fmt.Sprintf(" `%s`", c))
        if i != len(cols) - 1 {
            buf.WriteByte(',')
        }
    }
    q.cols = buf.String()
    return
}

func (q *QuerySet)setTableName(t string)  {
    if len(t) == 0 {
        panic("表名不能为空")
    }
    q.table = fmt.Sprintf("`%s`", t)
    return
}

func (q *QuerySet)setLimit(min, max uint) {
    q.limit = fmt.Sprintf("%d, %d", min, max)
    return
}

func (q *QuerySet) sql() string {
    var buf bytes.Buffer
    buf.WriteString("select ")
    if len(q.cols) == 0 {
        buf.WriteByte('*')
    } else {
        buf.WriteString(q.cols)
    }
    buf.WriteString(fmt.Sprintf(" from %s ", q.table))
    if len(q.where) > 0 {
        buf.WriteString(fmt.Sprintf("where %s ", q.where))
    }
    if len(q.order) > 0 {
        buf.WriteString(fmt.Sprintf("order by %s ", q.order))
    }
    if q.limit != "" {
        buf.WriteString(fmt.Sprintf("limit %s ", q.limit))
    }
    
    return buf.String()
}