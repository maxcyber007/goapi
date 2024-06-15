package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Binding from JSON
type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

type Tbl_User struct {
	gorm.Model
	Username string
	Password string
	Fullname string
	Avatar   string
}

func main() {
	dsn := "root:1234@tcp(127.0.0.1:3306)/db_gojwt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Tbl_User{})

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", func(c *gin.Context) {
		var json Register
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Check User Exists
		var userExist Tbl_User
		db.Where("username = ?", json.Username).First(&userExist)
		if userExist.ID > 0 {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Create User Failed"})
		}

		//Create User
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 0)
		user := Tbl_User{Username: json.Username, Password: string(encryptedPassword), Fullname: json.Fullname, Avatar: json.Avatar}
		db.Create(&user) // pass pointer of data to Create
		if user.ID > 0 {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Create User Success", "UserID": user.ID})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Create User Failed"})
		}

	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
