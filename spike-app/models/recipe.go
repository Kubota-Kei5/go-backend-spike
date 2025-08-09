package models

type Recipe struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"not null"`
	Servings      int    `gorm:"not null"`
	CookingTime   int    `gorm:"not null"`
}

type Ingredients struct {
	ID           uint    `gorm:"primaryKey"`
	Name         string  `gorm:"not null"`
	Enerc_kcal   float64 `gorm:"not null"`
	Biot         float64 `gorm:"not null"`
	Ca 	         float64 `gorm:"not null"`
	Chocdf       float64 `gorm:"not null"`
	Chole        float64 `gorm:"not null"`
	Cu           float64 `gorm:"not null"`
	Fat          float64 `gorm:"not null"`
	Fe           float64 `gorm:"not null"`
	Fib          float64 `gorm:"not null"`
	Fol          float64 `gorm:"not null"`
	K            float64 `gorm:"not null"`
	Mg           float64 `gorm:"not null"`
	Mn           float64 `gorm:"not null"`
	Na           float64 `gorm:"not null"`
	Ncal_eq	     float64 `gorm:"not null"`
	Nia          float64 `gorm:"not null"`
	P            float64 `gorm:"not null"`
	Pantac       float64 `gorm:"not null"`
	Protein      float64 `gorm:"not null"`
	Ribf         float64 `gorm:"not null"`
	Thia         float64 `gorm:"not null"`
	Vitb12       float64 `gorm:"not null"`
	Vitb6        float64 `gorm:"not null"`
	Vitc         float64 `gorm:"not null"`
	Vitba        float64 `gorm:"not null"`
	Vitd         float64 `gorm:"not null"`
	Vite         float64 `gorm:"not null"`
	Vitk         float64 `gorm:"not null"`
	Zn           float64 `gorm:"not null"`
}

type RecipeIngredient struct {
	RecipeID     uint    `gorm:"not null"`
	IngredientID uint    `gorm:"not null"`
	Quantity     float64 `gorm:"not null"`

	Recipe		Recipe   `gorm:"foreignKey:RecipeID";references:ID;constraint:OnDelete:CASCADE"`
	Ingredient	Ingredients `gorm:"foreignKey:IngredientID";references:ID;constraint:OnDelete:RESTRICT"`
}