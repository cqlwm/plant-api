package user

import (
	"errors"
	"fmt"
	"log"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
	"plant-api/net/one/test"
	"plant-api/net/one/util"
)

type UserServiceImpl struct {
}

// 登录验证
func (usi *UserServiceImpl) LoginCheck(form *entry.LoginForm) (string, error) {
	// 获取wxId
	id, _, b := util.WeiXinAuth(form.Code)
	if !b {
		return "", errors.New(config.WX_CODE_ERROR.Message())
	}

	userDo := entry.SysUser{}
	ok := util.BeanTo(*form, &userDo)
	userDo.OpenId = id
	userDo.UUID = util.Uuid(id)

	if !ok {
		log.Println(config.JSON_SERIALIZE_ERROR)
		return "", errors.New(config.SERVER_ERROR.Message())
	}

	// find
	db := config.DataBase
	tx := db.Find(&entry.SysUser{
		OpenId: id,
	})

	if tx.RowsAffected == 0 {
		tx2 := db.Save(&userDo)
		if tx2.Error != nil {
			log.Println(tx.Error)
			// 保存错误
			return "", errors.New(config.SERVER_ERROR.Message())
		}
		tx2.Commit()
	} else {
		userItem := entry.SysUser{}
		rows, _ := tx.Rows()
		if rows.Next() {
			_ = db.ScanRows(rows, &userItem)
		}
		userDo.Id = userItem.Id
	}

	fmt.Println(userDo.Id)
	if userDo.Id == 0 {
		log.Println("未知错误 UserId == 0")
		return "", errors.New("未知错误 UserId == 0")
	}

	// 生成Token或续期Token
	token := test.BuildToken(userDo.Id)
	return token, nil
}
