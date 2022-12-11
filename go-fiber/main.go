package main

import (
	"net/http"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)


type customer struct {
	CustName    string `json:"customer"`
	PhoneNumber string `json:"phone"`
	Date        string `json:"date"`
	State       string `json:"state"`
	DeviceName  string `json:"device"`
	Problem     string `json:"problem"`
	Mail        string `json:"mail"`
	Price       string `json:"price"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) DeleteCustomer(context *fiber.Ctx) error {
	CustModel := models.Customer{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(CustModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete Customer",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Customer delete successfully",
	})
	return nil
}

func (r *Repository) CreateCustomer(context *fiber.ctx) error {
	Customer := customer{}
	// burda json encode-decode işlemi yapılıyor
	context.BodyParser(&Customer)
	err := r.DB.Create(&Customer).Error

	if err != nil {
		context.Status(http.StatusBadRequest).json(
			&fiber.Map{"message": "could not create customer"})
		return err
	}
	context.Status(http.StatusOK).json(&fiber.Map{
		"message": "customer added."})
	return nil
}

func (r *Repository) GetCustomers(context *fiber.Ctx) error {
	CustModel := &[]models.Customer{}

	err := r.DB.Find(CustModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get Customers"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Customers fetched successfully",
		"data":    CustModel,
	})
	return nil
}

// create customer update funtion
func (r *Repository) UpdateCustomer(context *fiber.Ctx) error {
	CustModel := models.Customer{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&CustModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the customer"})
		return err
	}
}

func (r *Repository) GetCustomerbyName(context *fiber.Ctx) error {

	id := context.Params("id")
	CustModel := &models.Customer{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(CustModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the customer"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Customer id fetched successfully",
		"data":    CustModel,
	})
	return nil
}

// func (r *Repository) GetCustomerbyName(context *fiber.ctx) error {
// 	Customer := customer{}
// 	err := r.DB.Where("cust_name = ?", context.Params("id")).First(&Customer).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).json(
// 			&fiber.Map{"message": "could not get customer"})
// 		return err
// 	}
// 	context.Status(http.StatusOK).json(&fiber.Map{
// 		"message": "customer found."})
// 	return nil
// }

func (r *Repository) SetupRoutes() {
	api := technic_app.Group("/app")
	api.Post("/createcustomerinfo", r.CreateCustomer)
	api.Delete("/deletetask/:id", r.DeleteCustomer)
	api.Get("/getcustomerinfo/:id", r.GetCustomerbyName)
	api.Get("/customers", r.GetCustomers)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateCustomers(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
