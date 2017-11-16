package dorm

import (
    "database/sql"
    "reflect"
    "fmt"
    "strings"
    "github.com/Sirupsen/logrus"
)

var log = logrus.New()

type Orm struct {
    DB *sql.DB
}

func Register(d interface{}) {
    val := reflect.ValueOf(d)
    ind := reflect.Indirect(val)
    mn := getModelName(ind)
    log.Info("Register ", mn)
    mi := &modelInfo{name: mn}
    parseField(ind, mi)
    modelCache.set(mn, mi)
}

func Open(driverName, dataSourceName string) (*Orm, error) {
    db, err := sql.Open(driverName, dataSourceName)
    if err != nil {
        return nil, err
    }
    err = db.Ping()
    if err != nil {
        return nil, err
    }
    return &Orm{
        DB: db,
    }, nil
}

func (d *Orm) Close() {
    d.DB.Close()
}

func getValueList(data interface{}, fList []string) []interface{} {
    val := reflect.ValueOf(data)
    var out []interface{}
    for _, k := range fList {
        v := val.Elem().FieldByName(k)
        switch k := v.Kind(); k {
        case reflect.Invalid:
            log.Info("invalid")
        case reflect.String:
            log.Info(v.String())
            out = append(out, v.String())
        case reflect.Float32, reflect.Float64:
            log.Info("float:", v.Float())
            out = append(out, v.Float())
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            log.Info("int:", v.Int())
            out = append(out, v.Int())
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            log.Info("uint:", v.Uint())
            out = append(out, v.Uint())
        default:
            log.Info("no match")
        }
    }
    return out
}

func (Orm *Orm) Create(data ...interface{}) {
    log.Info(data)
    for _, d := range data {
        val := reflect.ValueOf(d)
        ind := reflect.Indirect(val)
        //log.Info(val.Elem().fi)
        mn := getModelName(ind)
        log.Info(mn)
        mi := modelCache.get(mn)
        log.Info(mi.columns(true))
        cs := mi.columns(false)  // 不需要主键
        payload := make([]string, len(cs))
        for i, _ := range cs {
            payload[i] = "?"
        }
        log.Info(getValueList(d, cs))
        pl := getValueList(d, cs)
        //var vals []strings
        rawSql := fmt.Sprintf("insert into user(%s) values(%s)", strings.Join(cs, ","), strings.Join(payload, ","))
        log.Info(rawSql)
        o, err := Orm.DB.Exec(rawSql, pl...)
        errCheck(err)
        logResult(o)
    }
}

// 通过
func (o *Orm) Pk(d interface{}, id uint) error {
    mi := modelCache.getByInterface(d)
    log.Info(mi.fields)
    return nil
}

func (d *Orm) Retrieve(data interface{}, filter string) {
    
}

func (d *Orm) Update(data ...interface{}) {
    
}

func (d *Orm) Delete(data ...interface{}) {
    
}










