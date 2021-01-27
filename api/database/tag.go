package database

import (
	"gorm.io/gorm"
)

type Tag struct {
	ID     string `gorm:"column:id;primary_key"`
	Name   string `gorm:"column:name"`
	UserID string `gorm:"column:user_id"`
	TodoID string `gorm:"column:todo_id"`
}

func (t *Tag) TableName() string { return "tag" }

type TagDao interface {
	InsertOne(u *Tag) error
	FindAll() ([]*Tag, error)
	FindOne(id string) (*Tag, error)
	FindByTodoID(todoID string) ([]*Tag, error)
}

type tagDao struct {
	db *gorm.DB
}

func NewTagDao(db *gorm.DB) TagDao {
	return &tagDao{db: db}
}

func (d *tagDao) InsertOne(u *Tag) error {
	res := d.db.Create(u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (d *tagDao) FindAll() ([]*Tag, error) {
	var tags []*Tag
	res := d.db.Find(&tags)
	if err := res.Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (d *tagDao) FindOne(id string) (*Tag, error) {
	var tags []*Tag
	res := d.db.Where("id = ?", id).Find(&tags)
	if err := res.Error; err != nil {
		return nil, err
	}
	if len(tags) < 1 {
		return nil, nil
	}
	return tags[0], nil
}

func (d *tagDao) FindByTodoID(todoID string) ([]*Tag, error) {
	var tags []*Tag
	res := d.db.Table("tag").
		Select("tag.*").
		Where("todo_id = ?", todoID).Find(&tags)

	if err := res.Error; err != nil {
		return nil, err
	}
	if tags == nil || len(tags) == 0 {
		return nil, nil
	}
	return tags, nil
}
