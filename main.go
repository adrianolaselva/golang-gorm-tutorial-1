package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
	"time"
)

type Workplace struct {
	gorm.Model
	Name string `gorm:"size:50;not null"`
	Address string `gorm:"size:255;not null"`
	Phone *string `gorm:"size:20"`
	Workers []*Worker
}

type Worker struct {
	gorm.Model
	WorkplaceID uint `gorm:"not null"`
	Workplace *Workplace
	Name string `gorm:"size:61;not null"`
	Birthday time.Time `gorm:"type:datetime"`
	Phone *string `gorm:"size:20"`
	Recipes []*Recipe `gorm:"many2many:worker_recipes"`
}

type Recipe struct {
	gorm.Model
	Name string `gorm:"size:50"`
	Workers []*Worker `gorm:"many2many:worker_recipes"`
	Toppings []*Topping `gorm:"many2many:recipe_toppings"`
}

type Topping struct {
	gorm.Model
	Name string `gorm:"size:20; not null"`
	Pizzas []*Pizza `gorm:"many2many:recipe_toppings"`
}

type Size struct {
	gorm.Model
	Name string `gorm:"size:20;not null"`
}

type Pizza struct {
	gorm.Model
	RecipeID uint `gorm:"not null;unique_index:pizzas"`
	Recipe Recipe
	SizeID uint `gorm:"not null;unique_index:pizzas"`
	Size Size
	Price decimal.Decimal `gorm:"not null;type:decimal(10,2)"`
}

func Connect () (*gorm.DB, error) {
	return gorm.Open("mysql", "root:root@tcp(127.0.0.1:13306)/pizza?charset=utf8mb4&parseTime=true&loc=Local")
}

func BoxString(x string) *string {
	return &x
}

func Migrate(db *gorm.DB) {
	workplacePrototype := &Workplace{}
	workerPrototype := &Worker{}
	recipePrototype := &Recipe{}
	sizePrototype := &Size{}
	pizzaPrototype := &Pizza{}
	toppingsPrototype := Topping{}

	db.AutoMigrate(workplacePrototype, workerPrototype, recipePrototype, sizePrototype, pizzaPrototype, toppingsPrototype)

	db.Model(workerPrototype).AddForeignKey("workplace_id", "workplaces(id)", "RESTRICT", "CASCADE")
	db.Table("worker_recipes").AddForeignKey("worker_id", "workers(id)", "RESTRICT", "CASCADE")
	db.Table("worker_recipes").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "CASCADE")
	db.Table("recipe_toppings").AddForeignKey("recipe_id", "recipes(id)", "RESTRICT", "CASCADE")
	db.Table("recipe_toppings").AddForeignKey("topping_id", "toppings(id)", "RESTRICT", "CASCADE")
}

func Seed(db *gorm.DB) {
	cheese := &Topping{
		Name:   "Cheese",
	}
	tomatoeSauce := &Topping{
		Name: "Tomatoe Sauce",
	}
	onions := &Topping{
		Name: "Onions",
	}
	tomatoeSlices := &Topping{
		Name: "Tomatoe Slices",
	}
	hamSlices := &Topping{
		Name: "Ham Slices",
	}
	pepperoni := &Topping{
		Name: "Pepperoni",
	}

	recipe1 := &Recipe{
		Name: "Mozzarella",
		Toppings: []*Topping{
			tomatoeSauce,
			cheese,
		},
	}
	recipe2 := &Recipe{
		Name: "Onions",
		Toppings: []*Topping{
			onions,
			cheese,
		},
	}
	recipe3 := &Recipe{
		Name: "Napolitan",
		Toppings: []*Topping{
			tomatoeSauce,
			tomatoeSlices,
			cheese,
			hamSlices,
		},
	}
	recipe4 := &Recipe{
		Name: "Pepperoni",
		Toppings: []*Topping{
			tomatoeSauce,
			pepperoni,
			cheese,
		},
	}

	db.Save(recipe1)
	db.Save(recipe2)

	sizePersonal := Size{
		Name:  "Personal",
	}
	sizeSmall := Size{
		Name: "Small",
	}
	sizeMedium := Size{
		Name: "Medium",
	}
	sizeBig := Size{
		Name: "Big",
	}
	sizeExtraBig := Size{
		Name: "Extra Big",
	}

	db.Save(&sizePersonal)
	db.Save(&sizeSmall)
	db.Save(&sizeMedium)
	db.Save(&sizeBig)
	db.Save(&sizeExtraBig)

	workplace1 := &Workplace{
		Name: "Workplace One",
		Address: "Fake st. 123rd",
	}
	workplace2 := &Workplace{
		Name: "Workplace Two",
		Address: "Evergreen Terrace 742nd",
		Phone: BoxString("(11) 986066232"),
		Workers: []*Worker{
			{
				Name: "Adriano Moreira",
				Birthday: time.Date(1987, 2, 11, 0, 0, 0, 0, time.UTC),
				Recipes: []*Recipe{
					recipe1,
					recipe2,
					recipe3,
				},
			},
			{
				Name: "Ana paula",
				Birthday: time.Date(1987, 6, 25, 0, 0, 0, 0, time.UTC),
				Recipes: []*Recipe{
					recipe1,
					recipe2,
					recipe4,
				},
			},
		},
	}

	for sizeIndex, size := range []*Size{&sizePersonal, &sizeSmall, &sizeMedium, &sizeBig, &sizeExtraBig} {
		for recipeIndex, recipe := range []*Recipe{recipe1,recipe2,recipe3,recipe4} {
			db.Save(&Pizza{
				Recipe:   *recipe,
				Size:     *size,
				Price:    decimal.NewFromFloat((float64(0.1) * float64(recipeIndex + 1) + 1) + float64(sizeIndex) * 5),
			})
		}
	}
	db.Save(workplace1)
	db.Save(workplace2)

	fmt.Printf("Workplace created:\n%v\n%v\n", workplace1, workplace2)
	fmt.Printf("Recipes created:\n%v\n", []*Recipe{recipe1, recipe2,recipe3,recipe4})
}

func ListEvetything(db *gorm.DB) {
	var workplaces []Workplace
	t := db.Preload("Workers")
	t = t.Preload("Workers.Recipes")
	t = t.Preload("Workers.Recipes.Toppings")
	t.Find(&workplaces)

	for _, workplace := range workplaces {
		fmt.Printf("Workplace data: %v\n", workplace)

		for _, worker := range workplace.Workers {
			fmt.Printf("Worker data: %v\n", worker)

			for _, recipe := range worker.Recipes {
				fmt.Printf("Recipe data: %v\n", recipe)

				for _, topping := range recipe.Toppings {
					fmt.Printf("Topping data: %v\n", topping)
				}
			}
		}
	}
}

func ClearEverything(db *gorm.DB) {
	err := db.Delete(&Workplace{}).Error
	fmt.Println("Deleting all workplace records", err)
}

func main() {
	if db, err := Connect(); err != nil {
		fmt.Println("Failed to connect database: ", err.Error())
	} else {
		db.LogMode(true)
		defer db.Close()

		Migrate(db)
		Seed(db)
		ListEvetything(db)
		ClearEverything(db)
	}
}