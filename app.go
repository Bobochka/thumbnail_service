package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"net/url"

	"fmt"

	"log"

	"github.com/Bobochka/thumbnail_service/lib"
	"github.com/Bobochka/thumbnail_service/lib/service"
	"github.com/Bobochka/thumbnail_service/lib/transform"
)

type App struct {
	service *service.Service
}

type params struct {
	url    string
	width  int
	height int
}

func (app *App) thumbnail(w http.ResponseWriter, r *http.Request) {
	params, err := app.thumbnailParams(r)
	if err != nil {
		app.renderError(w, err)
		return
	}

	t := transform.NewLPad(params.width, params.height)

	img, err := app.service.Perform(params.url, t)

	if err != nil {
		app.renderError(w, err)
		return
	}

	app.renderImg(w, img)
}

const maxAllowedArea = 6000000 // px

func (app *App) thumbnailParams(r *http.Request) (params, error) {
	var res params

	res.url = r.URL.Query().Get("url")
	_, err := url.ParseRequestURI(res.url)
	if err != nil {
		msg := fmt.Sprintf("url %s is not valid", res.url)
		return params{}, lib.NewError(err, lib.InvalidParams, msg)
	}

	w := r.URL.Query().Get("width")
	res.width, err = strconv.Atoi(w)
	if err != nil || res.width <= 0 {
		err = fmt.Errorf("width %s is not valid: should be positive integer", w)
		return params{}, lib.NewError(err, lib.InvalidParams, err.Error())
	}

	h := r.URL.Query().Get("height")
	res.height, err = strconv.Atoi(h)
	if err != nil || res.height <= 0 {
		err = fmt.Errorf("height %s is not valid: should be positive integer", h)
		return params{}, lib.NewError(err, lib.InvalidParams, err.Error())
	}

	if res.height*res.width > maxAllowedArea {
		err = fmt.Errorf("requested size of %v x %v is too big", res.width, res.height)
		return params{}, lib.NewError(err, lib.InvalidParams, err.Error())
	}

	return res, nil
}

func (app *App) renderImg(w http.ResponseWriter, img []byte) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(img)))
	w.Write(img)
}

func (app *App) renderError(w http.ResponseWriter, err error) {
	code := 500
	msg := "sorry, but something went wrong"

	if codedError, ok := err.(lib.Error); ok {
		code = codedError.Code()
		msg = codedError.Msg()
	}

	w.WriteHeader(code)

	response := struct{ Error string }{msg}
	data, e := json.Marshal(response)

	//log.Printf("data: %+v\n", data)

	if e != nil {
		log.Println("error marshaling response: ", e.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
