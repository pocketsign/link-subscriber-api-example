package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
)

func (h *handler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	state := randomString()

	sess, err := h.store.Get(r, "oidc-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// セッションにstateを保存
	sess.Values["state"] = state
	if err := sess.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, h.conf.AuthCodeURL(state), http.StatusFound)
}

func (h *handler) handleCallback(w http.ResponseWriter, r *http.Request) {
	sess, err := h.store.Get(r, "oidc-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// stateが一致しているかどうか検証
	if r.URL.Query().Get("state") != sess.Values["state"] {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	// codeからtokenを取得
	oauth2Token, err := h.conf.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// id_tokenを取得
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token", http.StatusInternalServerError)
		return
	}

	// id_tokenの検証
	verifier := h.provider.Verifier(&oidc.Config{ClientID: h.conf.ClientID})
	idtoken, err := verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// セッションにsubjectを保存
	sess.Values["sub"] = idtoken.Subject

	if err := sess.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// subjectをキーにして取得したトークンを保存
	src := h.conf.TokenSource(r.Context(), oauth2Token)
	tokenSourceStore[idtoken.Subject] = src
	refreshTokenStore[idtoken.Subject] = oauth2Token.RefreshToken
	http.Redirect(w, r, "/", http.StatusFound)
}

func randomString() string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
