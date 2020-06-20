package controllers

import (
	"net/http"

	"github.com/zongpoljkk/go_jwt/utils"
)

type Controller struct{}

func (c Controller) ProtectedEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, "It worked!")
	}
}
