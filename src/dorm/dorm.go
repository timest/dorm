package dorm

import (
    "database/sql"
    "reflect"
    "fmt"
    "strings"
    "github.com/Sirupsen/logrus"
    "strconv"
)

var log = logrus.New()

type Orm struct {
    DB *sql.DB
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

func getField(d interface{}, name string) {
    val := reflect.ValueOf(d)
    ind := reflect.Indirect(val)
    field := ind.FieldByName("Id")
    log.Info("出现吧:", field)
    //for i := 0; i < ind.NumField(); i++ {
    //    log.Info("什么鬼？", ind.Field()
    //}
}
type BaseModel struct {
    Id uint     `orm:"pk"`
}
func getValueList(data interface{}, fList []string) []interface{} {
    debugInfo("为", fList, " 准备插入的数据")
    val := reflect.ValueOf(data) // data 是 &Object
    var out []interface{}
    for _, k := range fList {
        CheckFiled:
        v := reflect.Indirect(val).FieldByName(k)
        switch v.Kind() {
        case reflect.String:
            //log.Info("string:", v.String())
            out = append(out, v.String())
        case reflect.Float32, reflect.Float64:
            //log.Info("float:", v.Float())
            out = append(out, v.Float())
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            //log.Info("int:", v.Int())
            out = append(out, v.Int())
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            //log.Info("uint:", v.Uint())
            out = append(out, v.Uint())
        case reflect.Ptr: // 外键
            out = append(out, v.Elem().FieldByName("Id").Uint())
        default:
            if strings.Index(k, "_id") != -1 {
                k = k[0:len(k)-3]
                goto CheckFiled
            }
            log.Info("没有匹配的字段名:", k)
        }
    }
    if len(out) != len(fList) {
        log.Fatal("为字段", fList, "准备材料出现问题,只采集到了", out)
    }
    return out
}

func (Orm *Orm) Create(data ...interface{}) {
    for _, d := range data {
        mi := modelCache.get(d)
        cols := mi.columns(false)  // 不需要主键
        payload := make([]string, len(cols))
        for i, _ := range cols {
            payload[i] = "?"
        }
        log.Info("Create: 待填充数据:", getValueList(d, cols))
        pl := getValueList(d, cols)
        //var vals []strings
        rawSql := fmt.Sprintf("insert into %s(%s) values(%s)", mi.table, strings.Join(cols, ","), strings.Join(payload, ","))
        log.Info("query:", rawSql)
        o, err := Orm.DB.Exec(rawSql, pl...)
        errCheck(err)
        logResult(o)
    }
}



// 通过
func (o *Orm) Pk(d interface{}, id uint) error {
    mi := modelCache.get(d)
    cols := mi.columns(true)
    fList := mi.fs(true)  // field list
    if len(cols) != len(fList) {
        log.Fatal("怎么可能？")
    }
    //log.Info(mi.fields, cols)
    ref := make([]interface{}, len(cols))
    for i, _ := range ref {
        var t interface{}
        ref[i] = &t
    }
    query := fmt.Sprintf("select %s from %s where id = %d", strings.Join(cols, ","), mi.table, id)
    log.Info("query:", query)
    err := o.DB.QueryRow(query).Scan(ref...)
    if err != nil {
        log.Fatal(err)
    }
    //log.Info("ref:",ref)
    val := reflect.ValueOf(d)
    ind := reflect.Indirect(val)
    for i, r := range ref {
        xxx := ind.FieldByName(cols[i])
        v := reflect.Indirect(reflect.ValueOf(r)).Interface()
        //switch vv := v.(type) {
        //case []byte:
        //    log.Info("[]byte", string(vv))
        //case string:
        //    log.Info("string", string(vv))
        //}
        switch fList[i].typ {
        case "string":
            xxx.SetString(string(v.([]uint8)))
        case "uint", "uint16", "uint32", "uint64":
            _v, _ := strconv.ParseUint(string(v.([]uint8)), 10, 64)
            xxx.SetUint(_v)
        case "int", "int16", "int32", "int64":
            _v, _ := strconv.ParseInt(string(v.([]uint8)), 10, 64)
            xxx.SetInt(_v)
        case "float32", "float64":
            _v, _ := strconv.ParseFloat(string(v.([]uint8)), 64)
            xxx.SetFloat(_v)
        default:
            log.Info("ops，忘了考虑这个类型:", fList[i].typ)
        }
    }
    debugInfo("PK获取到的数据:", d)
    return nil
}

func (d *Orm) Retrieve(data interface{}, filter string) {
    
}

func (d *Orm) Update(data ...interface{}) {
    
}

func (d *Orm) Delete(data ...interface{}) {
    
}










