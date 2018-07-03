package dorm

// model and file
import (
	"reflect"
	"bytes"
	"fmt"
	"strings"
)

var dbfiled = map[string]string{
	"string":  "varchar",
	"uint":    "int",
	"uint8":   "int",
	"uint16":  "int",
	"uint32":  "int",
	"uint64":  "int",
	"int":     "int",
	"int8":    "int",
	"int16":   "int",
	"int32":   "int",
	"int64":   "int",
	"float32": "double",
	"float64": "double",
}

// 存储model的信息
type modelInfo struct {
	name   string       // model的名字
	table  string       // 表格名称
	fields fieldTable   // models的field集合
	typ    reflect.Type //ind.Type()
}

func (mi *modelInfo) String() string {
	return fmt.Sprintf("[name:%s, tableName: %s, fields:%v, typ: %s]", mi.name, mi.table, mi.fields.cols(nil), mi.typ.Name())
}

//存储register过的model信息
var modelCache = &_modelCache{
	cache: make(map[string]*modelInfo),
}

type _modelCache struct {
	cache map[string]*modelInfo
}

func (m *_modelCache) getByName(name string) *modelInfo {
	for k := range m.cache {
		if k == name {
			return m.cache[k]
		}
	}
	return nil
}

func (m *_modelCache) get(d interface{}) *modelInfo {
	val := reflect.ValueOf(d)
	_struct := reflect.Indirect(val)
	mn := getModelName(_struct)
	mi := m.getByName(mn)
	if mi == nil {
		log.Fatalf("model %s 未被注册register!", mn)
	}
	return mi
}

func tableName(r reflect.Value) string {
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

// 每个字段一个field
type field struct {
	name  string
	typ   string
	pk    bool
	fk    bool   // foreignkey
	rel   string // 存放外键的 mn: modelname
}

func (f *field) String() string {
	return fmt.Sprintf("%s[ type:%s pk:%v ];", f.name, f.typ, f.pk)
}

type fieldTable []*field

type fieldFilter struct {
	pk  bool
	raw bool // 是否需要将 外键 的字段加上  _id 后缀, true
}

var PK_NEED = true
var RAW_NEED = true
var PK_NOT_NEED = false
var RAW_NOT_NEED = false

func (ft fieldTable) cols(ff *fieldFilter) []string {
	var fs []string
	for _, f := range ft {
		fname := f.name
		if ff == nil {
			// 如果 fieldFilter 为空，返回 包含ID 和 不含 _id 的外键名
			fs = append(fs, fname)
		} else {
			// 不需要主键
			if ff.pk == false && f.pk == true {
				continue
			}
			// 外键是否加上后缀
			if !ff.raw && f.fk {
				fs = append(fs, fname+"_id")
			} else {
				fs = append(fs, fname)
			}
		}
	}
	return fs
}

func (f fieldTable) index(i int) *field {
	if i < len(f) {
		return f[i]
	}
	return nil
}
