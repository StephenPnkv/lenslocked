
package models

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "golang.org/x/crypto/bcrypt"
  "lenslocked/hash"
  "lenslocked/rand"
  "strings"
  "regexp"
)

const hmacSecretKey = "hmac-secret-key"

var (
  //Custom errors for user validation
  ErrNotFound modelError = "models: resource not found."
  ErrInvalidID modelError = "models: ID is invalid."
  ErrPasswordInvalid modelError = "models: Password is invalid."
  ErrInvalidEmail modelError = "models: Invalid email."
  ErrEmailRequired modelError = "models: Email field is required."
  ErrEmailTaken modelError = "models: Email already taken."
  ErrPasswordTooShort modelError = "models: Password must be 8 characters or more."
  ErrPasswordRequired modelError = "models: A password is required."
  ErrRememberRequired modelError = "models: Remember token is required."
  ErrRememberTooShort modelError = "models: Remember token is too short."
  //Password salt
  usrPasswordPepper = "randomSecretString"
)

type modelError string

func (e modelError) Error() string{
  return string(e)
}

func (e modelError) Public() string{
  s := strings.Replace(string(e), "models: ", "", 1)
  split := strings.Split(s, " ")
  split[0] = strings.Title(split[0])
  return strings.Join(split, " ")
}

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
  ByID(ID uint) (*User, error)
  ByEmail(email string) (*User, error)
  ByRemember(token string) (*User, error)

  //CRUD operations
  Create(user *User) error
  Update(user *User) error
  Delete(ID uint) error
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
  emailRegex *regexp.Regexp
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
  uv := newUserValidator(ug, hmac)
  return &userService{
    UserDB: uv,
  }, nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
  return &userValidator{
    UserDB: udb,
    hmac: hmac,
    emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
  }
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
    uv.normalizeEmail,
    uv.emailFormat,
    uv.requireEmail,
    uv.emailAvailable,
    uv.passwordMinLength,
    uv.passwordRequired,)
  if err != nil{
    return err
  }

  return uv.UserDB.Create(user)

}

func (uv *userValidator) Update(user *User) error{
  if err := runUserValFns(user,
    uv.bcryptPassword,
    uv.hmacRemember,
    uv.normalizeEmail,
    uv.emailFormat,
    uv.requireEmail,
    uv.emailAvailable,
    uv.passwordMinLength,
    uv.passwordRequired,
    uv.passwordHashRequired);
    err != nil{
      return err
  }
  return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(ID uint) error{
  var user User
  user.ID = ID
  err := runUserValFns(&user, uv.idGreaterThan(0))
  if err != nil{
    return err
  }
  return uv.UserDB.Delete(ID)
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
    if user.ID <= n{
      return ErrInvalidID
    }
    return nil
  })
}

func (uv *userValidator) normalizeEmail(user *User) error{
  user.Email = strings.ToLower(user.Email)
  user.Email = strings.TrimSpace(user.Email)
  return nil
}

func (uv *userValidator) emailFormat(user *User) error{
  if user.Email == ""{
    return nil
  }
  if !uv.emailRegex.MatchString(user.Email){
    return ErrInvalidEmail
  }
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

func (uv *userValidator) requireEmail(user *User) error{
  if user.Email == ""{
    return ErrEmailRequired
  }
  return nil
}

func (uv *userValidator) emailAvailable(user *User) error{
  //If the record is not found, the email is available.
  //If the users ID does not match the ID in the database,
  //an error ErrEmailTaken is returned
  exist, err := uv.UserDB.ByEmail(user.Email)
  if err == ErrNotFound{
    return nil
  }

  if err != nil{
    return err
  }

  if user.ID != exist.ID{
    return ErrEmailTaken
  }
  return nil

}

func (uv *userValidator) bcryptPassword(user *User) error {
  if user.Password == ""{
    return nil
  }

  pwBytes := []byte(user.Password + usrPasswordPepper)
  hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
  if err != nil{
    return err
  }
  user.PasswordHash = string(hashedBytes)
  user.Password = ""
  return nil
}

func (uv *userValidator) passwordRequired(user *User) error{
  if user.Password == ""{
    return ErrPasswordRequired
  }
  return nil
}

func (uv *userValidator) passwordMinLength(user *User) error{
  if user.Password == ""{
    return nil
  }
  if len(user.Password) < 8{
    return ErrPasswordTooShort
  }
  return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error{
  if user.PasswordHash == ""{
    return ErrPasswordRequired
  }
  return nil
}

//If the remember token is less than 32 bytes, ErrRememberTooShort error is returned.
//If the remember token is empty, ErrRememberRequired error is returned.
func (uv *userValidator) rememberMinBytes(user *User) error{
  if user.Remember == ""{
    return ErrRememberRequired
  }

  n, err := rand.Nbytes(user.Remember)
  if err != nil{
    return err
  }
  if n < 32 {
    return ErrRememberTooShort
  }
  return nil
}


// Database Interaction Layer
// Create method
// Data has been validated and is inserted into the database.
func (ug *userGorm) ByID(ID uint) (*User, error){
  var user User
  db := ug.db.Where("ID = ?", ID)
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
  //If the password is invalid, it will return nil, ErrPasswordInvalid.
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
    return nil, ErrPasswordInvalid
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
func (ug *userGorm) Delete(ID uint) error{
  user := User{Model: gorm.Model{ID: ID}}
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
