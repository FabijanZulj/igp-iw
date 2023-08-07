package user

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// RegisterRequest contains data for registering
type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginRequest contains data for logging in
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthHandler contains the handler methods for all endpoints related to auth and registering
type AuthHandler struct {
	userService UserService
	validate    *validator.Validate
}

// NewAuthHandler creates a new instance of AuthHandler with provided validator instance and user userService
// note - validator instance should be used as a singleton as it is optimised for it with caching etc.
func NewAuthHandler(userService UserService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		validate:    validate,
	}
}

// Verify is the handler for verifying using the verification code sent over email
// It reads the email and the code from the route parameters
//
// Returns:
// status 400 - if the code or email is not provided in the route
// status 400 - if the code doesnt match the code in the database
// status 200 - if the account is verified
func (uh *AuthHandler) Verify() fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Params("email")
		code := c.Params("code")

		if email == "" || code == "" {
			log.Println("Provided email or code are not valid" + email + code)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err := uh.userService.VerifyEmail(email, code)
		if err != nil {
			log.Println("Could not verify email" + err.Error() + email + code)
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.SendString("Succesfully verified!")
	}
}

// Register is the handler for POST /register request.
// It parses the registration json and expects it to be valid as defined in RegisterRequest
// Returns:
// status 422 - if the json cannot be processed or it is not valid
// status 400 - if the registration process failed
// status 200 - if the registration was succesfull
func (uh *AuthHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reg := &RegisterRequest{}
		err := c.BodyParser(reg)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(
				fiber.Map{
					"errors": err.Error(),
				},
			)
		}
		err = uh.validate.Struct(reg)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(
				fiber.Map{
					"errors": err.Error(),
				},
			)
		}
		err = uh.userService.RegisterUser(reg, c.BaseURL())

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"errors": err.Error(),
				},
			)
		}
		reg.Password = ""
		return c.Status(fiber.StatusOK).JSON(reg)
	}
}

// Register is the handler for POST /login request.
// It parses and validates the json as LoginRequest and returns JWT if the information is correct
// Returns:
// status 422 - if the json is unprocessable or invalid
// status 400 - if the login process was unsuccesful (wrong password etc.)
// status 200 and JWT in the token field - if the login was sucessfull
func (uh *AuthHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logi := &LoginRequest{}
		err := c.BodyParser(logi)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(
				fiber.Map{
					"errors": err.Error(),
				},
			)
		}
		err = uh.validate.Struct(logi)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(
				fiber.Map{
					"errors": err.Error(),
				},
			)
		}

		token, err := uh.userService.Login(logi)
		if err != nil {
			log.Println("Could not login")
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"errors": err.Error(),
				},
			)
		}

		return c.JSON(fiber.Map{"token": token})
	}
}
