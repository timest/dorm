# dorm
golang 的 ORM，除了mysql，未做其他数据库的兼容。大繁至简，步步为营。会有更多特性加入。

## 使用

```go

// models.go 
type User struct {
    Name string     `default:"default_name"`
    Age uint16      `default:"18"`
    Score float64  `default:"11"`
}

type Post struct {
    BaseModel
    User *User
    Name string
}

func init() {
    orm.Register(new(User), new(Post), new(Message))
}

// main.go
func main() {
	orm, err = dorm.Open("mysql", "root:123456@/dorm")
	if err != nil {
		log.Fatal(err)
	}
	defer orm.Close()
    
    // 创建单个object
	u := new(User)
	orm.Defaults(u) // 自动填充default
	orm.Create(u)
	
	// 创建含外键的object
	u := new(User)
	orm.Pk(u, 4) // 获取ID为4的User
	
	p := new(Post)
	p.User = u
	p.Name = "moyi is shabi"
	
	orm.Create(p)
	
	// 获取含外键的object
	p := new(Post)
	orm.Pk(p, 1)
	fmt.Println(p.Name)  // 打印出 User.name
	
	// 检索
	var posts []Post
	orm.Query("name = 'hello'").Desc("id").All(&posts)
	

}


```

## structtag
如果tag值有key-value的，直接写在tag里，如：`default:"timest" size:"255"`。

如果是单独的值，需要包含在`orm`里，如果多值用`,`分隔,如`orm:"unique,null,pk"`