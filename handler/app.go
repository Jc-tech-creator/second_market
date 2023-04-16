package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"appstore/backend"
	"appstore/constants"
	"appstore/model"
	"appstore/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/olivere/elastic/v7"
	"github.com/pborman/uuid"
)

// func uploadHandler(w http.ResponseWriter, r *http.Request) {// the entrance for uploading app
//     // Parse r *http.Request to build app
//     fmt.Println("Received one upload request")
//     user := r.Context().Value("user")
//     claims := user.(*jwt.Token).Claims
//     username := claims.(jwt.MapClaims)["username"]

//     app := model.App{
//         Id:          uuid.New(),// we new a id from strach
//         User:        username.(string),
//         Title:       r.FormValue("title"),
//         Description: r.FormValue("description"),
//     }

//   price, err := strconv.Atoi(r.FormValue("price"))
//   // convert a string representation of an integer to an actual integer value
//     fmt.Printf("%v,%T", price, price)
//     if err != nil {
//         fmt.Println(err)
//     }
//     app.Price = price

//     file, _, err := r.FormFile("media_file")
//     if err != nil {
//         http.Error(w, "Media file is not available", http.StatusBadRequest)
//         fmt.Printf("Media file is not available %v\n", err)
//         return
//     }

//     err = service.SaveApp(&app, file)
//     if err != nil {
//         http.Error(w, "Failed to save app to backend", http.StatusInternalServerError)
//         fmt.Printf("Failed to save app to backend %v\n", err)
//         return
//     }

//     fmt.Println("App is saved successfully.")

//     fmt.Fprintf(w, "App is saved successfully: %s\n", app.Description)
// }


func searchHandler(w http.ResponseWriter, r *http.Request) {//original searchHandler, will rewrite
    fmt.Println("Received one search request")
    w.Header().Set("Content-Type", "application/json")
    //get params from request
    title := r.URL.Query().Get("title")
    description := r.URL.Query().Get("description")
 
    //call service to handle request
    var apps []model.App
    var err error
    apps, err = service.SearchApps(title, description)
    if err != nil {
        http.Error(w, "Failed to read Apps from backend", http.StatusInternalServerError)
        return
    }
 
    //construct response
    js, err := json.Marshal(apps)
    //it serializes the Go data into JSON format
    if err != nil {
        http.Error(w, "Failed to parse Apps into JSON format", http.StatusInternalServerError)
        return
    }
    w.Write(js)
 }

 func checkoutHandler(w http.ResponseWriter, r *http.Request) {// original checkoutHandler, will rewrite
    fmt.Println("Received one checkout request")
    w.Header().Set("Content-Type", "text/plain")
    //get params from request
    appID := r.FormValue("appID")
    //call service
    s, err := service.CheckoutApp(r.Header.Get("Origin"), appID)
    if err != nil {
        fmt.Println("Checkout failed.")
        w.Write([]byte(err.Error()))
        return
    }
    //construct URL
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(s.URL))
 
    fmt.Println("Checkout process started!")
 }
 
 // below one is new handler for upload production for flagcamp
 func uploadProductionHandler(w http.ResponseWriter, r *http.Request) {// the entrance for uploading production
    // Parse r *http.Request to build production
    fmt.Println("Received one upload request")
    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]
    //convert lat and lon into GeoPoint
    lat := r.FormValue("lat")
    lon := r.FormValue("lon")
    location, err := elastic.GeoPointFromString(fmt.Sprintf("%s,%s", lat, lon))
    if err != nil {
        fmt.Printf("Failed to convert lat and lon into geo.point %v\n", err)
        return
    }
    fmt.Println("location")

    production := model.Production{
        Id:          uuid.New(),// we new a id from strach
        User:        username.(string),
        Name:       r.FormValue("name"),
        Description: r.FormValue("description"),
        Location: location,
    }

  price, err := strconv.Atoi(r.FormValue("price"))
  // convert a string representation of an integer to an actual integer value
    fmt.Printf("%v,%T", price, price)
    if err != nil {
        fmt.Println(err)
    }
    production.Price = price
 
 
    file, _, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }
 //
    err = service.SaveProduction(&production, file)
    if err != nil {
        http.Error(w, "Failed to save production to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save production to backend %v\n", err)
        return
    }
 
 
    fmt.Println("Production is saved successfully.")
    fmt.Println(production.Id)
    fmt.Println(production.Name)
    fmt.Println(production.Price)
    fmt.Println(production.Location)
    fmt.Fprintf(w, "Production is saved successfully: %s\n", production.Description)
}
 
 // below is the filterHandler:
 func filterHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one search request")
    w.Header().Set("Content-Type", "application/json")
    //get params from request
    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]
    price := r.URL.Query().Get("price")
    distance := r.URL.Query().Get("distance")
    // username := r.URL.Query().Get("username")
    // Retrieve user's location from Elasticsearch using the provided username
    userLocation, error := getUserLocation(username.(string))
    if error != nil {
        http.Error(w, "Failed to retrieve user location from Elasticsearch", http.StatusInternalServerError)
        return
    }
    //call service to handle request
    var productions []model.Production

    productions, err := service.FilterProduction(price, distance, userLocation)
    
    if err != nil {
        http.Error(w, "Failed to read Productions from backend", http.StatusInternalServerError)
        return
    }
    fmt.Println("this is searched productions");
    fmt.Println(productions);
    //construct response
    js, err := json.Marshal(productions)
    //it serializes the Go data into JSON format
    if err != nil {
        http.Error(w, "Failed to parse Productions into JSON format", http.StatusInternalServerError)
        return
    }
    w.Write(js)
 }

 func getUserLocation(username string) (*elastic.GeoPoint, error) {//this function is necessary to get UserLocation
    // Implement logic to query Elasticsearch for the user's location using their username
    // user_prod *model.User_prod
    // Return the user's location as an elastic.GeoPoint
    query := elastic.NewTermQuery("username", username)
    searchedUsers, err := backend.ESBackend.ReadFromES(query, constants.USER_PROD_INDEX)
    if err != nil {
        fmt.Printf("Failed to read production from backend %v\n", err)
    }
    var utype model.User_prod
    var geopoints [] *elastic.GeoPoint
    for _, item := range searchedUsers.Each(reflect.TypeOf(utype)) {
        u := item.(model.User_prod)
        geopoints = append(geopoints, u.Location)
    }
    // fmt.Printf("geopoints[0]", geopoints[0])
    return geopoints[0], err
}
 