package service

import (
	"fmt"
	"reflect"

	"appstore/backend"
	"appstore/constants"
	"appstore/model"

	"github.com/olivere/elastic/v7"
)

// func CheckUser(username, password string) (bool, error) {old api, already rewrite below
//     query := elastic.NewBoolQuery()
//     query.Must(elastic.NewTermQuery("username", username))
//     query.Must(elastic.NewTermQuery("password", password))
//     searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
//     if err != nil {
//         return false, err
//     }

//     var utype model.User
//     for _, item := range searchResult.Each(reflect.TypeOf(utype)) {//有点多此一举这边，上面搜索能搜到就已经验证成功了
//         u := item.(model.User)
//         if u.Password == password {
//             fmt.Printf("Login as %s\n", username)
//             return true, nil
//         }
//     }
//     return false, nil
// }

// func AddUser(user *model.User) (bool, error) {old api, already rewrite below
//     query := elastic.NewTermQuery("username", user.Username)
//     searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)// if there is record, username exist, we can not add users
//     if err != nil {
//         return false, err
//     }

//     if searchResult.TotalHits() > 0 {//if there is record, username exist, we can not add users, return false;
//         if err != nil {
//             return false, err
//         }
//         return false, nil
//     }

//     err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
//     if err != nil {
//         return false, err
//     }
//     fmt.Printf("User is added: %s\n", user.Username)
//     return true, nil
// }
//below is for flag camp, AddNewUser
func NewCheckUser(username, password string) (bool, error) {
    query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("username", username))
    query.Must(elastic.NewTermQuery("password", password))
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_PROD_INDEX)
    if err != nil {
        return false, err
    }


    var utype model.User_prod
    for _, item := range searchResult.Each(reflect.TypeOf(utype)) {//有点多此一举这边，上面搜索能搜到就已经验证成功了
        u := item.(model.User_prod)
        if u.Password == password {
            fmt.Printf("Login as %s\n", username)
            return true, nil
        }
    }
    return false, nil
}



  func AddNewUser(user_prod *model.User_prod) (bool, error) {
    query := elastic.NewTermQuery("username", user_prod.Username)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_PROD_INDEX)// if there is record, username exist, we can not add users
    if err != nil {
        return false, err
    }

    if searchResult.TotalHits() > 0 {//if there is record, username exist, we can not add users, return false;
        if err != nil {
            return false, err
        }
        return false, nil
    }

    err = backend.ESBackend.SaveToES(user_prod, constants.USER_PROD_INDEX, user_prod.Username)
    if err != nil {
        return false, err
    }
    fmt.Printf("User_prod is added: %s\n", user_prod.Username)
    return true, nil
}
