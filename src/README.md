# dorm
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
    dorm.Register(new(User))
}

// main.go
func main() {
    orm := new(dorm.Dorm)
    db, err := sql.Open("mysql", "root:123456@/dorm")
    errCheck(err)
    errCheck(db.Ping())
    defer db.Close()
    orm.DB = db
    
    createUser(orm)
}

func createUser(o *dorm.Dorm) {
    
    u := new(User)
    dorm.Defaults(u) // 自动填充default值
    o.Create(u)
}

```