package dorm

import (
    "strconv"
    "database/sql"
    "reflect"
    "fmt"
)

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func logResult(o sql.Result) {
    a, _ := o.LastInsertId()
    b, _ := o.RowsAffected()
    log.Info("db.exec insert:", a, b)
}


func errCheck(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func getModelName(d reflect.Value) string {
    dtype := d.Type()
    return fmt.Sprintf("%s.%s", dtype.PkgPath(), dtype.Name())
}

// 自动填充strcut的default tag
func Defaults(m interface{}) {
    ps := reflect.ValueOf(m)
    if ps.Kind() != reflect.Ptr {
        return
    }
    s := ps.Elem()
    for i := 0; i < s.NumField(); i++ {
        f := s.Field(i)
        sf := s.Type().Field(i)
        d := sf.Tag.Get("default")
        if len(d) == 0 {
            continue
        }
        if f.IsValid() && f.CanSet() {
            switch f.Kind() {
            case reflect.String:
                f.SetString(d)
            case reflect.Bool:
                v, err := strconv.ParseBool(d)
                if err != nil {
                    log.Warn("%s 的default 非bool值")
                    continue
                }
                f.SetBool(v)
            case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                v, err := strconv.ParseUint(d, 10, 64)
                if err != nil {
                    log.Warn("%s 的default为非uint值")
                    continue
                }
                if !f.OverflowUint(v) {
                    f.SetUint(v)
                }
            case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                v, err := strconv.ParseInt(d, 10, 64)
                if err != nil {
                    log.Warn("%s 的default为非int值")
                    continue
                }
                if !f.OverflowInt(v) {
                    f.SetInt(v)
                }
            case reflect.Float32, reflect.Float64:
                v, err := strconv.ParseFloat(d, 64)
                if err != nil {
                    log.Warn("%s 的default为非float值")
                    continue
                }
                if !f.OverflowFloat(v) {
                    f.SetFloat(v)
                }
            default:
                log.Warn("ops, 没有考虑到:", f.Kind())
            }
            
        }
    }
    return
}
