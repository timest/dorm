# Orm
golang 的orm

## 使用

```go

// models.go 
type User struct {
    Name string     `default:"default"`
    Age uint16      `default:"18"`
    Score float64  `default:"11"`
}
func init() {
    Orm.Register(new(User))
}

// main.go
func main() {
    orm := new(Orm.Orm)
    db, err := sql.Open("mysql", "root:123456@/Orm")
    errCheck(err)
    errCheck(db.Ping())
    defer db.Close()
    orm.DB = db
    
    createUser(orm)
}

func createUser(o *Orm.Orm) {
    
    u := new(User)
    Orm.Defaults(u) // 自动填充default值
    o.Create(u)
}

```

## structtag
如果tag值有key-value的，直接写在tag里，如：`default:"timest" size:"255"`。

如果是单独的值，需要包含在`orm`里，如果多值用`,`分隔,如`orm:"unique,null,pk"`