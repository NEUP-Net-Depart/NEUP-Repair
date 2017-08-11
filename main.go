package main

import (
	"net/http"

	"github.com/NEUP-Net-Depart/NEUP-Repair/handler/middleware"

	"github.com/NEUP-Net-Depart/NEUP-Repair/config"
	"github.com/NEUP-Net-Depart/NEUP-Repair/handler"
	"github.com/NEUP-Net-Depart/NEUP-Repair/model"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

func init() {
	db, err := sqlx.Open("mysql", config.GlobalConfig.DSN)
	if err != nil {
		log.Fatal(err)
	}
	model.GlobalDB = db
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("main module inited")
}

func main() {
	router := httprouter.New()
	router.POST("/api/v1/orders", handler.AddOrder)
	router.PUT("/api/v1/orders/:secret", handler.FinishOrder)
	router.GET("/api/v1/orders/:secret", handler.OrdersSecret)
	router.PUT("/api/v1/orders/:secret/rate", handler.RateOrder)
	router.GET("/api/v1/orders/:secret/comments", handler.Comments)
	router.GET("/api/v1/orders", handler.OrdersByPage)
	router.POST("/api/v1/auth", handler.Auth)
	router.GET("/api/v1/announce", handler.Announce)
	router.PUT("/api/v1/announce", handler.SetAnnounce)
	m := negroni.New()
	m.Use(negroni.NewStatic(http.Dir("resources/NEUP-Repair-frontend/app")))
	m.UseFunc(middleware.CORS)
	m.UseFunc(middleware.CheckAuth)
	m.UseFunc(middleware.Logger)
	m.UseHandler(router)

	log.Fatal(http.ListenAndServe(config.GlobalConfig.ServAddr, m))
}
