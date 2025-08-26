package roomRepository

import (
	"fmt"

	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/gorm"
)

type GormRoomRepository struct {
	db *gorm.DB
}

func NewGormRoomRepository(db *gorm.DB) IRoomRepository {
	return &GormRoomRepository{db: db}
}

func (r *GormRoomRepository) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Find(&rooms).Error
	return rooms, err
}

func (r *GormRoomRepository) SaveRooms(rooms []models.Room) error {
	for _, room := range rooms {
		if err := r.db.Save(&room).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *GormRoomRepository) DeleteRoomByNumber(number int) error {
	return r.db.Where("number = ?", number).Delete(&models.Room{}).Error
}

func (r *GormRoomRepository) GetAvailableRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Where("is_available = ?", true).Find(&rooms).Error
	return rooms, err
}

func (r *GormRoomRepository) AddRoom(room models.Room) error {
	return r.db.Create(&room).Error
}

func (r *GormRoomRepository) GetRoomNumberByBookingID(bookingID string) (string, error) {
	var booking models.Booking
	err := r.db.First(&booking, "id = ?", bookingID).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", booking.RoomNum), nil
}

func (r *GormRoomRepository) RoomExists(number int) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Room{}).
		Where("number = ?", number).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
