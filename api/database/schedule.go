package database

import (
	"gorm.io/gorm"
)

type Schedule struct {
	ID     string `gorm:"column:id;primary_key"`
	Title  string `gorm:"column:title"`
	UserID string `gorm:"column:user_id"`
}

func (u *Schedule) TableName() string {
	return "schedule"
}

type ScheduleDao interface {
	InsertOne(u *Schedule) error
	FindAll() ([]*Schedule, error)
	FindByUserID(userID string) ([]*Schedule, error)
	FindOne(id string) (*Schedule, error)
}

type scheduleDao struct {
	db *gorm.DB
}

func NewScheduleDao(db *gorm.DB) ScheduleDao {
	return &scheduleDao{db: db}
}

func (d *scheduleDao) InsertOne(u *Schedule) error {
	res := d.db.Create(u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}
func (d *scheduleDao) FindAll() ([]*Schedule, error) {
	var schedules []*Schedule
	res := d.db.Find(&schedules)
	if err := res.Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (d *scheduleDao) FindOne(id string) (*Schedule, error) {
	var schedules []*Schedule
	res := d.db.Where("id = ?", id).Find(&schedules)
	if err := res.Error; err != nil {
		return nil, err
	}
	if len(schedules) < 1 {
		return nil, nil
	}
	return schedules[0], nil
}

func (d *scheduleDao) FindByUserID(userID string) ([]*Schedule, error) {
	var schedules []*Schedule
	res := d.db.Where("user_id = ?", userID).Find(&schedules)
	if err := res.Error; err != nil {
		return nil, err
	}
	return schedules, nil
}
