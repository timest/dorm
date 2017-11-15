package dorm

import (
    "database/sql"
    "reflect"
    "fmt"
    "bytes"
    "strings"
    "github.com/Sirupsen/logrus"
    "strconv"
)
var err error
var log = logrus.New()
var dbfiled = map[string]string{
    "string": "varchar",
    "uint16": "int",
    "float64": "double",
}

func init() {
    
}

type modelInfo struct {
    name string
    fields []*field
}

func (m *modelInfo) fieldList(d ...interface{}) []string {
    var fs []string
    for _, f := range m.fields {
        if _, ok := dbfiled[f.typ]; ok {
            fs = append(fs, f.name)
        }
        
    }
    return fs
}

var modelCache = &_modelCache{
    cache: make(map[string]*modelInfo),
}

type _modelCache struct {
    cache map[string]*modelInfo
}

func (m *_modelCache) get(name string) *modelInfo {
    for k := range m.cache {
        if k == name {
            return m.cache[k]
        }
    }
    return nil
}

func (m *_modelCache) set(name string, mi *modelInfo) bool {
    m.cache[name] = mi
    return true
}

func (m *_modelCache) String() string {
    var buf bytes.Buffer
    for k := range m.cache {
        buf.WriteString(fmt.Sprintf("[%s] { %s : %v }", k, m.get(k).name, m.get(k).fields))
    }
    return buf.String()
}

func Register(d interface{}) {
    val := reflect.ValueOf(d)
    ind := reflect.Indirect(val)
    mn := getModelName(ind)
    mi := &modelInfo{name: mn}
    for i := 0; i < ind.Type().NumField(); i++ {
        sf := ind.Type().Field(i)
        f := &field{
            name: sf.Name,
            typ: sf.Type.Name(),
        }
        mi.fields = append(mi.fields, f)
    }
    modelCache.set(mn, mi)
}

type field struct {
    name string
    typ string
}

func (f *field) String() string {
    return fmt.Sprintf("%s[ %s ];", f.name, f.typ)
}


type Orm interface {
    Create(...interface{})
    Retrieve(interface{}, string)
    Update(...interface{})
    Delete(...interface{})
}

type Dorm struct {
    DB *sql.DB
}

func getValueList(data interface{}, fList []string) []interface{} {
    val := reflect.ValueOf(data)
    var out []interface{}
    for _, k := range fList {
        v := val.Elem().FieldByName(k)
        switch k:= v.Kind(); k {
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
            out = append(out,v.Uint())
        default:
            log.Info("no match")
        }
    }
    return out
}

func (dorm *Dorm) Create(data ...interface{}) {
    log.Info(data)
    for _, d := range data {
        val := reflect.ValueOf(d)
        ind := reflect.Indirect(val)
        //log.Info(val.Elem().fi)
        mn := getModelName(ind)
        log.Info(mn)
        mi := modelCache.get(mn)
        log.Info(mi.fieldList())
        payload := make([]string, len(mi.fieldList()))
        for i, _ := range mi.fieldList() {
            payload[i] = "?"
        }
        log.Info(getValueList(d, mi.fieldList()))
        pl := getValueList(d, mi.fieldList())
        //var vals []strings
        rawSql := fmt.Sprintf("insert into user(%s) values(%s)", strings.Join(mi.fieldList(), ","), strings.Join(payload, ","))
        log.Info(rawSql)
        o, err := dorm.DB.Exec(rawSql, pl...)
        errCheck(err)
        logResult(o)
    }
}

func (d *Dorm) Retrieve(data interface{}, filter string) {
    
}

func (d *Dorm) Update(data ...interface{}) {
    
}

func (d *Dorm) Delete(data ...interface{}) {
    
}

var _ Orm = new(Dorm)




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
                f.SetUint(v)
            case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                v, err := strconv.ParseInt(d, 10, 64)
                if err != nil {
                    log.Warn("%s 的default为非int值")
                    continue
                }
                f.SetInt(v)
            case reflect.Float32, reflect.Float64:
                v, err := strconv.ParseFloat(d, 64)
                if err != nil {
                    log.Warn("%s 的default为非float值")
                    continue
                }
                f.SetFloat(v)
            default:
                log.Warn("ops, 没有考虑到:", f.Kind())
            }
            
        }
    }
    return
}









