package services

import (
	"errors"
	"fmt"
	"productshop/datamodels"
	"productshop/repositories"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetIdUser(userId int64) (*datamodels.User, error)
	GetNameUser(userName string) (*datamodels.User, error)
	AddUser(user *datamodels.User) (int64, error)
	IsLoginSuccess(userName string, pwd string) (*datamodels.User, bool)
}

type UserManagerService struct {
	UserRepository repositories.IUserRepository
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserManagerService{
		UserRepository: repository,
	}

}

func (s *UserManagerService) GetIdUser(userId int64) (*datamodels.User, error) {
	panic("implement me")
}

func (s *UserManagerService) GetNameUser(userName string) (*datamodels.User, error) {
	panic("implement me")
}

func (s *UserManagerService) AddUser(user *datamodels.User) (int64, error) {
	if user == nil || user.UserName == "" || user.HashPassword == "" {
		return 0, errors.New("用户信息不完整！")
	}
	enPwd, err := GeneratePwd(user.HashPassword)
	if err != nil {
		return 0, errors.New("密码异常！")
	}
	user.HashPassword = string(enPwd)
	return s.UserRepository.Insert(user)
}

func (s *UserManagerService) IsLoginSuccess(userName string, uiPwd string) (user *datamodels.User, isOk bool) {
	if userName == "" || uiPwd == "" {
		return nil, false
	}
	user, err := s.UserRepository.Select(userName)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	if isOk, _ := ValidatePwd(user.HashPassword, uiPwd); isOk {
		return user, true
	}
	return nil, false
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