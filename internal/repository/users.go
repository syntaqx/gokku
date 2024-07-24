package repository

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid; primarykey; default:uuid_generate_v4()"`
	Login     string    `gorm:"uniqueIndex; not null"`
	Email     string    `gorm:"uniqueIndex; not null"`
	Password  string    `gorm:"not null"`
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	u.ID = id

	if u.AvatarURL == "" {
		u.AvatarURL = "https://avatars.githubusercontent.com/u/6037730?v=4"
	}

	return
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type UsersRepository struct {
	db *gorm.DB
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login credentials")
)

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *UsersRepository) GetUserByLogin(login string) (*User, error) {
	var user User
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *UsersRepository) GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *UsersRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *UsersRepository) ValidateUser(identifier, password string) (*User, error) {
	user, err := r.GetUserByEmail(identifier)
	if err != nil {
		user, err = r.GetUserByLogin(identifier)
		if err != nil {
			return nil, ErrInvalidLogin
		}
	}

	if !user.CheckPassword(password) {
		return nil, ErrInvalidLogin
	}

	return user, nil
}

func (r *UsersRepository) Seed() {
	initialEmail := "admin@example.com"
	_, err := r.GetUserByEmail(initialEmail)
	if err == nil {
		// User already exists
		return
	}

	user := &User{
		Login: "admin",
		Email: initialEmail,
	}
	user.SetPassword("password") // In a real application, use a more secure password

	_ = r.CreateUser(user) // Ignoring the error for simplicity
}
