package website

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) Migrate(ctx context.Context) error {
	m := &Website{}
	return r.db.WithContext(ctx).AutoMigrate(&m)
}

func (r *GormRepository) Create(ctx context.Context, website Website) (*Website, error) {
	if err := r.db.WithContext(ctx).Create(&website).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &website, nil
}

func (r *GormRepository) All(ctx context.Context) ([]Website, error) {
	var websites []Website

	if err := r.db.WithContext(ctx).Find(&websites).Error; err != nil {
		return nil, err
	}

	return websites, nil
}

func (r *GormRepository) GetByName(ctx context.Context, name string) (*Website, error) {
	var website Website
	if err := r.db.WithContext(ctx).Where("name = ?", name).Find(&website).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	return &website, nil
}

func (r *GormRepository) Update(ctx context.Context, id int64, updatedWebsite Website) (*Website, error) {
	result := r.db.WithContext(ctx).Where("id = ?", id).Save(&updatedWebsite)
	if err := result.Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		return nil, ErrUpdatedFailed
	}

	return &updatedWebsite, nil
}

func (r *GormRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Website{}, id)
	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ErrDeletedFailed
	}

	return nil
}
