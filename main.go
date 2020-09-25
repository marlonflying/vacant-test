package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	connPort = "8080"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"getUsers",
		"GET",
		"/users",
		getUsers,
	},
	Route{
		"getUser",
		"GET",
		"/user/{id}",
		getUser,
	},
	Route{
		"addUser",
		"POST",
		"/user/add",
		addUser,
	},
	Route{
		"updateUser",
		"PUT",
		"/user/update",
		updateUser,
	},
	Route{
		"deleteUser",
		"DELETE",
		"/user/delete/{id}",
		deleteUser,
	}}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	} `json:"company"`
}

type Users []User

var users []User

func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Print("Get all users...")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vasrs := mux.Vars(r)
	id := vasrs["id"]
	index := getIndex(id)
	if index == -1 {
		log.Print("User not found with Id: ", id)
		http.Error(w, "User not found!", 404)
		return
	}
	user := users[index]
	log.Print("Get user Id :: ", id)
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vasrs := mux.Vars(r)
	id := vasrs["id"]
	log.Print("deleting user Id :: ", id)
	index := getIndex(id)
	if index == -1 {
		http.Error(w, "User not found!", 404)
		return
	}
	users = append(users[:index], users[index+1:]...)
	json.NewEncoder(w).Encode(users)
}

func getIndex(id string) int {
	user := User{}
	for i := 0; i < len(users); i++ {
		user = users[i]
		if strconv.Itoa(user.Id) == id {
			return i
		}
	}
	return -1
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print("error decoding user data :: ", err)
		return
	}
	var isUpsert = true
	for idx, usr := range users {
		if usr.Id == user.Id {
			isUpsert = false
			log.Printf("Updating user id :: %s ", user.Id)
			users[idx] = user
			break
		}
	}
	if isUpsert {
		log.Printf("Upserting user id :: %s", user.Id)
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print("Error decoding user data :: ", err)
		return
	}
	log.Printf("Adding user id :: %d ", user.Id)
	users = append(users, user)
	json.NewEncoder(w).Encode(users)
}

func AddRoutes(router *mux.Router) *mux.Router {
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func main() {
	muxRouter := mux.NewRouter().StrictSlash(true)
	router := AddRoutes(muxRouter)
	log.Printf("Server listening at port: %s", connPort)
	err := http.ListenAndServe(":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}
