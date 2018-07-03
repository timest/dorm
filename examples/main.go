package main

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/Sirupsen/logrus"
    "dorm"
    "fmt"
)

var err error
var orm *dorm.Orm

var log = logrus.New()
// 用户(user) 可以发一个帖子（post）， 一个post可以有多个 留言（message），每个message都有User和post外键

func main() {
    orm, err = dorm.Open("mysql", "root:123456@/dorm")
    if err != nil {
        log.Fatal(err)
    }
    defer orm.Close()
    
    //createUser()
    //getUser()
    //createPost()
    //getPost()
    //createMessage()
    
     //queryUser()
    queryPost()
    queryOnePost()
}


type Queryset struct {
    where string
    table string
}

func queryPost() {
    var posts []Post
    orm.Query("name = 'I can go home'").Desc("id").All(&posts)
    for _, p := range posts {
        log.Info("pp: ", p.Id, p.Name, p.User.Age)
    }
    
}

func queryOnePost() {
    var post Post
    orm.Query("name = 'I can go home'").Desc("id").First(&post)
    log.Info("out: ", post)
}

func queryUser() {
    var users []User
    orm.Query("name = 'ccceo'").All(&users)
    log.Info("检索的结果是:", users)
    for _, u := range users {
        log.Info(u.Id)
    }
    // todo: age = 123
    // age < 12
    // age > 234
    // score == 23423.0234
    // score > 1234
    // score < 1234
    
}

func getPost() {
    p := new(Post)
    orm.Pk(p, 1)
    fmt.Println(p.User.Name)
}

func createUser() {
    u1 := new(User)
    orm.Defaults(u1) // 自动填充default
    orm.Create(u1)
}

func createPost() {
    u := new(User)
    orm.Pk(u, 4)
    
    p := new(Post)
    p.User = u
    p.Name = "this is name"
    
    orm.Create(p)
    
}

func getUser() {
    u := new(User)
 
    err = orm.Pk(u, 9)
    log.Info("U:", u, u.Id)
}

func createMessage() {
    u := new(User)
    orm.Pk(u, 6)
    
    p := new(Post)
    orm.Pk(p, 1)

    m := &Message {
        User: u,
        Post: p,
        Content: "this is content)",
    }
    orm.Create(m)
}