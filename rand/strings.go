
package rand

import (
  "crypto/rand"
  "encoding/base64"
)

const RememberTokenBytes = 32

//Function returns a remember token
func RememberToken() (string, error){
  return String(RememberTokenBytes)
}

//Utility function that returns a slice of randomly generated bytes
//using the crypto/rand package.
func Bytes(n int) ([]byte, error){
  bytes := make([]byte,n)
  _, err := rand.Read(bytes)
  if err != nil{
    return nil, err
  }

  return bytes, nil
}

//Utility function that returns the number of bytes in the remember token.
func Nbytes(base64String string) (int, error){
  b, err := base64.URLEncoding.DecodeString(base64String)
  if err != nil{
    return -1, err
  }
  return len(b), nil
}

//Function generates a base64 encoded string
func String(nBytes int) (string, error){
  bytes, err := Bytes(nBytes)
  if err != nil{
    return "", err
  }
  return base64.URLEncoding.EncodeToString(bytes), nil
}
