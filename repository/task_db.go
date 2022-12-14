package repository

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type taskRepositoryDB struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return taskRepositoryDB{db: db}
}

// Get Task
func (r taskRepositoryDB) GetAll(search string, sort string) ([]Task, error) {
	tx := r.db
	if search != "" { // search
		tx = tx.Where("INSTR(CONCAT_WS('|', title, description, status), ?)", search)
	}
	if sort == "title" || sort == "created_at" || sort == "status" { // sort
		tx = tx.Order(sort)
	}
	task := []Task{}
	tx = tx.Find(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return task, nil
}

func (r taskRepositoryDB) GetImage(uid *uuid.UUID) (*Task, error) {
	task := Task{}
	tx := r.db.Where("uid = ?", uid).Preload(clause.Associations).Find(&task)
	if tx.Error != nil {
		logrus.Errorln("select task error: ", tx.Error)
		return nil, tx.Error
	}
	return &task, nil
}

// Create Task
func (r taskRepositoryDB) CreateTask(task Task) (*Task, error) {
	tx := r.db.Create(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &task, nil
}

func (r taskRepositoryDB) CreateTaskImage(taskimg TaskImage) (*TaskImage, error) {
	tx := r.db.Create(&taskimg)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &taskimg, nil
}

// Update Task
func (r taskRepositoryDB) UpdateTask(uid *uuid.UUID, task Task) error {
	tx := r.db.Model(&Task{}).Where("uid = ?", uid).Updates(task)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r taskRepositoryDB) UpdateTaskImage(uid *uuid.UUID, taskImg TaskImage) error {
	tx := r.db.Model(&TaskImage{}).Where("uid = ?", uid).Updates(taskImg)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
