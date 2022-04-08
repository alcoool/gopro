package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	actions("jon.doe@company.com", "1234")
	//go actions("jon.doe@company.com", "1234")
	//go actions("jon.doe2@company.com", "1234")
}

func actions(email string, password string) {
	token := login(email, password)

	addToCart(token, 1, 2)
	addToCart(token, 2, 1)

	checkout(token)
}

func login(email string, password string) string {
	data := url.Values{}
	data.Set("email", email)
	data.Set("password", password)

	resp, err := http.PostForm("http://localhost:8080/login", data)
	if err != nil {
		log.Fatalln(err)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "jwt" {
			err := resp.Body.Close()
			if err != nil {
				log.Fatalln(err)
			}
			return cookie.Value
		}
	}

	log.Fatalln(errors.New("no token found"))

	return ""
}

func addToCart(token string, productId int, quantity int) {
	data := url.Values{}
	data.Set("productId", strconv.Itoa(productId))
	data.Set("quantity", strconv.Itoa(quantity))

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/add-to-cart", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "jwt="+token)

	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Print(body)

	err = resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func checkout(token string) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/checkout", strings.NewReader(""))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Cookie", "jwt="+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	log.Print(body)
}
