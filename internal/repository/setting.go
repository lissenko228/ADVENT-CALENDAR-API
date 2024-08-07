package repository

type (
	SettingDTO struct {
		Month       int  `form:"month" validate:"required,min=1,max=12"`
		ShowAllDays bool `form:"showAllDays"`
	}
	Setting struct {
		ID          uint `json:"id"`
		Month       int  `json:"month" gorm:"not null"`
		ShowAllDays bool `json:"showAllDays" gorm:"type:boolean;not null;default:false"`
	}

	Clicks struct {
		IP string `json:"ip" gorm:"primaryKey;not null;unique"`
	}
)

var SettingService = new(Setting)

func (s *Setting) Get() (*Setting, error) {
	if err := DB.Model(s).First(s).Error; err != nil {
		return &Setting{}, err
	}

	return s, nil
}

func (s Setting) Update(toUpdate Setting) error {
	return DB.Model(&s).Where("1").Updates(&toUpdate).Error
}
