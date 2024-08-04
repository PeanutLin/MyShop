package repositories

import (
	"errors"
	"productshop/datamodels"
	"productshop/db"

	"gorm.io/gorm"
)

type IUserRepository interface {
	// 检查数据库连接
	Conn() error
	// 插入用户信息
	Insert(user *datamodels.User) (int64, error)
	// 删除用户信息
	Delete(int64) bool
	// 更新用户信息
	Update(user *datamodels.User) error
	// 根据主键查找用户
	SelectById(id int64) (user *datamodels.User, err error)
	// 根据用户名称查找用户
	Select(userName string) (user *datamodels.User, err error)
}

type UserManagerRepository struct {
	sqlConn *gorm.DB
}

// 新建用户操作接口
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserManagerRepository{
		sqlConn: db,
	}
}

// 检查数据库连接
func (u *UserManagerRepository) Conn() error {
	if u.sqlConn == nil {
		db, err := db.NewMysqlConn()
		if err != nil {
			return err
		}
		u.sqlConn = db
	}
	return nil
}

// 插入用户信息
func (u *UserManagerRepository) Insert(user *datamodels.User) (userID int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}

	err = u.sqlConn.Create(user).Error
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

// 删除用户信息
func (u *UserManagerRepository) Delete(int64) bool {
	// TODO:
	panic("implement me")
}

// 更新用户信息
func (u *UserManagerRepository) Update(user *datamodels.User) error {
	//TODO:
	panic("implement me")
}

// 根据用户名称查找用户
func (u *UserManagerRepository) Select(userName string) (*datamodels.User, error) {
	if userName == "" {
		return nil, errors.New("条件不能为空!")
	}

	if err := u.Conn(); err != nil {
		return nil, err
	}

	user := &datamodels.User{}
	err := u.sqlConn.Where("userName=?", userName).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// 根据主键查找用户
func (u *UserManagerRepository) SelectById(userID int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return nil, err
	}

	err = u.sqlConn.Where("ID=?", userID).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
