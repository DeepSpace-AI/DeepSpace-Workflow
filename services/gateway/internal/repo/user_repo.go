package repo

import (
	"context"
	"errors"
	"time"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}

func (r *UserRepo) List(ctx context.Context, page, pageSize int, search, role, status string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.WithContext(ctx).Model(&model.User{})

	if search != "" {
		query = query.Where("email LIKE ? OR id IN (SELECT user_id FROM user_profiles WHERE display_name LIKE ? OR full_name LIKE ?)", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("id desc").Preload("Profile").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, id int64, loginTime time.Time) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("last_login_at", loginTime).Error
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}
