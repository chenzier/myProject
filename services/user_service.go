package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"product/datamodels"
	"product/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{UserRepository: repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {
	//检查密码是否匹配

	var err error
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	//这个函数接受一个用户密码作为输入，并返回一个经过哈希加密后的密码和err
	//在这个函数中，它使用了 bcrypt 哈希算法来生成密码的哈希值。
	//bcrypt.GenerateFromPassword 函数接受两个参数：
	//	用户提供的密码和哈希的计算成本（cost），它返回一个经过哈希加密的密码（作为字节数组）和err
	//	哈希的计算成本是一个整数，用于控制哈希计算的复杂度，也就是哈希函数的迭代次数，从而影响哈希结果的安全性。
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost) //使用了 bcrypt 哈希算法来实现密码的加密和验证
}

func ValidatePassword(userPassword string, hashed string) (isOk bool, err error) {
	//这个函数用于验证用户提供的密码和之前存储的哈希值是否匹配。
	//它接受两个参数：
	//	用户提供的密码
	//	之前存储的经过哈希加密的密码（哈希值）。
	//	在这个函数中，它使用了 bcrypt.CompareHashAndPassword 函数来比较用户提供的密码和之前存储的哈希值是否匹配。
	//	如果匹配，那么返回 true，否则返回 false 和一个相应的错误，通常是密码不匹配的错误信息
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误")
	}
	return true, nil
}
