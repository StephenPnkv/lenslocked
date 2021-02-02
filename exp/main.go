
package main

import (
  _ "database/sql"
  "fmt"
  _ "github.com/lib/pq"
  "../models"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  dbname= "lenslocked_dev"
)


func main(){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

  us, err := models.NewUserService(psqlInfo)
  if err != nil{
    panic(err)
  }
  defer us.Close()
  us.DestructiveReset()

  user := models.User{
    Name: "Stephen Penkov",
    Email: "stephenpnkv@gmail.com",
  }
  if err := us.Create(&user); err != nil{
    panic(err)
  }

  usr, err := us.ByEmail("stephenpnkv@gmail.com")
  if err != nil{
    panic(err)
  }

  user.Name = "John Smith"
  if err := us.Update(&user); err != nil{
    panic(err)
  }

  err = us.Delete(1)
  if err != nil{
    panic(err)
  }

  usr, err = us.ByEmail("stephenpnkv@gmail.com")
  if err != nil{
    panic(err)
  }
  fmt.Println(usr)
}
