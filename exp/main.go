
package main

import (
  _ "database/sql"
  "fmt"
  _ "github.com/lib/pq"
  "lenslocked/models"
  "log"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  dbname= "lenslocked_dev"
)

func logError(err error){
  if err != nil {
    log.Fatal(err)
  }
}

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
    Password: "abc123",
  }
  err = us.Create(&user)
  logError(err)

  fmt.Printf("%+v\n", user)
  if user.Remember == ""{
    panic("Invalid remember token.")
  }

  //check remember token
  user2, err := us.ByRememberToken(user.Remember)
  logError(err)
  fmt.Printf("%+v\n", *user2)

}
