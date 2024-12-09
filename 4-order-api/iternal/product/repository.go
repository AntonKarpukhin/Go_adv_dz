package product

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"orderApi/pkg/db"
)

type ProductRepository struct {
	DataBase *db.Db
}

func NewProductRepository(dataBase *db.Db) *ProductRepository {
	return &ProductRepository{
		DataBase: dataBase,
	}
}

func (repo *ProductRepository) Create(product *Product) (*Product, error) {
	result := repo.DataBase.DB.Create(product)

	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) GetByName(name string) (*Product, error) {
	var product Product
	result := repo.DataBase.DB.Where("name = ?", name).First(&product)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &product, nil
}

func (repo *ProductRepository) Update(product *Product) (*Product, error) {
	result := repo.DataBase.DB.Clauses(clause.Returning{}).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) Delete(id uint) error {
	result := repo.DataBase.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ProductRepository) GetById(id uint) (*Product, error) {
	var product Product
	result := repo.DataBase.DB.First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}
