
package main

import (
  _ "database/sql"
  "fmt"
  _ "github.com/lib/pq"
//  "lenslocked/models"
  "lenslocked/rand"
  "lenslocked/hash"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  dbname= "lenslocked_dev"
)


func main(){
  //fmt.Println(rand.String(10))
  //fmt.Println(rand.RememberToken())

  hmac := hash.NewHMAC("secret-key")
  fmt.Println(hmac.Hash("HashingThisString"))

}
