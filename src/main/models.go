package main

import (
    "reflect"
    "fmt"
)
type Model interface {
    FF(Model)
}

type BaseModel struct {
    Model
    Id uint     `orm:"pk"`
}

func (bm *BaseModel) FF(m Model) {
    val := reflect.ValueOf(m)
    ind := reflect.Indirect(val)
    log.Info(val, ind)
    for i := 0; i < ind.NumField(); i++ {
        log.Info(ind.Type().Field(i).Name)
    }
}
type User struct {
    BaseModel
    Name string     `default:"default"`
    Age uint16      `default:"18"`
    Score float64  `default:"11"`
}

func (u *User) TableName() string {
    return "user"
}

func (u *User) String() string {
    return fmt.Sprintf("Name: %s   Age: %d  Score: %f", u.Name, u.Age, u.Score)
}

type Post struct {
    BaseModel
    User *User
    Name string
}

func (p *Post) TableName() string {
    return "post"
}

type Message struct {
    BaseModel
    User *User
    Post *Post
    Content string `default:"没有填写短信内容"`
}

func (m *Message) TableName() string {
    return "message"
}

func init() {
    orm.Register(new(User), new(Post), new(Message))
}