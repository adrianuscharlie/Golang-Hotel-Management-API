package model

type Room struct {
	ID         uint `gorm:"primaryKey"`
	RoomTypeID uint
	RoomType   RoomType `gorm:"foreignKey:RoomTypeID;references:ID"`
	Status     string
	Number     string `gorm:"unique"`
}

type RoomType struct {
	ID          uint `gorm:"primaryKey"`
	Price       float64
	Capacity    int
	Description string
	Name        string
}
