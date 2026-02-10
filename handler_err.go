package main

import (
	"net/http"

	"github.com/Aym-Aymen777/RSS-Aggregator/utils"
)

func handlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 400, "Something went wrong")
}