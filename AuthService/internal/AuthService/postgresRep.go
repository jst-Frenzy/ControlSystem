package AuthService

import (
	"gorm.io/gorm"
)

//go:generate mockgen -source=postgresRep.go -destination=mocks/mockPostgres.go

type AuthPostgresRepo interface {
	GetUser(string) (User, error)
	CreateUser(User) (User, error)
	SaveRefreshToken(RefreshToken) error
	GetUserByRefreshToken(string) (User, error)

	ChangeRole(int, string) (User, error)
}

type authPostgresRepo struct {
	db *gorm.DB
}

func NewAuthPostgresRepo(db *gorm.DB) AuthPostgresRepo {
	return &authPostgresRepo{db: db}
}

func (r *authPostgresRepo) GetUser(email string) (User, error) {
	var err error
	var foundUser User
	if err = r.db.Table("users").Where("email = ?", email).First(&foundUser).Error; err != nil {
		return User{}, err
	}
	return foundUser, nil
}

func (r *authPostgresRepo) CreateUser(u User) (User, error) {
	var err error
	if err = r.db.Table("users").Create(&u).Error; err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *authPostgresRepo) SaveRefreshToken(token RefreshToken) error {
	return r.db.Table("refresh_tokens").Create(&token).Error
}

func (r *authPostgresRepo) GetUserByRefreshToken(refreshTokenHash string) (User, error) {
	var user User
	var err error
	if err = r.db.Table("users").
		Where("id = (select user_id from refresh_tokens where token_hash = ?)", refreshTokenHash).
		First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *authPostgresRepo) ChangeRole(id int, newRole string) (User, error) {
	var u User
	if err := r.db.Table("users").Where("id = ?", id).Update("role", newRole).Error; err != nil {
		return User{}, err
	}
	if err := r.db.Table("users").Where("id = ?", id).First(&u).Error; err != nil {
		return User{}, err
	}
	return u, nil
}
