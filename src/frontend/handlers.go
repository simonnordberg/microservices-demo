package main

import (
	"fmt"
	"html/template"
	"net/http"

	pb "simonnordberg.com/demoshop/frontend/genproto"
)

var (
	templates = template.Must(template.New("").
		Funcs(template.FuncMap{
			"renderMoney": renderMoney,
		}).ParseGlob("templates/*.html"))
)

func (fe *frontendServer) debugHandler(w http.ResponseWriter, r *http.Request) {
	err := fe.getEnvironments(r.Context())
	if err != nil {
		log.Debugf("error: %v", err)
	}
}

func (fe *frontendServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	/*
		currencies, err := fe.getCurrencies(r.Context())
		if err != nil {
			log.Fatalf("could not retrieve currencies: %v", err)
			return
		}
	*/
	products, err := fe.getProducts(r.Context())
	if err != nil {
		log.Fatalf("could not retrieve products: %v", err)
		return
	}

	type productView struct {
		Item *pb.Product
		//		Price *pb.Money
	}
	ps := make([]productView, len(products))

	for i, p := range products {
		/*
				price, err := fe.convertCurrency(r.Context(), p.GetPriceUsd(), "SEK")
				if err != nil {
					log.Fatalf("could to do currency conversion: %v", err)
					return
				}
			ps[i] = productView{p, price}
		*/
		ps[i] = productView{p}
	}

	if err := templates.ExecuteTemplate(w, "home", map[string]interface{}{
		"currencies": []string{"USD"},
		"products":   ps,
	}); err != nil {
		log.Error(err)
	}
}

func renderMoney(money *pb.Money) string {
	currencyLogo := renderCurrencyLogo(money.GetCurrencyCode())
	return fmt.Sprintf("%s%d.%02d", currencyLogo, money.GetUnits(), money.GetNanos()/10000000)
}

func renderCurrencyLogo(currencyCode string) string {
	logos := map[string]string{
		"USD": "$",
		"CAD": "$",
		"JPY": "¥",
		"EUR": "€",
		"TRY": "₺",
		"GBP": "£",
	}

	logo := "$" //default
	if val, ok := logos[currencyCode]; ok {
		logo = val
	}
	return logo
}
