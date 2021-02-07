
package models

import (
  "log"
  "errors"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "golang.org/x/crypto/bcrypt"
  "lenslocked/hash"
  "lenslocked/rand"
)

var (
  ErrNotFound = errors.New("models: resource not found!")
  ErrInvalidID = errors.New("models: ID is invalid!")
  ErrInvalidPassword = errors.New("models: Password is invalid!")
  userPwPepper = "randomSecretString"
)

type User struct{
  gorm.Model
  Name string
  Email string `gorm:"not null;unique_index"`
  Password string `gorm:"-"`
  PasswordHash string `gorm:"not null"`
  Remember string `gorm:"-"`
  RememberHash string `gorm:"not null;unique_index"`
}


type UserService struct{
  db *gorm.DB
}

func first(db *gorm.DB, dst interface{}) error{
  err := db.First(dst).Error
  if err == gorm.ErrRecordNotFound{
    return ErrNotFound
  }
  return err
}

func (us *UserService) ByID(id uint) (*User, error){
  var user User
  db := us.db.Where("id = ?", id)
  err := first(db,&user)
  if err != nil{
    return nil, err
  }
  return &user, nil

}

func (us *UserService) ByEmail(email string) (*User, error){
  var user User
  db := us.db.Where("email = ?", email)
  err := first(db, &user)
  return &user, err
}

func (us *UserService) AutoMigrate() error{
  if err := us.db.AutoMigrate(&User{}).Error; err != nil{
    return err
  }
  return nil
}

func (us *UserService) DestructiveReset() error{
  err := us.db.DropTableIfExists(&User{}).Error
  if err != nil{
    return err
  }
  us.AutoMigrate()
  return nil
}

//Create user
func (us *UserService) Create(user *User) error{
  pwBytes := []byte(user.Password + userPwPepper)
  hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
  if err != nil{
    log.Panicln(err)
  }
  user.PasswordHash = string(hashedBytes)
  user.Password = ""
  return us.db.Create(user).Error
}

func (us *UserService) Authenticate(email, password string) (*User, error){
  //This function authenticats a user with a given email and password
  //If the email is invalid, it will return nil, ErrNotFound.
  //If the password is invalid, it will return nil, ErrInvalidPassword.
  //If both email and password is valid, it will return user, nil.
  //If another error occurs, nil, error is returned.
  foundUser, err := us.ByEmail(email)
  if err != nil{
    return nil, err
  }

  err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password + userPwPepper))
  switch err{
  case nil:
    return foundUser, nil
  case bcrypt.ErrMismatchedHashAndPassword:
    return nil, ErrInvalidPassword
  default:
    return nil, err
  }

}

//Update user
func (us *UserService) Update(user *User) error{
  return us.db.Save(user).Error
}

//Delete user
func (us *UserService) Delete(id uint) error{
  if id == 0{
    return ErrInvalidID
  }
  user := User{Model: gorm.Model{ID: id}}
  return us.db.Delete(&user).Error
}

func NewUserService(connectionInfo string) (*UserService, error){
  db, err := gorm.Open("postgres", connectionInfo)
  if err != nil{
    return nil, err
  }
  db.LogMode(true)
  return &UserService{
    db: db,
  }, nil
}

func (us *UserService) Close() error{
  return us.db.Close()
}
