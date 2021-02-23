
package model

import (
  "github.com/jinzhu/gorm"
)

type Gallery struct{
  gorm.Model
  UserID uint `gorm:"not_nul;index"`
  Title string `gorm:"not_null"`
}

//Gallery service provides database interactions with the Gallery resource
type GalleryService interface{
  GalleryDB
}

type GalleryDB interface{
  Create(gallery *Gallery) error
}

type galleryGorm struct{
  db *gorm.DB
}

func (gg *galleryGorm) Create(g *Gallery) error{
  return nil 
}
