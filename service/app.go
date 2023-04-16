package service

import (
	"appstore/backend"
	"appstore/constants"
	"appstore/model"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"

	"github.com/olivere/elastic/v7"
	"github.com/stripe/stripe-go/v74"
)

//below are all original api related to search, which I would rewrite or already rewrite
func SearchApps(title string, description string) ([]model.App, error) {
   if title == "" {
       return SearchAppsByDescription(description)
   }
   if description == "" {
       return SearchAppsByTitle(title)
   }


   query1 := elastic.NewMatchQuery("title", title)
   query2 := elastic.NewMatchQuery("description", description)
   query := elastic.NewBoolQuery().Must(query1, query2)// requires both match queries to be satisfied.
   searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
   if err != nil {
       return nil, err
   }


   return getAppFromSearchResult(searchResult), nil
}


func SearchAppsByTitle(title string) ([]model.App, error) {
   query := elastic.NewMatchQuery("title", title)
   query.Operator("AND")// query will look for documents that contain all words in title field
   if title == "" {
       query.ZeroTermsQuery("all")
   }
   searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
   if err != nil {
       return nil, err
   }

   return getAppFromSearchResult(searchResult), nil
}

func SearchAppsByDescription(description string) ([]model.App, error) {
   query := elastic.NewMatchQuery("description", description)
   query.Operator("AND")
   if description == "" {
       query.ZeroTermsQuery("all")
   }
   searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
   if err != nil {
       return nil, err
   }


   return getAppFromSearchResult(searchResult), nil
}

func SearchAppsByID(appID string) (*model.App, error) {
   query := elastic.NewTermQuery("id", appID)
   searchResult, err := backend.ESBackend.ReadFromES(query, constants.APP_INDEX)
   if err != nil {
       return nil, err
   }
   results := getAppFromSearchResult(searchResult)
   if len(results) == 1 {
       return &results[0], nil
   }
   return nil, nil
}


func getAppFromSearchResult(searchResult *elastic.SearchResult) []model.App {
   var ptype model.App
   var apps []model.App
   for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
       p := item.(model.App)
       apps = append(apps, p)
   }
   return apps
}

// func SaveApp(app *model.App, file multipart.File) error {// old api, already rewrite
//     //Call Stripe to get PriceID and ProductID
//     productID, priceID, err := backend.CreateProductWithPrice(app.Title, app.Description, int64(app.Price*100))
//     if err != nil {
//         fmt.Printf("Failed to create Product and Price using Stripe SDK %v\n", err)
//         return err
//     }
//     app.ProductID = productID
//     app.PriceID = priceID
//     fmt.Printf("Product %s with price %s is successfully created", productID, priceID)
     
//     //Call GCS to store file and get URL and add to property of app
//     medialink, err := backend.GCSBackend.SaveToGCS(file, app.Id)
//     if err != nil {
//         return err
//     }
//     app.Url = medialink
    
//     // Save Everything to ES
//     err = backend.ESBackend.SaveToES(app, constants.APP_INDEX, app.Id)
//     if err != nil {
//         fmt.Printf("Failed to save app to elastic search with app index %v\n", err)
//         return err
//     }
//     fmt.Println("App is saved successfully to ES app index.")
 
//     return nil
//  }
 func CheckoutApp(domain string, appID string) (*stripe.CheckoutSession, error) {
    // map appId to priceId
    app, err := SearchAppsByID(appID)
    //find the app
    if err != nil {
        return nil, err
    }
    if app == nil {
        return nil, errors.New("unable to find app in elasticsearch")
    }
    // call stripe checkout app
    return backend.CreateCheckoutSession(domain, app.PriceID)

 }

 // flag camp
 // saveProduction for flag camp
 func SaveProduction(production *model.Production, file multipart.File) error {
    //Call Stripe to get PriceID and ProductID
    productID, priceID, err := backend.CreateProductWithPrice(production.Name, production.Description, int64(production.Price*100))
    if err != nil {
        fmt.Printf("Failed to create Product and Price using Stripe SDK %v\n", err)
        return err
    }
    production.ProductId = productID
    production.PriceId = priceID
    fmt.Printf("Product %s with price %s is successfully created", productID, priceID)
     
    //Call GCS to store file and get URL and add to property of production
    medialink, err := backend.GCSBackend.SaveToGCS(file, production.Id)
    if err != nil {
        return err
    }
    production.Url = medialink
    
    // Save Everything to ES, store to Production_INDEX
    err = backend.ESBackend.SaveToES(production, constants.PRODUCTION_INDEX, production.Id)
    if err != nil {
        fmt.Printf("Failed to save production to elastic search with production index %v\n", err)
        return err
    }
    fmt.Println("Production is saved successfully to ES production index.")
 
    return nil
 }
// filterProduction for flag camp
 func FilterProduction(price string, distance string, userLocation *elastic.GeoPoint) ([]model.Production, error) {
    fmt.Println("Here, below is the price, distance, and userLocation")
    fmt.Println(price)
    fmt.Println(distance)
    fmt.Println(userLocation)
	if price == "-1" { // I don't know how front end express none input price filter, here -1 just means default
        distance_int, err := strconv.Atoi(distance)
        if err != nil {
            fmt.Println("Error converting string to int64:", err)
            return nil, err
        } 
		return SearchProductionsByDistance(distance_int, userLocation)
	}
	if distance == "-1" { // I don't know how front end express none input distance filter, here -1 just means default
        fmt.Println("We enter distance == -1")
        price_int, err := strconv.Atoi(price)
        if err != nil {
            fmt.Println("Error converting string to int64:", err)
            return nil, err
        }     
		return SearchProductionsByPrice(price_int)
	}
    fmt.Println("We have both distance and price, we would filter on both")
	searchResult, err := filterProductionByBothDistanceAndPrice(distance, userLocation, price)
	if err != nil {
		return nil, err
	}
	return getProductionsFromSearchResult(searchResult), nil
}

func SearchProductionsByPrice(price_int int) ([]model.Production, error) {
    return nil, nil
}

func filterProductionByBothDistanceAndPrice(distance string, userLocation *elastic.GeoPoint, price string)(*elastic.SearchResult, error){
    return nil, nil
}

func getProductionsFromSearchResult(searchResult *elastic.SearchResult) []model.Production {
    fmt.Println("we enter getProductionsFromSearchResult")
	var ptype model.Production
	var productions []model.Production
	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(model.Production)
		productions = append(productions, p)
	}
    fmt.Println("this is returned productions from getProductionsFromSearchResult")
    fmt.Println(productions)
	return productions
}

func SearchProductionsByDistance(distance int, userLocation *elastic.GeoPoint) ([]model.Production, error) {
	// this function would return the array of Production within the range of distance of user location.
	// Create a Geo Distance Query
    fmt.Println("Hello, we enter the SearchProductionsByDistance")
    fmt.Println(distance)
    fmt.Println(userLocation)
	geoDistanceQuery := elastic.NewGeoDistanceQuery("location").
		Distance(fmt.Sprintf("%dm", distance*1609)). // Convert miles to meters
		Point(userLocation.Lat, userLocation.Lon)

	// geoDistanceQuery.Operator("AND") // query will look for documents that contain all words in title field
	searchResult, err := backend.ESBackend.ReadFromES(geoDistanceQuery, constants.PRODUCTION_INDEX)
    fmt.Println("this is searchResult")
    fmt.Println(searchResult)
	if err != nil {
		return nil, err
	}
	return getProductionsFromSearchResult(searchResult), nil
}
 