
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

const hmacSecretKey = "hmac-secret-key"

var (
  ErrNotFound = errors.New("models: resource not found!")
  ErrInvalidID = errors.New("models: ID is invalid!")
  ErrInvalidPassword = errors.New("models: Password is invalid!")
  userPwPepper = "randomSecretString"
)

//Test: userGorm implements UserDB interface
var _ UserDB = &userGorm{}

type User struct{
  gorm.Model
  Name string
  Email string `gorm:"not null;unique_index"`
  Password string `gorm:"-"`
  PasswordHash string `gorm:"not null"`
  Remember string `gorm:"-"`
  RememberHash string `gorm:"not null;unique_index"`
}

//UserDB interface interacts with the database.
type UserDB interface{
  //Methods to query db
  ByID(int uint) (*User, error)
  ByEmail(int uint) (*User, error)
  ByRemember(int uint) (*User, error)

  //CRUD operations
  Create(user *User) error
  Update(user *User) error
  Delete(id uint) error
  //Close DB connection
  Close() error
  //Migration
  AutoMigrate() error
  DestructiveReset() error
}

type userGorm struct{
  db *gorm.DB
  hmac hash.HMAC
}

type userService{
  //Authenticate will verify if the provided email is valid. If valid,
  //the user is returned, otherwise will receive an error.
  Authenticate (email string) (*User, error)
  UserDB
}

func newUserGorm(connectionInfo string) (*userGorm, error){
  db, err := gorm.Open("postgres", connectionInfo)
  if err != nil{
    return nil, err
  }
  db.LogMode(true)
  hmac := hash.NewHMAC(hmacSecretKey)
  return &userGorm{
    db: db,
    hmac: hmac,
  }, nil
}

type userValidator struct{
  UserDB
}

type UserService struct{
  UserDB
}

func first(db *gorm.DB, dst interface{}) error{
  err := db.First(dst).Error
  if err == gorm.ErrRecordNotFound{
    return ErrNotFound
  }
  return err
}

func (ug *userGorm) ByID(id uint) (*User, error){
  var user User
  db := ug.db.Where("id = ?", id)
  err := first(db, &user)
  if err != nil{
    return nil, err
  }
  return &user, nil

}

func (ug *userGorm) ByRememberToken(token string) (*User, error){
  var user User
  rememberHash := ug.hmac.Hash(token)

  err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
  if err != nil{
    return nil, err
  }
  return &user, nil
}

func (ug *userGorm) ByEmail(email string) (*User, error){
  var user User
  db := ug.db.Where("email = ?", email)
  err := first(db, &user)
  return &user, err
}

func (ug *userGorm) AutoMigrate() error{
  if err := ug.db.AutoMigrate(&User{}).Error; err != nil{
    return err
  }
  return nil
}

func (ug *userGorm) DestructiveReset() error{
  err := ug.db.DropTableIfExists(&User{}).Error
  if err != nil{
    return err
  }
  ug.AutoMigrate()
  return nil
}

//Create user
func (ug *userGorm) Create(user *User) error{
  pwBytes := []byte(user.Password + userPwPepper)
  hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
  if err != nil{
    log.Panicln(err)
  }
  user.PasswordHash = string(hashedBytes)
  user.Password = ""

  if user.Remember == ""{
    token, err := rand.RememberToken()
    if err != nil {
      return err
    }
    user.Remember = token
  }
  user.RememberHash = ug.hmac.Hash(user.Remember)

  return ug.db.Create(user).Error
}

func (us *userService) Authenticate(email, password string) (*User, error){
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
func (ug *userGorm) Update(user *User) error{
  return ug.db.Save(user).Error
}

//Delete user
func (ug *userGorm) Delete(id uint) error{
  if id == 0{
    return ErrInvalidID
  }
  user := User{Model: gorm.Model{ID: id}}
  return ug.db.Delete(&user).Error
}

func NewUserService(connectionInfo string) (*UserService, error){
  ug, err := newUserGorm(connectionInfo)
  if err != nil{
    return nil, err
  }
  return &userService{
    UserDB: &userValidator{
      UserDB: ug,
    },
  },nil
}

func (ug *userGorm) Close() error{
  return ug.db.Close()
}
