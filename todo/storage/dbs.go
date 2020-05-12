package storage

import (
	"database/sql"
	"todo"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"


)


type Options struct{
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (o Options) Validate() error {
	switch {
	case o.Host == "":
		return errors.New("host is required")
	case o.Port == 0:
		return errors.New("port is required")
	case o.User == "":
		return errors.New("user is required")
	case o.Password == "":
		return errors.New("password is required")
	case o.Name == "":
		return errors.New("name is required")
	default:
		return nil
	}
}
type Client struct {
	db *sql.DB
}

func New(opts Options) (*Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		opts.Host, opts.Port, opts.User, opts.Password, opts.Name)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) AddUser(user todo.User) error {
	_, err := c.db.Exec(`INSERT INTO users (id, name) VALUES (?, ?)`, user.ID, user.Name)
	return err
}

func (c *Client) Users() ([]todo.User, error) {
	rows, err := c.db.Query(`SELECT name FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []todo.User
	for rows.Next() {
		var user todo.User
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (c *Client) RemoveUser(userID string) error {
	_, err := c.db.Exec(`DELETE FROM users WHERE id = ?`, userID)
	return err
}

func (c *Client) UpdateUser(user todo.User) error {
	_, err := c.db.Exec(`UPDATE users SET name = ? WHERE id = ?`, user.Name, user.ID)
	return err
}

func (c *Client) User(userID string) (todo.User, error) {
	row := c.db.QueryRow(`SELECT id, name FROM users WHERE id = ?`, userID)
	var user todo.User
	err := row.Scan(&user.ID, &user.Name)
	return user, err
}
