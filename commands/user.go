package commands

import (
	"github.com/bookandmusic/docs/models"
)

func CreateSuperUser(username, password string) error {
	// 查询数据库是否存在该用户
	userModel := models.NewUser()
	user, _ := userModel.FindByIdentifier(username)
	// 如果用户不存在，创建新用户
	if user == nil {
		return userModel.Create(username, password, "", true)
	} else {

		updates := map[string]interface{}{
			"username": username,
			"password": password,
			"is_admin": true,
		}
		return user.Update(updates)
	}
}
