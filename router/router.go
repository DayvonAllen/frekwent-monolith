package router

import (
	"freq/handlers"
	"freq/middleware"
	"freq/repository"
	"freq/services"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	ah := handlers.AuthHandler{AuthService: services.NewAuthService(repository.NewAuthRepoImpl())}
	ch := handlers.CouponHandler{CouponService: services.NewCouponService(repository.NewCouponRepoImpl())}
	ih := handlers.LoginIpHandler{LoginIpService: services.NewLoginIpService(repository.NewLoginIpRepoImpl())}
	crh := handlers.CustomerHandler{CustomerService: services.NewCustomerService(repository.NewCustomerRepoImpl())}
	ph := handlers.PurchaseHandler{PurchaseService: services.NewPurchaseService(repository.NewPurchaseRepoImpl())}
	prh := handlers.ProductHandler{ProductService: services.NewProductService(repository.NewProductRepoImpl())}
	eh := handlers.EmailHandler{EmailService: services.NewEmailService(repository.NewEmailRepoImpl())}
	mmh := handlers.MailMemberHandler{MailMemberService: services.NewMailMemberService(repository.NewMailMemberRepoImpl())}

	prometheus := fiberprometheus.New("freq")
	prometheus.RegisterAt(app, "/metrics")

	app.Use(prometheus.Middleware)
	app.Use(recover.New())
	app.Use(pprof.New())

	api := app.Group("", logger.New())

	auth := api.Group("/iriguchi/auth")
	auth.Post("/login", ah.Login)
	auth.Get("/logout", middleware.IsLoggedIn, ah.Logout)

	subscribe := api.Group("/email")
	subscribe.Post("/subscribe", mmh.Create)

	product := api.Group("/products")
	product.Post("/buy", ph.Purchase)
	product.Post("/ids", prh.FindAllByProductIds)
	product.Get("/:id", prh.FindByProductId)
	product.Get("/category/:category", prh.FindAllByCategory)
	product.Get("", prh.FindAll)

	optOut := api.Group("/list/optout")
	optOut.Get("/:email", crh.UpdateOptInStatus)

	purchase := api.Group("/iriguchi/purchases")
	purchase.Get("/transactions/:state", middleware.IsLoggedIn, ph.CalculateTransactionsByState)
	purchase.Get("/:id", middleware.IsLoggedIn, ph.FindByPurchaseById)
	purchase.Put("/shipped/:id", middleware.IsLoggedIn, ph.UpdateShippedStatus)
	purchase.Put("/delivered/:id", middleware.IsLoggedIn, ph.UpdateDeliveredStatus)
	purchase.Put("/address/:id", middleware.IsLoggedIn, ph.UpdatePurchaseAddress)
	purchase.Put("/tracking/:id", middleware.IsLoggedIn, ph.UpdateTrackingNumber)
	purchase.Get("", middleware.IsLoggedIn, ph.FindAll)

	items := api.Group("/iriguchi/items")
	items.Put("/name/:id", middleware.IsLoggedIn, prh.UpdateName)
	items.Get("/name", middleware.IsLoggedIn, prh.FindByProductName)
	items.Put("/description/:id", middleware.IsLoggedIn, prh.UpdateDescription)
	items.Put("/quantity/:id", middleware.IsLoggedIn, prh.UpdateQuantity)
	items.Put("/ingredients/:id", middleware.IsLoggedIn, prh.UpdateIngredients)
	items.Put("/price/:id", middleware.IsLoggedIn, prh.UpdatePrice)
	items.Put("/category/:id", middleware.IsLoggedIn, prh.UpdateCategory)
	items.Delete("/delete/:id", middleware.IsLoggedIn, prh.DeleteById)
	items.Post("", middleware.IsLoggedIn, prh.Create)

	email := api.Group("/iriguchi/email")
	email.Get("/get/:email", middleware.IsLoggedIn, eh.FindAllByEmail)
	email.Get("/status/:status", middleware.IsLoggedIn, eh.FindAllByStatus)
	email.Get("/members", middleware.IsLoggedIn, mmh.FindAll)
	email.Delete("/members/:id", middleware.IsLoggedIn, mmh.DeleteById)
	email.Get("", middleware.IsLoggedIn, eh.FindAll)

	coupon := api.Group("/iriguchi/coupon")
	coupon.Post("/send/:emailType", middleware.IsLoggedIn, eh.SendEmail)
	coupon.Post("/mass/send", middleware.IsLoggedIn, eh.MassCouponEmail)
	coupon.Post("", middleware.IsLoggedIn, ch.Create)
	coupon.Get("", middleware.IsLoggedIn, ch.FindAll)
	coupon.Get("/code/:code", middleware.IsLoggedIn, ch.FindByCode)
	coupon.Delete("/code/:code", middleware.IsLoggedIn, ch.DeleteByCode)

	ip := api.Group("/iriguchi/ip")
	ip.Get("/get/:ip", middleware.IsLoggedIn, ih.FindByIp)
	ip.Get("", middleware.IsLoggedIn, ih.FindAll)

	customer := api.Group("/iriguchi/customer")
	customer.Post("/send/:emailType", middleware.IsLoggedIn, eh.SendEmail)
	customer.Get("/name", middleware.IsLoggedIn, crh.FindAllByFullName)
	customer.Get("/optin", middleware.IsLoggedIn, crh.FindAllByOptInStatus)
	customer.Get("", middleware.IsLoggedIn, crh.FindAll)
}

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	SetupRoutes(app)

	return app
}
