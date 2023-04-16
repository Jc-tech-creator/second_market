package model

import "github.com/olivere/elastic/v7"

//original struct, no use now
type App struct {//original struct, no use now
 Id          string `json:"id"`
 User        string `json:"user"`
 Title       string `json:"title"`
 Description string `json:"description"`
 Price       int    `json:"price"`
 Url         string `json:"url"`
 ProductID   string `json:"product_id"`
 PriceID     string `json:"price_id"`
}

type User struct {//original struct, no use now
    Username string `json:"username"`
    Password string `json:"password"`
    Age      int64  `json:"age"`
    Gender   string `json:"gender"`
}

//below is new for flagcamp
type Production struct {
	Id          string           `json:"id"`
	User        string           `json:"user"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       int              `json:"price"`
	Url         string           `json:"url"`
	ProductId   string           `json:"product_id"`
	PriceId     string           `json:"price_id"`
	Location    *elastic.GeoPoint `json:"location"`
}

type User_prod struct {
	Username string           `json:"username"`
	Password string           `json:"password"`
	Age      int64            `json:"age"`
	Gender   string           `json:"gender"`
	Location *elastic.GeoPoint `json:"location"`
}

type UserProdRequest struct {
    Username string  `json:"username"`
    Password string  `json:"password"`
    Age      int64   `json:"age"`
    Gender   string  `json:"gender"`
    Lat      float64 `json:"lat"`
    Lon      float64 `json:"lon"`
}

