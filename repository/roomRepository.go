package repository

import (
	"hms-backend/model"
	"hms-backend/request"

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
	FindAvailable(params request.RoomFilterParams) ([]model.Room, error)
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
	err := r.db.Preload("RoomType").Where("number = ?", s).First(&room).Error
	return &room, err
}

func (r *roomRepository) FindAvailable(params request.RoomFilterParams) ([]model.Room, error) {
	var rooms []model.Room

	// 1. Start query chain. Preload loads the RoomType data efficiently after the query is done.
	query := r.db.Model(&model.Room{}).Preload("RoomType")

	// 2. JOIN with room_types table to allow filtering on price and category name.
	//    We use the actual table name `room_types` for clarity.
	query = query.Joins("JOIN room_types ON room_types.id = rooms.room_type_id")

	// 3. Create the subquery to find all IDs of rooms that are UNAVAILABLE.
	//    This date logic correctly finds ALL overlapping bookings.
	subquery := r.db.Table("bookings").Select("room_id").
		Where("NOT (check_out_date <= ? OR check_in_date >= ?)", params.CheckIn, params.CheckOut)

	// 4. Filter out the unavailable rooms from the main query.
	query = query.Where("rooms.id NOT IN (?)", subquery)

	// --- Dynamically add the rest of the user's filters ---

	// 5. Filter by category (RoomType name) if provided.
	if params.Category != "" {
		query = query.Where("room_types.name = ?", params.Category)
	}

	// 6. Filter by minimum price if provided.
	if params.MinPrice > 0 {
		query = query.Where("room_types.price >= ?", params.MinPrice)
	}

	// 7. Filter by maximum price if provided.
	//    CORRECTED LOGIC: The filter should apply if MaxPrice is a positive number.
	if params.MaxPrice > 0 {
		query = query.Where("room_types.price <= ?", params.MaxPrice)
	}

	// Execute the fully constructed query
	err := query.Find(&rooms).Error
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
