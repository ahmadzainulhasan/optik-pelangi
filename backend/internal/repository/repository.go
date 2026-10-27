package repository

import "gorm.io/gorm"

type Repository[T any] struct {
}

func (repository *Repository[T]) Save(tx *gorm.DB, entity T) error {
	return tx.Save(&entity).Error
}

func (repository *Repository[T]) Delete(tx *gorm.DB, entity T) error {
	return tx.Delete(&entity).Error
}

func (repository *Repository[T]) CountById(tx *gorm.DB, id any) (int64, error) {
	var total int64
	err := tx.Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (repository *Repository[T]) FindByID(tx *gorm.DB, id any) (T, error) {
	var entity T
	err := tx.Where("id = ?", id).Take(&entity).Error
	if err != nil {
		return entity, err
	}
	return entity, nil
}
	