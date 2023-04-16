package handler

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitRouter() http.Handler {
    jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            return []byte(mySigningKey), nil
        },
        SigningMethod: jwt.SigningMethodHS256,
    })

    router := mux.NewRouter()
    //  those api which have jwtMiddleware.Handler needs token to access
    // router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
    // below are two original routes, which I will rewrite later
    router.Handle("/checkout", jwtMiddleware.Handler(http.HandlerFunc(checkoutHandler))).Methods("POST")    
    router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
    
    // router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
    // router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

    // below are new routes for flag camp
    router.Handle("/uploadprod", jwtMiddleware.Handler(http.HandlerFunc(uploadProductionHandler))).Methods("POST")
    router.Handle("/newsignup", http.HandlerFunc(newSignupHandler)).Methods("POST")
    router.Handle("/newsignin", http.HandlerFunc(newSigninHandler)).Methods("POST")
    // router.Handle("/newsearch", jwtMiddleware.Handler(http.HandlerFunc(newSearchHandler))).Methods("GET")
    router.Handle("/filter", jwtMiddleware.Handler(http.HandlerFunc(filterHandler))).Methods("GET")
    
    originsOk := handlers.AllowedOrigins([]string{"*"})
    headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

    return handlers.CORS(originsOk, headersOk, methodsOk)(router)

}