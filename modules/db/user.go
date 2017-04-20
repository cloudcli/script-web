package db

import (
	"fmt"
	"script-web/modules/db/models"
	"github.com/jinzhu/gorm"
	//"script-web/modules/util"
)

func (this *MysqlDB) GetUserFullInfo(email string) (models.User, error){
	var user models.User

	rawQuery := fmt.Sprintf(`select * from users where email="%s"`, email)

	if err := this.Db.Raw(rawQuery).Scan(&user).Error; err != nil {
		return user, err
	}

	//rawQueryGroups := fmt.Sprintf(`select * from groups where user_id="%d"`, user.ID)

	//if err := this.Db.Raw(rawQueryGroups).Scan(&user.Groups).Error; err != nil {
	//	return user, err
	//}

	//rawQueryPackages := fmt.Sprintf(`select * from packages where user_id="%d"`, user.ID)

	//if err := this.Db.Raw(rawQueryPackages).Scan(&user.Packages).Error; err != nil {
	//	return user, err
	//}

	return user, nil
}

func (this *MysqlDB) FindUserByKey(key string, name string) (models.User, error) {
	var user models.User
	if err := this.Db.Where(key + " = ?", name).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (this *MysqlDB) IsUserNameExist(userName string) (bool, error) {

	var user models.User

	if err := this.Db.Find(&user, "name = ?", userName).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}


func (this *MysqlDB) IsUserEmailExist(email string) (bool, error) {

	var user models.User

	if err := this.Db.Find(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}


func (this *MysqlDB) CreateUser(email string, name string, password string) (*models.User, error) {
	user, err := models.NewUser()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := user.EncryptPassword(password)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	user.Email = email
	user.Name = name

	//
	//
	//enpassword, err := util.EncodePassword(password)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//user := &models.User {
	//	Name: name,
	//	Email: email,
	//	Password: enpassword,
	//}

	if err := this.Db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

