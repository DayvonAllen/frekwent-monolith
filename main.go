package main

import (
	"fmt"
	"freq/database"
	"freq/models"
	"freq/repository"
	"freq/router"
	"github.com/globalsign/mgo/bson"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var scheduler *cron.Cron

func init() {
	conn := database.Sess

	user := new(models.User)
	err := conn.DB(database.DB).C(database.ADMIN).Find(bson.D{{"email", "admin@admin.com"}}).One(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			admin := new(models.User)
			admin.Email = "admin@admin.com"
			admin.Username = "admin"
			admin.Password = "password"
			admin.CreatedAt = time.Now()
			admin.UpdatedAt = time.Now()

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
			admin.Password = string(hashedPassword)

			err = conn.DB(database.DB).C(database.ADMIN).Insert(admin)
			if err != nil {
				return
			}
		}
	}
	models.Instance = models.CreateMailer()
	go models.Instance.ListenForMail()

	// cron job for resending failed emails
	scheduler = cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))

	_, err = scheduler.AddFunc("@every 4h30m", func() {
		// TODO check every 4 hours and 30 minutes for failed emails and try to resend them
		stat := models.Failed
		emails, err := repository.EmailRepoImpl{}.FindAllByStatusRaw(&stat)
		if err != nil {
			fmt.Println("none found")
		}

		var wg sync.WaitGroup
		wg.Add(len(*emails))

		for _, email := range *emails {
			email := email
			go func(e models.Email) {
				defer wg.Done()
				models.SendMessage(&email)
			}(email)
		}
	})

	if err != nil {
		panic(err)
	}

	scheduler.Start()
	//_, err = conn.PurchaseCollection.DeleteMany(context.TODO(), bson.M{})
	//if err != nil {
	//	return
	//}
	//
	//_, err = conn.CustomerCollection.DeleteMany(context.TODO(), bson.M{})
	//if err != nil {
	//	return
	//}
	//_, err = conn.EmailCollection.DeleteMany(context.TODO(), bson.M{})
	//if err != nil {
	//	return
	//}
}

func main() {
	app := router.Setup()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}
}
