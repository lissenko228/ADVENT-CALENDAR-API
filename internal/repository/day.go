package repository

import (
	"advent-calendar/pkg/utils"
)

type (
	Day struct {
		ID          uint         `json:"id"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		Attachments []Attachment `json:"attachments,omitempty"`
	}

	DayDTO struct {
		Title       string `json:"title" form:"title" validate:"required,min=5"`
		Description string `json:"description" form:"description" validate:"required,min=5"`
	}

	DayUPD struct {
		Title         string `json:"title" form:"title" validate:"min=5"`
		Description   string `json:"description" form:"description" validate:"min=5"`
		AttachmentIds []uint `json:"attachmentIds" form:"attachmentIds"`
	}
)

var DayService = new(Day)

func (d Day) Create(day DayDTO, files []utils.File) error {

	if len(files) > 0 {
		for _, file := range files {
			d.Attachments = append(d.Attachments, Attachment{
				Label: file.OriginalName,
				URL:   file.Destination,
				Type:  file.FileType,
			})
		}
	}

	d.Title = day.Title
	d.Description = day.Description

	return DB.Model(&d).Create(&d).Error
}

func (d Day) GetAll(params Params, where Day) ([]Day, error) {
	var days []Day

	query := DB.Model(&d).Where("id <= ?", where.ID).Preload("Attachments")

	if params.Limit > 0 {
		query = query.Limit(params.Limit).Offset((params.Page - 1) * params.Limit)
	}

	if err := query.Find(&days).Error; err != nil {
		return nil, err
	}

	return days, nil
}

func (d Day) Update(day DayDTO, files []utils.File) error {

	if len(files) > 0 {
		for _, file := range files {
			d.Attachments = append(d.Attachments, Attachment{
				Label: file.OriginalName,
				URL:   file.Destination,
				Type:  file.FileType,
			})
		}
	}

	d.Title = day.Title
	d.Description = day.Description

	return DB.Model(&d).Updates(&d).Error
}

func (d Day) Get(where Day) (Day, error) {
	if err := DB.Model(&d).Preload("Attachments").Where(where).First(&d).Error; err != nil {
		return Day{}, err
	}

	return d, nil
}