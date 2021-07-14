package dto

import "admin_demo/model"

type UserDto struct {
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}


/**
将User转换为UserDto
 */
func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}
