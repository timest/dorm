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
    fields []*field // models的field集合
}

func (m *modelInfo) columns(pk bool) []string {
    var fs []string
    for _, f := range m.fields {
        if pk == false && f.pk == true { // 不需要主键
            continue
        }
        fs = append(fs, f.name)
    }
    return fs
}

//存储register过的model信息
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

func (m *_modelCache) getByInterface(d interface{}) *modelInfo {
    val := reflect.ValueOf(d)
    if val.Kind() != reflect.Ptr {
        return nil
    }
    mn := getModelName(val.Elem())
    mi := m.get(mn)
    // 如果 cache里还没有注册过model，自动注册
    if mi == nil {
        mi = &modelInfo{name: mn}
        parseField(val.Elem(), mi)
        modelCache.set(mn, mi)
    }
    return mi
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

func parseField(e reflect.Value, mi *modelInfo) {
    for i := 0; i < e.NumField(); i++ {
        fd := e.Field(i)
        // 如果遇到struct(如basemodel)就递归进去继续读取字段
        if fd.Kind() == reflect.Struct {
            parseField(fd, mi)
        } else {
            sf := e.Type().Field(i)
            // 只处理能识别的类型 (目前是：string uint int float )  interface 直接continue
            if _, ok := dbfiled[sf.Type.Name()]; !ok {
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
}

func (f *field) String() string {
    return fmt.Sprintf("%s[ type:%s pk:%v ];", f.name, f.typ, f.pk)
}
