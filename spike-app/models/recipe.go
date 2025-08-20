package models

import "gorm.io/gorm"

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

type Recipe struct {
	ID          uint   `gorm:"primaryKey" json:"ID" form:"id"`
	Title       string `gorm:"not null" json:"Title" form:"title"`
	Servings    int    `gorm:"not null" json:"Servings" form:"servings"`
	CookingTime int    `gorm:"not null" json:"CookingTime" form:"cooking_time"`
}

type Ingredients struct {
	ID         uint    `gorm:"primaryKey"`
	Name       string  `gorm:"not null"`
	Enerc_kcal float64 `gorm:"not null"`
	Biot       float64 `gorm:"not null"`
	Ca         float64 `gorm:"not null"`
	Chocdf     float64 `gorm:"not null"`
	Chole      float64 `gorm:"not null"`
	Cu         float64 `gorm:"not null"`
	Fat        float64 `gorm:"not null"`
	Fe         float64 `gorm:"not null"`
	Fib        float64 `gorm:"not null"`
	Fol        float64 `gorm:"not null"`
	K          float64 `gorm:"not null"`
	Mg         float64 `gorm:"not null"`
	Mn         float64 `gorm:"not null"`
	Na         float64 `gorm:"not null"`
	Ncal_eq    float64 `gorm:"not null"`
	Nia        float64 `gorm:"not null"`
	P          float64 `gorm:"not null"`
	Pantac     float64 `gorm:"not null"`
	Protein    float64 `gorm:"not null"`
	Ribf       float64 `gorm:"not null"`
	Thia       float64 `gorm:"not null"`
	Vitb12     float64 `gorm:"not null"`
	Vitb6      float64 `gorm:"not null"`
	Vitc       float64 `gorm:"not null"`
	Vitba      float64 `gorm:"not null"`
	Vitd       float64 `gorm:"not null"`
	Vite       float64 `gorm:"not null"`
	Vitk       float64 `gorm:"not null"`
	Zn         float64 `gorm:"not null"`
}

type RecipeIngredient struct {
	RecipeID     uint    `gorm:"not null"`
	IngredientID uint    `gorm:"not null"`
	Quantity     float64 `gorm:"not null"`

	Recipe     Recipe      `gorm:"foreignKey:RecipeID;references:ID;constraint:OnDelete:CASCADE"`
	Ingredient Ingredients `gorm:"foreignKey:IngredientID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (r *Recipe) Create() (*Recipe, error) {
	if err := db.Create(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Recipe) Save() (*Recipe, error) {
	if err := db.Save(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Recipe) Delete() error {
	if err := db.Where("id = ?", r.ID).Delete(&r).Error; err != nil {
		return err
	}
	return nil
}

func GetAllRecipes() ([]Recipe, error) {
	var recipes []Recipe
	if err := db.Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func GetRecipeByID(id string, recipe *Recipe) error {
	if err := db.First(&recipe, id).Error; err != nil {
		return err
	}
	return nil
}
