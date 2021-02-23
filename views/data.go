package views

import (
  "log"
)

type Data struct{
  Alert *Alert
  Yield interface{}
}

type Alert struct{
  Message string
  Level string
}

const (
	AlertLvlError = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo = "info"
	AlertLvlSuccess = "success"
  AlertMsgGeneric = "Oops, something went wrong!" +
    "Please try again and contact us if the problem persists."
)

type PublicError interface{
  error
  Public() string
}

func (d *Data) SetAlert(err error){
  var msg string
  if pErr, ok := err.(PublicError); ok{
    msg = pErr.Public()
  }else{
    log.Println(err)
    msg = AlertMsgGeneric
  }
  d.Alert = &Alert{
    Level: AlertLvlError,
    Message: msg,
  }
}
