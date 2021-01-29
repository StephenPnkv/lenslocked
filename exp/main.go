
package main

import (
  "database/sql"
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

  usr, err := us.ByID(1)
  if err != nil{
    panic(err)
  }

  fmt.Println(usr)
}
