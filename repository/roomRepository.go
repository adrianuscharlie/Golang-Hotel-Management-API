package repository

import (
	"hms-backend/model"
	"time"

	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

type RoomRepository interface {
	FindAll() ([]model.Room, error)
	FindByID(id int) (*model.Room, error)
	Create(room *model.Room) error
	Update(room *model.Room) error
	Delete(id int) error
	FindByNumber(number string) (*model.Room, error)
	FindAvailable(checkIn, checkOut time.Time) ([]model.Room, error)
	ChangeStatus(id int, status string) error
	CreateRoomType(roomType *model.RoomType) error
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db}
}

func (r *roomRepository) FindAll() ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.Preload("RoomType").Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) FindByID(id int) (*model.Room, error) {
	var room model.Room
	err := r.db.Preload("RoomType").Where("id = ?", id).First(&room).Error
	return &room, err
}

func (r *roomRepository) Create(m *model.Room) error {
	return r.db.Create(m).Error
}

func (r *roomRepository) Update(m *model.Room) error {
	return r.db.Save(m).Error
}

func (r *roomRepository) Delete(id int) error {
	return r.db.Delete(&model.Room{}, id).Error
}

func (r *roomRepository) FindByNumber(s string) (*model.Room, error) {
	var room model.Room
	err := r.db.Preload("RoomType").Where("number = ?").First(&room).Error
	return &room, err
}

func (r *roomRepository) FindAvailable(checkIn, checkOut time.Time) ([]model.Room, error) {
	var rooms []model.Room

	err := r.db.Raw(`
		SELECT * FROM rooms 
		WHERE id NOT IN (
			SELECT room_id FROM bookings
			WHERE NOT (
				check_out <= ? OR check_in >= ?
			)
		)
	`, checkIn, checkOut).Scan(&rooms).Error
	for i := range rooms {
		r.db.Model(&rooms[i]).Association("RoomType").Find(&rooms[i].RoomType)
	}
	return rooms, err
}

func (r *roomRepository) ChangeStatus(id int, status string) error {
	room, err := r.FindByID(id)
	if err != nil {
		return err
	}
	room.Status = status
	return r.db.Save(&room).Error
}

func (r *roomRepository) CreateRoomType(roomType *model.RoomType) error {
	return r.db.Create(&roomType).Error
}
