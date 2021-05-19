package data

import (
	"database/sql"
	"errors"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

//dsn="username:password@tcp(addr)/database_name"
const dsn = "root:mysql@tcp(127.0.0.1:3306)/chat"

var MYDB *MyDB

type MyDB struct {
	db *sql.DB
}

func GetDB() *MyDB {
	if MYDB != nil {
		return MYDB
	}

	var once sync.Once

	once.Do(func() {
		err := newDB()
		if err != nil {
			log.Println(err)
		}
	})

	return MYDB
}

func newDB() error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	MYDB = &MyDB{
		db: db,
	}
	return nil
}

func (m *MyDB) Register(username, password string) error {
	sqlStr := "insert into user(user_name,nickname,password) values(?,?,?)"

	ret, err := m.db.Exec(sqlStr, username, username, password)
	if err != nil {
		return err
	}

	id, err := ret.LastInsertId()
	if err != nil {
		return err
	}

	log.Println("insert new user success ", id)

	return nil
}

func (m *MyDB) Login(username, password string) error {
	sqlStr := "select password from user where user_name= ?"

	var pwd string
	err := m.db.QueryRow(sqlStr, username).Scan(&pwd)
	if err != nil {
		return err
	}
	if pwd == password {
		log.Println("login success")
		return nil
	}
	return errors.New("username and password not match")
}

func (m *MyDB) NewRoom(roomname string) error {
	sqlStr := "insert into room(room_name) values(?)"

	ret, err := m.db.Exec(sqlStr, roomname)
	if err != nil {
		return err
	}

	id, err := ret.LastInsertId()
	if err != nil {
		return err
	}
	log.Println("insert new room success ", id)
	return nil
}

func (m *MyDB) NewMessage(username, roomname, message string) error {
	sqlStr := "insert into history(message,user_name,room_name) values(?,?,?)"

	ret, err := m.db.Exec(sqlStr, message, username, roomname)
	if err != nil {
		return err
	}

	id, err := ret.LastInsertId()
	if err != nil {
		return err
	}

	log.Println("insert new history success ", id)
	return nil
}

func (m *MyDB) History(roomname string) (*sql.Rows, error) {
	sqlStr := `select message,user_name 
			   from (select message,user_name,created_time from history order by created_time desc limit 0,10)temp
			   order by created_time asc`

	return m.db.Query(sqlStr)
}

func (m *MyDB) DelRoom(roomname string) error {
	sqlStr := "delete from room where room_name=?"

	ret, err := m.db.Exec(sqlStr, roomname)
	if err != nil {
		return err
	}

	n, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("delete room nums ", n)
	return nil
}

func (m *MyDB) AddUser(username, roomname string) error {
	sqlStr := "insert into user_room(user_name,room_name) values(?,?)"

	ret, err := m.db.Exec(sqlStr, username, roomname)
	if err != nil {
		return err
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return err
	}

	log.Println("add user to room success ", id)
	return nil
}

func (m *MyDB) GetRoom(username string) ([]string, error) {
	sqlStr := "select room_name from user_room where user_name=?"

	rows, err := m.db.Query(sqlStr, username)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	var temp string

	for rows.Next() {
		err := rows.Scan(&temp)
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}

	return result, nil
}
