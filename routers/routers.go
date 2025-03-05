package routers

import (
	"database/sql"
	"ecommerce-project-go/auth"
	"ecommerce-project-go/controllers"
	"ecommerce-project-go/repository"
	"ecommerce-project-go/service"

	"github.com/gin-gonic/gin"
)

func StartServer(root string, db *sql.DB) {

	r := gin.Default()
	users := r.Group("/users")
	category := r.Group("/category")
	Product := r.Group("/product")
	transaction := users.Group("/transaction")

	// user endpoint handler
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := controllers.NewUserHandler(userService)
	// all user
	users.POST("/register", userHandler.RegisterUser)
	users.POST("/login", userHandler.Login)
	users.PUT("/edit", auth.MiddlewareUserAuth(userService), userHandler.UpdateUser)
	users.DELETE("/delete", auth.MiddlewareUserAuth(userService), userHandler.DeleteUser)

	// admin
	users.GET("/get-all-users", auth.MiddlewareUserAuth(userService), userHandler.GetAllUsers)

	// categories endpoint handler
	categoryRepository := repository.NewCategoriesRepository(db)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := controllers.NewCatHandler(categoryService)

	// all user
	category.GET("/get-all", categoryHandler.GetAllCategories)
	category.GET("/:category_id/items", categoryHandler.GetAllProductsByCatId)

	// admin
	category.POST("/add", auth.MiddlewareUserAuth(userService), categoryHandler.InsertCategory)
	category.PUT("/edit/:category_id", auth.MiddlewareUserAuth(userService), categoryHandler.UpdateCategory)
	category.DELETE("/delete/:category_id", auth.MiddlewareUserAuth(userService), categoryHandler.DeleteCategories)

	// Product endpoint handler
	ProductRepository := repository.NewProductRepository(db)
	ProductService := service.NewProductService(ProductRepository)
	ProductHandler := controllers.NewProductHandler(ProductService)
	// all user
	Product.GET("/get-all", ProductHandler.GetAll)
	Product.GET("/get/:product_id", ProductHandler.GetById)

	// admin
	Product.POST("/add", auth.MiddlewareUserAuth(userService), ProductHandler.InsertProduct)
	Product.PUT("/edit/:product_id", auth.MiddlewareUserAuth(userService), ProductHandler.UpdateProduct)
	Product.DELETE("/delete/:product_id", auth.MiddlewareUserAuth(userService), ProductHandler.DeleteProduct)

	// transactions endpoint handler
	transactionRepository := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepository)
	transactionHandler := controllers.NewTransactionHandler(transactionService)

	// user
	transaction.GET("/get-all", auth.MiddlewareUserAuth(userService), transactionHandler.GetAll)

	// example /get?status=Paid
	transaction.GET("/get", auth.MiddlewareUserAuth(userService), transactionHandler.GetByStatus)
	transaction.POST("/create", auth.MiddlewareUserAuth(userService), transactionHandler.CreateTransaction)

	// example /1?action=pay or /1?action=cancel
	transaction.PUT("/:trans_id", auth.MiddlewareUserAuth(userService), transactionHandler.UpdateTransaction)
	transaction.PUT("/admin/:trans_id", auth.MiddlewareUserAuth(userService), transactionHandler.UpdateAdmin)
	// admin
	transaction.GET("/admin/get-all", auth.MiddlewareUserAuth(userService), transactionHandler.GetAllAdmin)
	transaction.GET("/admin/get", auth.MiddlewareUserAuth(userService), transactionHandler.GetByStatusAdmin)

	r.Run(root)
}
