package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	tmpl "github.com/pocketsign/link-subscriber-api-example/template"
)

const (
	// 各種リソースのID
	nameResourceID      = "50f423aa-8ea8-4385-81a1-a7d8da1c4939"
	birthDateResourceID = "612c6a3c-212a-402e-8c89-792a4b3e7889"
	addressResourceID   = "5419c333-2863-4d57-9272-5f6c6c8989ec"
	genderResourceID    = "5ac941bf-88f3-47e3-b2e7-3bb0522aa14b"
)

type templateUserResources struct {
	Name      string
	BirthDate string
	Address   string
	Gender    string
	Custom    string
}

func getTemplateUserResources(ctx context.Context, sess *sessions.Session) *templateUserResources {
	res := &templateUserResources{
		Name:      "データの取得に失敗しました",
		BirthDate: "データの取得に失敗しました",
		Address:   "データの取得に失敗しました",
		Gender:    "データの取得に失敗しました",
		Custom:    "",
	}

	// 各種ユーザーリソースをAPIを通して取得
	name, err := getUserResource(ctx, sess, nameResourceID)
	if err != nil {
		log.Println(err)
	} else {
		res.Name = name.Value
	}

	birthDate, err := getUserResource(ctx, sess, birthDateResourceID)
	if err != nil {
		log.Println(err)
	} else {
		res.BirthDate = birthDate.Value
	}

	address, err := getUserResource(ctx, sess, addressResourceID)
	if err != nil {
		log.Println(err)
	} else {
		res.Address = address.Value
	}

	gender, err := getUserResource(ctx, sess, genderResourceID)
	if err != nil {
		log.Println(err)
	} else {
		res.Gender = gender.Value
	}

	custom, err := getUserResource(ctx, sess, customResourceID)
	if err != nil {
		log.Println(err)
	} else {
		res.Custom = custom.Value
	}

	return res
}

func execTemplate(ctx context.Context, sess *sessions.Session, w http.ResponseWriter) {
	tmpl, err := template.ParseFS(tmpl.F, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := getTemplateUserResources(ctx, sess)

	err = tmpl.Execute(w, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
