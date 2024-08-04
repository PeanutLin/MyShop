package base_repository

import (
	"context"
	"productshop/product_shop/middleware/mysql"
	"productshop/product_shop/middleware/mysql/gen/model"
	"productshop/product_shop/middleware/mysql/gen/query"
)

type IUserRepository interface {
	// 插入用户信息
	Insert(ctx context.Context, user *model.User) (userID int64, err error)
	// // 删除用户信息
	// Delete(int64) bool
	// 根据主键查找用户
	SelectByID(ctx context.Context, userID int64) (user *model.User, err error)
	// 根据用户名称查找用户
	SelectByName(ctx context.Context, userName string) (user *model.User, err error)
}

type UserRepository struct {
	Q *query.Query
}

// 新建用户操作接口
func NewUserRepository() IUserRepository {
	return &UserRepository{
		Q: mysql.QueryDB,
	}
}

// 插入用户信息
func (u *UserRepository) Insert(ctx context.Context, user *model.User) (userID int64, err error) {
	err = u.Q.User.WithContext(ctx).Save(user)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

// // 删除用户信息
// func (u *UserManagerRepository) Delete(int64) bool {
// 	// TODO:
// 	panic("implement me")
// }

// // 更新用户信息
// func (u *UserManagerRepository) Update(user *model.User) error {
// 	//TODO:
// 	panic("implement me")
// }

// 根据用户名称查找用户
func (u *UserRepository) SelectByName(ctx context.Context, userName string) (user *model.User, err error) {
	userPO := u.Q.User
	user, err = userPO.Where(userPO.UserName.Eq(userName)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 根据主键查找用户
func (u *UserRepository) SelectByID(ctx context.Context, userID int64) (user *model.User, err error) {
	userPO := u.Q.User
	user, err = userPO.Where(userPO.ID.Eq(int32(userID))).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}
