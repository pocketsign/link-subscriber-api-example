package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	sess, err := h.store.Get(r, "oidc-session")
	if err != nil {
		log.Println(err)
	}

	execTemplate(r.Context(), sess, w)
}

func (h *handler) handlePut(w http.ResponseWriter, r *http.Request) {
	sess, err := h.store.Get(r, "oidc-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := struct {
		Value string `json:"value"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// valueに渡す値を一度jsonに変換
	b, err := json.Marshal(v.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 変換したjsonをvalueに渡す
	rv := UserResource{
		Value: string(b),
	}

	if err := putUserResource(r.Context(), sess, customResourceID, rv); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
