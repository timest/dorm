package dorm
// model and file
import (
    "reflect"
    "bytes"
    "fmt"
    "strings"
)

var dbfiled = map[string]string{
    "string": "varchar",
    "uint": "int",
    "uint8": "int",
    "uint16": "int",
    "uint32": "int",
    "uint64": "int",
    "int": "int",
    "int8": "int",
    "int16": "int",
    "int32": "int",
    "int64": "int",
    "float32": "double",
    "float64": "double",
}

// 存储model的信息
type modelInfo struct {
    name   string   // model的名字
    table  string   // 表格名称
    fields []*field // models的field集合
}


func (m *modelInfo) columns(pk bool) []string {
    var _fs []string
    for _, f := range m.fields {
        if pk == false && f.pk == true {
            // 不需要主键
            continue
        }
        _fs = append(_fs, f.name)
    }
    return _fs
}

func (m *modelInfo) fs(pk bool) []*field {
    var _fs []*field
    for _, f := range m.fields {
        if pk == false && f.pk == true {
            // 不需要主键
            continue
        }
        _fs = append(_fs, f)
    }
    return _fs
}

//存储register过的model信息
var modelCache = &_modelCache{
    cache: make(map[string]*modelInfo),
}

type _modelCache struct {
    cache map[string]*modelInfo
}

func (m *_modelCache) _get(name string) *modelInfo {
    for k := range m.cache {
        if k == name {
            return m.cache[k]
        }
    }
    return nil
}

func (m *_modelCache) getByValue(val reflect.Value) *modelInfo {
    if val.Kind() != reflect.Ptr {
        return nil
    }
    _struct := val.Elem()
    mn := getModelName(_struct)
    mi := m._get(mn)
    // 如果 cache里还没有注册过model，自动注册
    if mi == nil {
        mi = &modelInfo{name: mn, table: tableName(val)}
        parseField(val.Elem(), mi)
        debugInfo("成功注册model:", mi, "字段", mi.columns(true))
        modelCache.set(mn, mi)
    }
    return mi
}

func (m *_modelCache) get(d interface{}) *modelInfo {
    // 对象的地址，ptr
    val := reflect.ValueOf(d)
    return m.getByValue(val)
}

func tableName(r reflect.Value) string{
    var m reflect.Value
    m = r.MethodByName("TableName") // (x *Object)
    if m.IsValid() {
        return m.Call([]reflect.Value{})[0].String()
    }
    ind := reflect.Indirect(r)
    m = ind.MethodByName("TableName") // (x Object)
    if m.IsValid() {
        return m.Call([]reflect.Value{})[0].String()
    }
    t := getModelName(ind)
    return strings.ToLower(strings.Replace(t, ".", "_", -1))
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

// 遍历struct的field，将name 和 type 写入到mi的field里
// 如果遇到embed的struct，递归进去读取（如ID）
// 顺便处理field的tag属性，如 pk
func parseField(e reflect.Value, mi *modelInfo) {
    for i := 0; i < e.NumField(); i++ {
        fd := e.Field(i)
        // 如果遇到struct(如basemodel)就递归进去继续读取字段
        if fd.Kind() == reflect.Interface {
            continue
        }else if fd.Kind() == reflect.Struct {
            parseField(fd, mi)
        } else {
            sf := e.Type().Field(i)
            // 只处理能识别的类型 (目前是：string uint int float )  interface 直接continue
            if _, ok := dbfiled[sf.Type.Name()]; !ok {
                if strings.Index(sf.Type.String(), ".") != -1 { // todo: 有 . 就是外键吗？如xxx.xxx
                    var rel string = sf.PkgPath
                    if strings.HasPrefix(sf.PkgPath, "*") {
                        rel = rel[1:]
                    }
                    f := &field{
                        name: sf.Name + "_id",
                        typ: "uint64",
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
}

type field struct {
    name string
    typ  string
    pk   bool
    fk   bool  // foreignkey
    rel  string // 存放外键的 mn: modelname
}

func (f *field) String() string {
    return fmt.Sprintf("%s[ type:%s pk:%v ];", f.name, f.typ, f.pk)
}
