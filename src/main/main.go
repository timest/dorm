package main

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/Sirupsen/logrus"
    "dorm"
)

var err error
var orm *dorm.Orm

var log = logrus.New()
// 用户(user) 可以发一个帖子（post）， 一个post可以有多个 留言（message），每个message都有User和post外键
//
//func createPost(u *User, name string) {
//    o, err := db.Exec(`insert into post(user_id, name) values(?, ?)`, u.Id, name)
//    errCheck(err)
//    logResult(o)
//}
//
//func getUser(id int) (u *User, err error) {
//    u = new(User)
//    err = db.QueryRow(`select id, name, age, score from user where id = ?`, id).Scan(&u.Id, &u.Name, &u.Age, &u.Score)
//    return
//}
//
//func getPost(id int) (p *Post, err error) {
//    p = new(Post)
//    var user_id int
//    err = db.QueryRow(`select id, user_id, name from post where id = ?`, id).Scan(&p.Id, &user_id, &p.Name)
//    u, _ := getUser(user_id)
//    p.User = u
//    return
//}

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
    createMessage()
}

func getPost() {
    p := new(Post)
    orm.Pk(p, 1)
    
    //u1 := &User{
    //    Name: "test1",
    //}
    //val1 := reflect.ValueOf(u1)
    //ind1 := reflect.Indirect(val1)
    //log.Info(ind1.Type())
    ////name := "main.user"
    //
    //v := reflect.New(reflect.TypeOf(*u1))
    //v2 := v.Elem().Interface().(User)
    //
    //log.Info(v,v2)
    //log.Info(u1)
    //
    ////v3 := reflect.New(val1.Type())
    ////log.Info(v3)
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
    p.Name = "moyi is shabi"
    
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
    
    //m := new(Message)
    //orm.Defaults(m)
    //m.User = u
    //m.Post = p
    //
    //orm.Create(m)
    m := &Message {
        User: u,
        Post: p,
        Content: "张信哲（再多的苦我也愿意背)",
    }
    orm.Create(m)
}