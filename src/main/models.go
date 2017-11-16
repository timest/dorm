package main

import (
    "reflect"
    "dorm"
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

type Post struct {
    BaseModel
    User *User
    Name string
}

type Message struct {
    BaseModel
    User *User
    Post *Post
    Content string
}

func init() {
    dorm.Register(new(User))
}