package main

import (
	"log"
	"net/http"
)

func (h *handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	sess, err := h.store.Get(r, "oidc-session")
	if err != nil {
		log.Println(err)
	}

	err = refresh(r.Context(), h.conf, sess)
	if err != nil {
		log.Println(err)
	}

	execTemplate(r.Context(), sess, w)
}
