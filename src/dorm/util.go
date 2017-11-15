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

