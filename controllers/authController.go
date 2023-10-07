package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shohan-joarder/go-jwt-starter/database"
	"github.com/shohan-joarder/go-jwt-starter/models"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate
// var DB *sql.DB

func init()  {
	validate = validator.New()
	database.InitDB()
}

type CustomValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type CustomValidationErrors struct {
	Errors []CustomValidationError `json:"errors"`
}

func Register(c *gin.Context) {
	// var newUser models.UserRegistrations
	var userValidation models.UserValidations
	
	if err := c.ShouldBind(&userValidation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user input
	if err := validate.Struct(userValidation); err != nil {
		fmt.Println( err.(validator.ValidationErrors))
		validationErrors := []CustomValidationError{}
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			validationErrors = append(validationErrors, CustomValidationError{
				Field: field,
				Error: fmt.Sprintf("Invalid value for '%s' field. Tag: '%s'", field, tag),
			})
		}
		c.JSON(http.StatusBadRequest, CustomValidationErrors{Errors: validationErrors})
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(userValidation.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error generating password"})
		return
	}

	insertQuery := "INSERT INTO users (first_name,last_name,email,phone,password) VALUES(?,?,?,?,?)"
	
	_, err = database.DB.Exec(insertQuery, &userValidation.FirstName, &userValidation.LastName, &userValidation.Email,&userValidation.Phone,&password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database", "err": err.Error()})
		return
	}
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error registered user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}



func Login(c *gin.Context) {
	var validUser models.Login

	if err := c.ShouldBind(&validUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user input
	if err := validate.Struct(validUser); err != nil {
		fmt.Println( err.(validator.ValidationErrors))
		validationErrors := []CustomValidationError{}
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			validationErrors = append(validationErrors, CustomValidationError{
				Field: field,
				Error: fmt.Sprintf("Invalid value for '%s' field. Tag: '%s'", field, tag),
			})
		}
		c.JSON(http.StatusBadRequest, CustomValidationErrors{Errors: validationErrors})
		return
	}

	// Query the database for the user's hashed password
	var hashedPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE email = ?", validUser.Email).Scan(&hashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the provided password with the hashed password from the database
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(validUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	 
	// jwt_secret:=os.Getenv("JWT_KEY")

	// fmt.Println(secretKey,token,tokenString)

	token, err :=CreateToken(validUser.Email)
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to generate token"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"message": "Login successfully","token":token})
}

func CreateToken(email string) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
