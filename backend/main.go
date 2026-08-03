package main

import (
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/Anyeling0620/OnlineOJ/backend/router"
)

func main() {
	defer func() {
		err := models.SqlDB.Close()
		if err != nil {
			fmt.Println("close db error:", err)
		}
	}()
	r := router.Router()

	_ = r.Run(":8080")

}
