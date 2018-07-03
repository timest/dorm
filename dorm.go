package dorm

import (
    "database/sql"
    "reflect"
    "fmt"
    "strings"
    "github.com/Sirupsen/logrus"
    "strconv"
    "errors"
    "bytes"
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

func (o *Orm) Close() {
    o.DB.Close()
}

func (o *Orm) Query(d ...interface{}) *QuerySet {
    q := &QuerySet{
        orm: o,
    }
    return q.Query(d...)
}

func register(d interface{}) {
    val := reflect.ValueOf(d)
    ind := reflect.Indirect(val)
    mn := getModelName(ind)
    mi := &modelInfo{name: mn, table: tableName(val), typ: ind.Type()}
    parseField(ind, mi)
    modelCache.set(mn, mi)
}

// 遍历struct的field，将name 和 type 写入到mi的field里
// 如果遇到embed的struct，递归进去读取（如ID）
// 顺便处理field的tag属性，如 pk
func parseField(e reflect.Value, mi *modelInfo) {
    for i := 0; i < e.NumField(); i++ {
        fd := e.Field(i) // 字段内容
        // 如果遇到struct(如basemodel)就递归进去继续读取字段
        if fd.Kind() == reflect.Interface {
            continue
        } else if fd.Kind() == reflect.Struct {
            parseField(fd, mi)
        } else {
            sf := e.Type().Field(i) // 字段名称
            // 只处理能识别的类型 (目前是：string uint int float )  interface 直接continue
            if _, ok := dbfiled[sf.Type.Name()]; !ok {
                // 判断是否为外键
                if strings.Index(sf.Type.String(), ".") != -1 {
                    // todo: 有 . 就是外键吗？如xxx.xxx
                    var rel string = sf.PkgPath
                    if strings.HasPrefix(sf.PkgPath, "*") {
                        rel = rel[1:]
                    }
                    f := &field{
                        name: sf.Name,
                        typ: "ptr",
                        fk: true,
                        rel: rel,
                    }
                    // parseTag(f, &sf.Tag)  // todo: 外键应该没有特殊的 tag
                    mi.fields = append(mi.fields, f)
                }
                continue
            }
            f := &field{
                name: sf.Name,
                typ: sf.Type.Name(),
            }
            parseTag(f, &sf.Tag)
            mi.fields = append(mi.fields, f)
        }
    }
}

func parseTag(f *field, tag *reflect.StructTag) {
    t := tag.Get("orm")
    tList := strings.Split(t, ",")
    for _, k := range tList {
        switch k {
        case "pk":
            f.pk = true
        }
    }
    if fn := tag.Get("field"); fn != "" {
        f.name = fn
    }
}

func (o *Orm) Register(data ...interface{}) {
    for _, d := range data {
        register(d)
    }
}

type BaseModel struct {
    Id uint     `orm:"pk"`
}

// 通过数据库的字段名，查询struct的字段名
func structFieldName(ind reflect.Value, fieldName string) string {
    for i := 0; i < ind.NumField(); i++ {
        if ind.Type().Field(i).Name == fieldName {
            return fieldName
        }
        if ind.Type().Field(i).Tag.Get("field") == fieldName {
            return ind.Type().Field(i).Name
        }
    }
    return fieldName
}

func getValueList(data interface{}, fList []string) []interface{} {
    log.Info("为", fList, " 准备插入的数据")
    val := reflect.ValueOf(data) // data 是 &Object
    ind := reflect.Indirect(val)
    var out []interface{}
    for _, k := range fList {
        CheckFiled:
        var v reflect.Value
        // 如果有field tag的，用tag代替name
        v = ind.FieldByName(structFieldName(ind, k))
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
                k = k[0:len(k) - 3]
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

func (o *Orm) Create(data ...interface{}) {
    for _, d := range data {
        mi := modelCache.get(d)
        cols := mi.fields.cols(&fieldFilter{PK_NOT_NEED, RAW_NOT_NEED})  // 不需要主键
        payload := make([]string, len(cols))
        for i, _ := range cols {
            payload[i] = "?"
        }
        log.Info("Create: 待填充数据:", getValueList(d, cols))
        pl := getValueList(d, cols)
        rawSql := fmt.Sprintf("insert into %s(%s) values(%s)", mi.table, strings.Join(cols, ","), strings.Join(payload, ","))
        log.Info("query:", rawSql)
        out, err := o.DB.Exec(rawSql, pl...)
        errCheck(err)
        logResult(out)
    }
}

// 目前只支持单个修改
func (o *Orm) Update(d interface{}) error {
    //slice  := false
    ind := reflect.Indirect(reflect.ValueOf(d))
    typ := ind.Type()
    //if ind.Kind() == reflect.Slice {
    //    slice = true
    //    typ = ind.Type().Elem()
    //}
    mi := modelCache.getByName(typ.Name())
    cols := mi.fields.cols(&fieldFilter{PK_NOT_NEED, RAW_NOT_NEED})  // 不需要主键
    pl := getValueList(d, cols)
    if len(cols) != len(pl) || len(cols) == 0 {
        log.Info(len(cols), len(pl))
        return errors.New("字段数量不匹配")
    }
    var buf bytes.Buffer
    buf.WriteString("update " + mi.table + " set ")
    for i, k := range cols {
        //log.Info(pl[i].(type))
        //buf.WriteString(k + " = `" + pl[i].(string) + "`")
        if i != 0 {
            buf.WriteString(", ")
        }
        buf.WriteString(k + " = ")

        switch t := pl[i].(type) {
        case string:
            buf.WriteString(`"` + t +`"`)
        case int, int8, int16, int32:
            // todo
        case int64:
            buf.WriteString(strconv.FormatInt(pl[i].(int64), 10))
        case uint, uint8, uint16, uint32, uint64:
            buf.WriteString(strconv.Itoa(int(pl[i].(uint))))
        default:
            return errors.New(k + " 是什么类型?目前没有识别")
        }
    }
    v := ind.FieldByName("Id")
    buf.WriteString(" where id = " + strconv.FormatUint(v.Uint(), 10))
    log.Info("SQL:", buf.String())
    o.DB.Exec(buf.String())
    _, err := o.DB.Exec(buf.String())
    errCheck(err)
    return err
}

func (o *Orm) retrieve(d interface{}, q *QuerySet) error {
    slice  := false
    ind := reflect.Indirect(reflect.ValueOf(d))
    typ := ind.Type()
    if ind.Kind() == reflect.Slice {
        slice = true
        typ = ind.Type().Elem()
    }
    mi := modelCache.getByName(typ.Name())
    rawCols := mi.fields.cols(&fieldFilter{pk: PK_NEED, raw: RAW_NEED})
    //log.Info("query:", q.sql())
    rows, err := o.DB.Query(q.sql())
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        ref := make([]interface{}, len(rawCols))
        for i, _ := range ref {
            var t interface{}
            ref[i] = &t
        }
        var obj reflect.Value
        if slice {
            obj = reflect.New(ind.Type().Elem()).Elem()
        } else {
            obj = ind
        }
        
        rows.Scan(ref...)
        for i, r := range ref {
            xxx := obj.FieldByName(structFieldName(obj, rawCols[i]))
            v := reflect.Indirect(reflect.ValueOf(r)).Interface()
            switch mi.fields.index(i).typ {
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
            case "ptr":
                obj := reflect.New(modelCache.getByName(rawCols[i]).typ)
                _v, _ := strconv.ParseUint(string(v.([]uint8)), 10, 64)
                o.Pk(obj.Interface(), uint(_v))
                xxx.Set(obj)
            default:
                log.Info("ops，忘了考虑这个类型:", mi.fields.index(i).typ)
            }
        }
        if slice {
            ind.Set(reflect.Append(ind, obj))
        }
    }
    return nil
}

func (o *Orm) Pk(d interface{}, id uint) error {
    q := &QuerySet{
        orm: o,
        where: fmt.Sprintf("id = %d", id),
    }
    q.marshal(d)
    return o.retrieve(d, q)
}

// 自动填充strcut的default tag
func (o *Orm) Defaults(m interface{}) {
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

func logResult(o sql.Result) {
    if o == nil {
        return
    }
    a, _ := o.LastInsertId()
    b, _ := o.RowsAffected()
    log.Info("db.exec insert:", a, b)
}

func errCheck(err error) {
    if err != nil {
        log.Print(err)
    }
}

func getModelName(d reflect.Value) string {
    dtype := d.Type()
    return dtype.Name()
    //return fmt.Sprintf("%s.%s", dtype.PkgPath(), dtype.Name())
}







