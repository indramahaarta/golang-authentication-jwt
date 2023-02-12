package maincontroller

import (
	"net/http"

	"github.com/indramahaarta/golang-authentication-jwt/helpers"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

type M map[string]interface{}

func Product(w http.ResponseWriter, r *http.Request) {
	var product []M

	m1 := M{
		"name":  "sepatu",
		"price": "10000",
		"stock": "10",
	}

	m2 := M{
		"name":  "sepatu",
		"price": "10000",
		"stock": "10",
	}

	m3 := M{
		"name":  "sepatu",
		"price": "10000",
		"stock": "10",
	}

	product = append(product, m1)
	product = append(product, m2)
	product = append(product, m3)

	helpers.ResponseJSON(w, http.StatusAccepted, product)
}
