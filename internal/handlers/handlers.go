package handlers

import (
	"errors"
	"html/template"
	"main/internal/logger"
	"main/internal/usecase"
	"net/http"
	"regexp"
)

type Handler struct {
	logic  usecase.UseCase
	logger logger.ILogger
}

func New(logic usecase.UseCase, logger logger.ILogger) *Handler {
	return &Handler{logic: logic, logger: logger}
}

func (h Handler) MainHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		h.logger.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, "")
}

func (h Handler) ResultHandler(w http.ResponseWriter, r *http.Request) {
	addres := r.FormValue("address")
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		h.logger.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data usecase.IPAddress
	regStr := `^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}/(1[0-9]{1}|2[0-9]{1}|3[0-2]{1}|[0-9])$`
	reg, err := regexp.Compile(regStr)
	if err != nil {
		data.Error = errors.New("Server error!")
		h.logger.Warn(err.Error())
		t.Execute(w, data)
		return
	}

	matched := reg.MatchString(addres)

	if !matched {
		data.Error = errors.New("Wrong Data!")
		t.Execute(w, data)
		return
	}

	data, err = h.logic.GetInfoByIP(addres)

	t.Execute(w, data)
}
