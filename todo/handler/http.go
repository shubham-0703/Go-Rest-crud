package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"todo"
	"todo/storage"
)
import "github.com/gorilla/mux"


type Client struct {
	db *storage.Client
}

func New(db *storage.Client) *Client {
	return &Client{db: db}
}

func (c *Client) Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/users", c.Users).Methods(http.MethodGet)
	r.HandleFunc("/users", c.AddUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{userID}", c.User).Methods(http.MethodGet)
	r.HandleFunc("/users/{userID}", c.DeleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/users/{userID}", c.UpdateUser).Methods(http.MethodPut)

	srv := &http.Server{Handler: r, Addr: ":8080"}
	return srv.ListenAndServe()
}

func (c *Client) Users(resp http.ResponseWriter, _ *http.Request) {
	users, err := c.db.Users()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	jsn, err := json.Marshal(users)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	resp.Write(jsn)
}

func (c *Client) User(resp http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["userID"]
	if userID == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := c.db.User(userID)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}

	jsn, err := json.Marshal(user)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
	}

	resp.Write(jsn)
}

func (c *Client) DeleteUser(resp http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["userID"]
	if userID == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	err := c.db.RemoveUser(userID)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
}

func (c *Client) UpdateUser(resp http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["userID"]
	if userID == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	var user todo.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(err.Error()))
		return
	}

	err = c.db.UpdateUser(user)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(err.Error()))
		return
	}
}

func (c *Client) AddUser(resp http.ResponseWriter, req *http.Request) {
	var user todo.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(err.Error()))
		return
	}

	user.ID = uuid.New().String()
	err = c.db.AddUser(user)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(err.Error()))
		return
	}
}
