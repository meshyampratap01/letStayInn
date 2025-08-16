package bookingRepository

import (
	"github.com/meshyampratap01/letStayInn/internal/models"
	"gorm.io/gorm"
)

type GormBookingRepository struct {
	db *gorm.DB
}

func NewGormBookingRepository(db *gorm.DB) BookingRepository {
	return &GormBookingRepository{db: db}
}

func (r *GormBookingRepository) GetAllBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Find(&bookings).Error
	return bookings, err
}

func (r *GormBookingRepository) SaveBookings(bookings []models.Booking) error {
	for _, booking := range bookings {
		if err := r.db.Save(&booking).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *GormBookingRepository) GetBookingsByUserID(userID string) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *GormBookingRepository) UpdateBooking(updated models.Booking) error {
	return r.db.Save(&updated).Error
}

func (r *GormBookingRepository) GetBookingByID(bookingID string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.First(&booking, "id = ?", bookingID).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *GormBookingRepository) GetActiveBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("status != ?", models.BookingStatusCancelled).Find(&bookings).Error
	return bookings, err
}

func (r *GormBookingRepository) CheckRoomBooked(roomNumber int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).Where("room_num = ? AND status = ?", roomNumber, models.BookingStatusBooked).Count(&count).Error
	return count > 0, err
}
