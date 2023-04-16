package handler

import (
	"appstore/model"
	"appstore/service"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/olivere/elastic/v7"

	jwt "github.com/form3tech-oss/jwt-go"
)
var mySigningKey = []byte("secret")
//below are original methods, I comment out
// func signinHandler(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("Received one signin request")
//     w.Header().Set("Content-Type", "text/plain")

//     //  Get User information from client
//     decoder := json.NewDecoder(r.Body)
//     var user model.User
//     if err := decoder.Decode(&user); err != nil {
//         http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
//         fmt.Printf("Cannot decode user data from client %v\n", err)
//         return
//     }
//     exists, err := service.CheckUser(user.Username, user.Password)
//     if err != nil {
//         http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
//         fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
//         return
//     }

//     if !exists {
//         http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
//         fmt.Printf("User doesn't exists or wrong password\n")
//         return
//     }

//     token := jwt.NewWithClaims(
//         jwt.SigningMethodHS256, //first part of token
//         jwt.MapClaims{//second part of token
//         "username": user.Username,
//         "exp":      time.Now().Add(time.Hour * 24).Unix(),
//     })

//     tokenString, err := token.SignedString(mySigningKey) // encode first two part and get complete token
//     if err != nil {
//         http.Error(w, "Failed to generate token", http.StatusInternalServerError)
//         fmt.Printf("Failed to generate token %v\n", err)
//         return
//     }
//     w.Write([]byte(tokenString))
// }

// func signupHandler(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("Received one signup request")
//     w.Header().Set("Content-Type", "text/plain")

//     decoder := json.NewDecoder(r.Body)
//     var user model.User
//     if err := decoder.Decode(&user); err != nil {
//         http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
//         fmt.Printf("Cannot decode user data from client %v\n", err)
//         return
//     }



//     if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {// some validation
//         http.Error(w, "Invalid username or password", http.StatusBadRequest)
//         fmt.Printf("Invalid username or password\n")
//         return
//     }

//     success, err := service.AddUser(&user)  
//     if err != nil {
//         http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
//         fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
//         return
//     }

//     if !success {
//         http.Error(w, "User already exists", http.StatusBadRequest)
//         fmt.Println("User already exists")
//         return
//     }
//     fmt.Printf("User added successfully: %s.\n", user.Username)
// }
 // for new flag camp, it is new signuphandler, add the Location
 func newSignupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signup request")
    w.Header().Set("Content-Type", "text/plain")

    decoder := json.NewDecoder(r.Body)
    var userProdReq model.UserProdRequest
    if err := decoder.Decode(&userProdReq); err != nil {
        http.Error(w, "Cannot decode userProdReq data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode userProdReq data from client %v\n", err)
        return
    }


    if userProdReq.Username == "" || userProdReq.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(userProdReq.Username) {// some validation
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

    location := elastic.GeoPointFromLatLon(userProdReq.Lat, userProdReq.Lon)

    user_prod := model.User_prod{
        Username: userProdReq.Username,
        Password: userProdReq.Password,
        Age:      userProdReq.Age,
        Gender:   userProdReq.Gender,
        Location: location,
    }

    success, err := service.AddNewUser(&user_prod)  
    if err != nil {
        http.Error(w, "Failed to save user_prod to Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user_prod to Elasticsearch %v\n", err)
        return
    }

    if !success {
        http.Error(w, "user_prod already exists", http.StatusBadRequest)
        fmt.Println("user_prod already exists")
        return
    }
    fmt.Printf("user_prod added successfully: %s.\n", user_prod.Username)
}

func newSigninHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signin request")
    w.Header().Set("Content-Type", "text/plain")

    //  Get User information from client
    decoder := json.NewDecoder(r.Body)
    var user_prod model.User_prod
    if err := decoder.Decode(&user_prod); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }
    exists, err := service.NewCheckUser(user_prod.Username, user_prod.Password)
    if err != nil {
        http.Error(w, "Failed to read user_prod from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user_prod from Elasticsearch %v\n", err)
        return
    }

    if !exists {
        http.Error(w, "User_prod doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User_prod doesn't exists or wrong password\n")
        return
    }

    token := jwt.NewWithClaims(
        jwt.SigningMethodHS256, //first part of token
        jwt.MapClaims{//second part of token
        "username": user_prod.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(mySigningKey) // encode first two part and get complete token
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return
    }
    w.Write([]byte(tokenString))
}
