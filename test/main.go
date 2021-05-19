package main

import (
	"Shaw/goWeb/chatRoom/data"
	"fmt"
)

func main() {
	db := data.GetDB()

	// err := db.Register("beijing", "23679")
	// err := db.Login("xiaoming", "1234596")
	// err := db.NewRoom("main")
	// err := db.DelRoom("main")

	// err := db.NewMessage("xiaoming", "main", "zhende")
	// err = db.NewMessage("lihua", "main", "今天")
	// err = db.NewMessage("lihua", "main", "真没")
	// err = db.NewMessage("xiaoming", "main", "干劲")
	// err = db.NewMessage("lihua", "main", "不是吗")

	// rows, err := db.History("main")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// var msg, name string
	// for rows.Next() {
	// 	err := rows.Scan(&name, &msg)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Println(msg, name)
	// 	}
	// }

	// err := db.NewRoom("科技")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = db.AddUser("lihua", "科技")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	ret, err := db.GetRoom("lihua")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret)

	fmt.Println("operation success")
}
