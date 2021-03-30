package main

import (
    "stt_back/bootstrap"
    "stt_back/router"
    "stt_back/settings"
    "fmt"
    "net/http"
)

func main() {

    // делаем автомиграцию
    bootstrap.FillDBTestData()

    if settings.IsDev() {
        fmt.Println("Running in DEV mode")
    } else {
        fmt.Println("Running in PROD mode")
    }

    runHttpServer()
}

func runHttpServer() {

	fmt.Println("API сервер запущен :" + settings.ServerPort)
	http.ListenAndServe("0.0.0.0:" + settings.ServerPort, router.Router())
}
