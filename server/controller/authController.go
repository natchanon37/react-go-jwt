package controller

import (
	"react-go-jwt/server/database"
	"react-go-jwt/server/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Hello(c *fiber.Ctx) error  {
	return c.SendString("Hello From fiber")
}


func Register (c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil{
		return err
	}

	/*Create User and hash password*/
	password , _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name: data["name"],
		Email: data["email"],
		Password: password,
	}

	//Insert data to database
	database.DB.Create(&user)
	return c.JSON(user)
}

var secretkey = []byte("secret")
func Login(c *fiber.Ctx) error {

	var data map[string]string
	var user models.User

	/* Get data from request and map in key-value */
	if err := c.BodyParser(&data); err != nil{
		return err
	}

	/*Check submit email in database and store in user variable*/
	database.DB.Where("email = ?",data["email"]).First(&user)
	if user.Id == 0{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message":"user not found",
		})
	}

	/*Check password*/

	if err := bcrypt.CompareHashAndPassword(user.Password,[]byte(data["password"])); err != nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message":"incorrect password",
		})
	}

	//JWT
	clams := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
		Issuer: strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour*24).Unix(), // 1day
	})

	token, err := clams.SignedString(secretkey)

	if err != nil{
		return err
	}

	/*Store token in cookie*/
	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)


	return c.JSON(fiber.Map{
		"message":"success",
	})
}

/*
	@desc Retrieve login user using only cookie
*/
func User(c *fiber.Ctx) error  {

	// 1. Retrieve a cookie that got when login
	cookie := c.Cookies("jwt")
	var user models.User

	//2. Retrieve a token from cookie
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{},func(t *jwt.Token) (interface{}, error) {
		return secretkey, nil
	})

	if err != nil{
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message":"unauthenticated",
		})
	}

	/*Get the user id from the token using Claims*/
	claims := token.Claims.(*jwt.StandardClaims) 
	/* We want a Issuer property which is a user Id so we need 
	to convert type "Claims" -> "StandardClaims"*/


	//Query login user by claims.Issuer
	database.DB.Where("id = ?",claims.Issuer).First(&user)

	return c.JSON(user)	
}

func Logout(c *fiber.Ctx) error{
	//Remove cookie by set the expire to past
	cookie := fiber.Cookie{
		 Name: "jwt",
		 Value: "",
		 MaxAge: -1,
		 Expires: time.Now().Add(-time.Hour),
		 HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"success",
	})

}