package repository

type (
	LoginDTO struct {
		Email    string `json:"email" form:"email" validate:"required,min=5,email"`
		Password string `json:"password" form:"password" validate:"required,min=5"`
	}

	User struct {
		ID           uint   `json:"id"`
		Email        string `json:"email" gorm:"unique;not null"`
		Password     string `json:"-" gorm:"not null;type:varchar(255)"`
		Role         string `json:"role" gorm:"not null; default:user"`
		RefreshToken string `json:"refreshToken" gorm:"not null"`
	}
)

var UserService = new(User)

func (u User) Get(where User) (User, error) {
	var user User
	err := DB.Where(&where).First(&user).Error
	return user, err
}

func (u User) Update(toUpdate User) error {
	return DB.Model(&u).Where(u).Updates(&toUpdate).Error
}
