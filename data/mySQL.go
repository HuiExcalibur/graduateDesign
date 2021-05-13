package data

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:mysql@tcp(127.0.0.1:3306)/laravel"

type MyDB struct {
	db *sql.DB
}

func NewDB() (*MyDB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {

		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &MyDB{
		db: db,
	}, nil
}

func (m *MyDB) Register(username, password string) {

}
func (m *MyDB) Login(username, password string) {

}
func (m *MyDB) NewRoom(roomname string) {

}
func (m *MyDB) NewMessage(username, roomname, message string) {

}
func (m *MyDB) History(roomname string) {

}
func (m *MyDB) DelRoom(roomname string) {

}
