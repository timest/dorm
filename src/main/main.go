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
    getUser()
}

func createUser() {
    u1 := new(User)
    dorm.Defaults(u1) // 自动填充default
    orm.Create(u1)
}

func createPost() {
    
}

func getUser() {
    u := new(User)
    err = orm.Pk(u, 1)
    if err != nil {
        log.Fatal(err)
    }
}
