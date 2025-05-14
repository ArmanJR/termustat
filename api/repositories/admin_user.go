package repositories

import (
	"context"
	"github.com/armanjr/termustat/api/dto"
	"github.com/armanjr/termustat/api/errors"
	"github.com/armanjr/termustat/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminUserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByStudentID(ctx context.Context, studentID string) (*models.User, error)
	FindByEmailOrStudentID(ctx context.Context, email, studentID string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, pagination *dto.PaginationQuery) (*dto.PaginatedList[models.User], error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
	UpdateEmailVerification(ctx context.Context, userID uuid.UUID, verified bool) error
	FindByUniversity(ctx context.Context, universityID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[models.User], error)
	FindByFaculty(ctx context.Context, facultyID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[models.User], error)
}

type adminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	var created models.User
	if err := r.db.WithContext(ctx).First(&created, user.ID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch created user")
	}

	return &created, nil
}

func (r *adminUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("user", id.String())
		}
		return nil, errors.Wrap(err, "database error")
	}
	return &user, nil
}

func (r *adminUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("user", "email: "+email)
		}
		return nil, errors.Wrap(err, "database error")
	}
	return &user, nil
}

func (r *adminUserRepository) FindByStudentID(ctx context.Context, studentID string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("student_id = ?", studentID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("user", "student_id: "+studentID)
		}
		return nil, errors.Wrap(err, "database error")
	}
	return &user, nil
}

func (r *adminUserRepository) FindByEmailOrStudentID(ctx context.Context, email, studentID string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ? OR student_id = ?", email, studentID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewNotFoundError("user", "email or student_id")
		}
		return nil, errors.Wrap(err, "database error")
	}
	return &user, nil
}

func (r *adminUserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	var updated models.User
	if err := r.db.WithContext(ctx).First(&updated, user.ID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch updated user")
	}

	return &updated, nil
}

func (r *adminUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete user")
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("user", id.String())
	}
	return nil
}

func (r *adminUserRepository) GetAll(ctx context.Context, pagination *dto.PaginationQuery) (*dto.PaginatedList[models.User], error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count users")
	}

	if err := query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch users")
	}

	return &dto.PaginatedList[models.User]{
		Items: users,
		Total: total,
		Page:  pagination.Page,
		Limit: pagination.Limit,
	}, nil
}

func (r *adminUserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", hashedPassword)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update password")
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("user", userID.String())
	}
	return nil
}

func (r *adminUserRepository) UpdateEmailVerification(ctx context.Context, userID uuid.UUID, verified bool) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("email_verified", verified)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to update email verification")
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("user", userID.String())
	}
	return nil
}

func (r *adminUserRepository) FindByUniversity(ctx context.Context, universityID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[models.User], error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{}).Where("university_id = ?", universityID)

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count users")
	}

	if err := query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch users")
	}

	return &dto.PaginatedList[models.User]{
		Items: users,
		Total: total,
		Page:  pagination.Page,
		Limit: pagination.Limit,
	}, nil
}

func (r *adminUserRepository) FindByFaculty(ctx context.Context, facultyID uuid.UUID, pagination *dto.PaginationQuery) (*dto.PaginatedList[models.User], error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{}).Where("faculty_id = ?", facultyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count users")
	}

	if err := query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch users")
	}

	return &dto.PaginatedList[models.User]{
		Items: users,
		Total: total,
		Page:  pagination.Page,
		Limit: pagination.Limit,
	}, nil
}
