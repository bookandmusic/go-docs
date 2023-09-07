package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
)

type User struct {
	gorm.Model
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsAdmin     bool   `json:"is_admin"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) FindByIdentifier(identifier string) (*User, error) {
	var user User
	err := global.GVA_DB.Where("username = ? OR email = ?", identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) FindByUserId(userId int) (*User, error) {
	var user User
	err := global.GVA_DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Login 用于用户登录验证
func (u *User) Login(identifier, password string) (*User, bool) {
	user, err := u.FindByIdentifier(identifier)
	if err != nil {
		return nil, false
	}
	// 验证密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == nil {
		return user, true
	} else {
		return nil, false
	}
}

// Create 用于创建用户并存储到数据库
func (u *User) Create(username, password, email string, isAdmin bool) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		// 其他用户信息字段的赋值
		IsAdmin: isAdmin,
	}

	if err := global.GVA_DB.Create(&newUser).Error; err != nil {
		return err
	}

	return nil
}

// Update 用于更新用户字段，如果密码字段存在，则加密更新，否则只更新其他字段
func (u *User) Update(updates map[string]interface{}) error {
	// 检查是否有密码字段
	newPassword, newPasswordExists := updates["password"]
	if newPasswordExists {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword.(string)), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hashedPassword)
	}

	// 更新字段
	tx := global.GVA_DB.Begin()
	if err := tx.Model(u).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (u *User) CheckPassword(password string) bool {
	// 验证密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}
