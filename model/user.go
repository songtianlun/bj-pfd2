package model

import (
	"bj-pfd2/com/constvar"
	"bj-pfd2/com/db"
	"bj-pfd2/com/log"
	"bj-pfd2/com/utils"
	"fmt"
	"net/http"
)

type User struct {
	BaseModel
	Role     constvar.UserType `json:"role" gorm:"column:role;default:1"`
	Name     string            `json:"name" gorm:"column:name;not null;default:-" validate:"min=1,max=128"`
	Email    string            `json:"email" gorm:"column:email;unique" validate:"max=64"`
	Password string            `json:"password" gorm:"column:password;not null" validate:"min=5,max=128"`
	Sessions []Session
}

type Session struct {
	BaseModel
	Email  string `json:"email" gorm:"column:email" validate:"max=64"`
	UserId uint64 `gorm:"column:uid;comment:UserID" `
}

func (u *User) Create() error {
	u.Password = utils.Encrypt(u.Password)
	if u.Email == "" {
		return fmt.Errorf("email is empty")
	}
	return db.DB.Create(&u).Error
}

func (u *User) CloseAllSession() {
	var sessions []Session
	err := db.DB.Model(&u).Association("Sessions").Find(&sessions)
	if err != nil {
		log.ErrorF("Error to close all session of user %s: %s", u.Name, err.Error())
		return
	}
	if len(sessions) > 0 {
		for _, session := range sessions {
			err := session.Delete()
			if err != nil {
				log.ErrorF("Failed to close session %s of user %s: %s", session.Id, u.Name, err.Error())
			}
		}
	}
}

func (u *User) CreateSession() (session Session, err error) {
	session = Session{
		Email:  u.Email,
		UserId: u.Id,
	}
	err = db.DB.Model(&u).Association("Sessions").Append(&session)
	return
}

func (u *User) Session() (session Session, err error) {
	err = db.DB.First(&session, u.Id).Error
	return
}

func (s *Session) User() (user User, err error) {
	err = db.DB.First(&user, s.UserId).Error
	return
}

func (s *Session) Delete() error {
	return db.DB.Delete(&s).Error
}

func UserByEmail(email string) (u User, err error) {
	d := db.DB.Where("email = ?", email).First(&u)
	err = d.Error
	return
}

func Check(sid string) (s Session, err error) {
	if sid == "" {
		err = fmt.Errorf("get a null session id")
		log.Error(err.Error())
		return
	}
	err = db.DB.Where("uuid = ?", sid).First(&s).Error
	return
}

func CheckSession(request *http.Request) (session Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err != nil {
		return Session{}, err
	}
	session, err = Check(cookie.Value)
	if err != nil {
		return Session{}, err
	}
	return
}

func GetUser(r *http.Request) (u User, err error) {
	sess, err := CheckSession(r)
	if err != nil {
		return User{}, err
	}
	u, err = sess.User()
	if err != nil {
		return User{}, err
	}
	return
}

func DeleteSession(uuid string) {
	db.DB.Where("uuid = ?", uuid).Delete(&Session{})
}
