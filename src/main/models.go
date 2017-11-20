package main

import (
    "fmt"
)

type BaseModel struct {
    Id uint     `orm:"pk"`
}

func (bm *BaseModel) FF() {
    fmt.Println("ff")
}

type User struct {
    BaseModel
    Name  string     `default:"default"`
    Age   uint16      `default:"18"`
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
    User    *User
    Post    *Post
    Content string `default:"没有填写短信内容"`
}

func (m *Message) TableName() string {
    return "message"
}

func init() {
    orm.Register(new(User), new(Post), new(Message))
}