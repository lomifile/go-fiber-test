package controllers

import (
	"github.com/auth-api/database"
	"github.com/auth-api/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
  "github.com/golang-jwt/jwt/v4"
  "strconv"
  "time"
)

const SecretKey = "secret"

func Register (c *fiber.Ctx) error {
  var data map[string]string
  if err := c.BodyParser(&data); err != nil {
    return err
  }
  
  password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

  user := models.User {
    Name: data["name"],
    Email: data["email"],
    Password: password ,
  }

  database.DB.Create(&user)

  return c.JSON(user)
}

func Login (c *fiber.Ctx) error {
  var data map[string]string

  if err := c.BodyParser(&data); err != nil {
    return err
  }

  var user models.User

  database.DB.Where("email = ?", data["email"]).First(&user)

  if user.Id == 0 {
    c.Status(fiber.StatusNotFound)
    return c.JSON(fiber.Map{
      "message":"User not found",
    })
  }

  if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
    c.Status(fiber.StatusBadRequest)
    return c.JSON(fiber.Map {
      "message":"Password not correct",
    })
  }
  
  claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
    Issuer: strconv.Itoa(int(user.Id)),
    ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
  })
  
  token, err := claims.SignedString([]byte(SecretKey))
  if err != nil {
    c.Status(fiber.StatusInternalServerError)
    return c.JSON(fiber.Map{
      "message":"Not logged in",
    })
  }
  
  cookie := fiber.Cookie{
    Name: "JWT",
    Value: token,
    Expires: time.Now().Add(time.Hour * 24),
    HTTPOnly: true,
  }
  
  c.Cookie(&cookie)

  return c.JSON(fiber.Map{
    "status":fiber.StatusOK,
    "message":"Login success",
    "data":user,
  })
}

func User (c *fiber.Ctx) error {
  cookie := c.Cookies("JWT")

  token , err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
    return []byte(SecretKey), nil
  })

  if err != nil {
    c.Status(fiber.StatusUnauthorized)
    return c.JSON(fiber.Map{
      "message":"unauthorized",
    })
  }

  claims := token.Claims.(*jwt.StandardClaims)
  
  var user models.User
  
  database.DB.Where("id = ?", claims.Issuer).First(&user)

  return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
  cookie := fiber.Cookie{
    Name: "JWT",
    Value: "",
    Expires: time.Now().Add(-time.Hour),
    HTTPOnly: true,
  }

  c.Cookie(&cookie)

  return c.JSON(fiber.Map{
    "status":fiber.StatusOK,
    "message":"Logout success",
  })
}

