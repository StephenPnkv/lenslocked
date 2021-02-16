
package models

import (
  "errors"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "golang.org/x/crypto/bcrypt"
  "lenslocked/hash"
  "lenslocked/rand"
  "strings"
)

const hmacSecretKey = "hmac-secret-key"

var (
  ErrNotFound = errors.New("models: resource not found!")
  ErrInvalidID = errors.New("models: ID is invalid!")
  ErrInvalidPassword = errors.New("models: Password is invalid!")
  usrPasswordPepper = "randomSecretString"
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

//UserDB interface interacts with the database.
type UserDB interface{
  //Methods to query db
  ByID(id uint) (*User, error)
  ByEmail(email string) (*User, error)
  ByRemember(token string) (*User, error)

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

type UserService interface{
  //Authenticate will verify if the provided email is valid. If valid,
  //the user is returned, otherwise will receive an error.
  Authenticate (email, password string) (*User, error)
  UserDB
}

type userValidator struct{
  UserDB
  hmac hash.HMAC
}

type userService struct{
  UserDB
}
//Test: userGorm implements UserDB interface
var _ UserDB = &userGorm{}
type userValFn func(*User) error

func NewUserService(connectionInfo string) (UserService, error){
  ug, err := newUserGorm(connectionInfo)
  if err != nil{
    return nil, err
  }
  hmac := hash.NewHMAC(hmacSecretKey)
  uv := &userValidator{
    UserDB: ug,
    hmac: hmac,
  }
  return &userService{
    UserDB: uv,
  }, nil
}

func newUserGorm(connectionInfo string) (*userGorm, error){
  db, err := gorm.Open("postgres", connectionInfo)
  if err != nil{
    return nil, err
  }
  db.LogMode(true)
  return &userGorm{
    db: db,
  }, nil
}

func first(db *gorm.DB, dst interface{}) error{
  err := db.First(dst).Error
  if err == gorm.ErrRecordNotFound{
    return ErrNotFound
  }
  return err
}


// Validation Layer
func (uv *userValidator) Create(user *User) error{
  err := runUserValFns(user,
    uv.bcryptPassword,
    uv.hmacRemember,
    uv.setRememberIfUnset,
    uv.normalizeEmail)
  if err != nil{
    return err
  }

  return uv.UserDB.Create(user)

}

func (uv *userValidator) Update(user *User) error{
  if err := runUserValFns(user,
    uv.bcryptPassword,
    uv.hmacRemember,
    uv.normalizeEmail);
    err != nil{
      return err
  }
  return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error{
  var user User
  user.ID = id
  err := runUserValFns(&user, uv.idGreaterThan(0))
  if err != nil{
    return err
  }
  return uv.UserDB.Delete(id)
}

func (uv *userValidator) ByRemember(token string) (*User, error){
  user := User{
    Remember: token,
  }
  if err := runUserValFns(&user, uv.hmacRemember); err != nil{
    return nil, err
  }
  return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) hmacRemember(user *User) error{
  if user.Remember == ""{
    return nil
  }
  user.RememberHash = uv.hmac.Hash(user.Remember)
  return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error{
  //If the user's remember token is unset, a token is generated using
  //the rand package and is set.
  if user.Remember != ""{
    return nil
  }
  token, err := rand.RememberToken()
  if err != nil{
    return err
  }
  user.Remember = token
  return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFn{
  return userValFn(func(user *User) error{
    if user.id <= n{
      return ErrInvalidID
    }
    return nil
  })
}

func (uv *userValidator) normalizeEmail(email string) error{
  user.Email = strings.ToLower(user.Email)
  user.Email = strings.TrimSpace(user.Email)
  return nil
}

func (uv *userValidator) ByEmail(email string) (*User, error){
  user := User{
    Email: email,
  }
  err := runUserValFns(&user, uv.normalizeEmail)
  if err != nil{
    return nil, err
  }
  return uv.UserDB.ByEmail(user.Email)
}


// Database Interaction Layer
// Create method
// Data has been validated and is inserted into the database.
func (ug *userGorm) ByID(id uint) (*User, error){
  var user User
  db := ug.db.Where("id = ?", id)
  err := first(db, &user)
  if err != nil{
    return nil, err
  }
  return &user, nil

}

func (ug *userGorm) ByRemember(rememberHash string) (*User, error){
  var user User
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

  err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password + usrPasswordPepper))
  switch err{
  case nil:
    return foundUser, nil
  case bcrypt.ErrMismatchedHashAndPassword:
    return nil, ErrInvalidPassword
  default:
    return nil, err
  }

}

// Create user
func (ug *userGorm) Create(user *User) error{
  return ug.db.Create(user).Error
}

// Update user
func (ug *userGorm) Update(user *User) error{
  return ug.db.Save(user).Error
}

// Delete user
func (ug *userGorm) Delete(id uint) error{
  user := User{Model: gorm.Model{ID: id}}
  return ug.db.Delete(&user).Error
}

// Close database
func (ug *userGorm) Close() error{
  return ug.db.Close()
}



//Utility functions
func runUserValFns(user *User, fns ...userValFn) error {
  for _, fn := range fns{
    if err := fn(user); err != nil{
      return err
    }
  }
  return nil
}
