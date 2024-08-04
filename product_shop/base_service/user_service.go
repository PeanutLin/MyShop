package base_service

import (
	"context"

	"productshop/product_shop/base_repository"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/mysql/gen/model"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	SelectUserByID(ctx context.Context, userID int64) (*model.User, error)
	SelectUserByName(ctx context.Context, userName string) (*model.User, error)
	AddUser(ctx context.Context, user *model.User) (int64, error)
	IsLoginSuccess(ctx context.Context, userName string, passord string) (user *model.User, isOK bool, err error)
}

type UserManagerService struct {
	UserRepository base_repository.IUserRepository
}

func NewUserService(repository base_repository.IUserRepository) IUserService {
	return &UserManagerService{
		UserRepository: repository,
	}
}

func (s *UserManagerService) SelectUserByID(ctx context.Context, userID int64) (*model.User, error) {
	panic("implement me")
}

func (s *UserManagerService) SelectUserByName(ctx context.Context, userName string) (*model.User, error) {
	panic("implement me")
}

func (s *UserManagerService) AddUser(ctx context.Context, user *model.User) (int64, error) {
	if user == nil || user.UserName == "" || user.Password == "" {
		return -1, errors.New("user information error")
	}
	enPwd, err := GeneratePwd(user.Password)
	if err != nil {
		return -1, errors.New("password error")
	}
	user.Password = string(enPwd)
	return s.UserRepository.Insert(ctx, user)
}

func (s *UserManagerService) IsLoginSuccess(ctx context.Context, userName string, uiPwd string) (user *model.User, isOK bool, err error) {
	if userName == "" || uiPwd == "" {
		return nil, false, errors.New("UserName or Password is empty")
	}

	user, err = s.UserRepository.SelectByName(context.Background(), userName)
	if err != nil {
		logs.Error("[SelectByName] error")
		return nil, false, errors.Wrap(err, "mysql error")
	}

	if isOK, _ = ValidatePwd(user.Password, uiPwd); isOK {
		return user, true, nil
	}
	return nil, false, nil
}

func ValidatePwd(dbPwd string, uiPwd string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(dbPwd), []byte(uiPwd)); err != nil {
		return false, err
	}
	return true, nil
}

func GeneratePwd(uiPwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(uiPwd), bcrypt.DefaultCost)
}
