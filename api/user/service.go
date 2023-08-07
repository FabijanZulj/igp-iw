package user

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/igp-iw/database"
	"github.com/igp-iw/notifications"
	"github.com/igp-iw/util"
	"golang.org/x/crypto/bcrypt"
)

// UserService is a service that provides ways of registering, verifying and logging in of users
type UserService interface {
	RegisterUser(*RegisterRequest, string) error
	VerifyEmail(string, string) error
	Login(*LoginRequest) (string, error)
}

type userService struct {
	query         database.Querier
	db            *sql.DB
	notifications notifications.NotificationPublisher
	jwtSecret     string
}

// NewUserService creates a new instance of UserService
func NewUserService(queries database.Querier, db *sql.DB,
	notifications notifications.NotificationPublisher, jwtSecret string,
) *userService {
	return &userService{
		query:         queries,
		db:            db,
		notifications: notifications,
		jwtSecret:     jwtSecret,
	}
}

// RegisterUser hashes the password and creates the entry in the database of the user that is registering
// and also creates the verification data (verification code that is randomly generated 20 length string)
// It also sends a email notification - verification email. That contains the url for verification
func (ur *userService) RegisterUser(rr *RegisterRequest, baseUrl string) error {
	ctx := context.Background()
	tx, err := ur.db.Begin()
	if err != nil {
		log.Println("Error while creating transaction for registering" + rr.Email)
		return err
	}
	defer tx.Rollback()

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(rr.Password), 14)
	if err != nil {
		log.Println("Error while creating password hash, for email" + rr.Email)
		return err
	}

	usr, err := ur.query.CreateUser(ctx, tx, database.CreateUserParams{
		Email:      rr.Email,
		Password:   string(hashedPw),
		Isverified: false,
	})
	if err != nil {
		log.Println("Error while creating user data" + err.Error())
		return err
	}

	verifyCode := util.RandomString(20)
	_, err = ur.query.CreateVerifyData(ctx, tx, database.CreateVerifyDataParams{
		Userid: usr.ID,
		Code:   verifyCode,
	})

	if err != nil {
		log.Println("Error while creating verify data" + err.Error())
		return err
	}
	err = ur.notifications.Publish(notifications.NotificationData{
		RawData: map[string]any{
			"redirectUrl": util.BuildVerificationUrl(baseUrl, rr.Email, verifyCode),
			"template":    "email_verification.html",
			"subject":     "Verify your email!",
		},
		NotificationType: "EMAIL_NOTIFICATION",
		Initiator:        rr.Email,
		Target:           rr.Email,
	})

	if err != nil {
		log.Println("Error while sending email confirmation notification")
		return err
	}

	tx.Commit()
	return nil
}

// VerifyEmail compares the provided verification code and the needed code in the verificationData table.
// Sets the user status to verified and sends a welcome email if the verification code is correct.
func (ur *userService) VerifyEmail(email, code string) error {
	ctx := context.Background()
	vd, err := ur.query.GetVerificationCode(ctx, ur.db, email)
	if err != nil {
		log.Println("Cannot get code for verification")
		return err
	}

	if vd.Code != code {
		return errors.New("cannot verify user with given code")
	}

	ur.query.UpdateUser(ctx, ur.db, database.UpdateUserParams{
		Isverified: true,
		Email:      email,
	})

	err = ur.notifications.Publish(notifications.NotificationData{
		RawData: map[string]any{
			"template": "welcome_email.html",
			"subject":  "Welcome!",
		},
		NotificationType: "EMAIL_NOTIFICATION",
		Initiator:        email,
		Target:           email,
	})

	if err != nil {
		log.Println("Error while sending email confirmation notification")
		return err
	}

	return nil
}

// Login compares the provided login password with the hashed pasword in the database and
// if the comparation is currect returns a JWT token
func (ur *userService) Login(lr *LoginRequest) (string, error) {
	ctx := context.Background()
	usr, err := ur.query.GetUserByEmail(ctx, ur.db, lr.Email)
	if err != nil {
		log.Println("Error while fetching user with given email")
		return "", err
	}

	failure := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(lr.Password))
	if usr.Isverified && failure == nil {
		claims := jwt.MapClaims{
			"email": usr.Email,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(ur.jwtSecret))
		if err != nil {
			log.Println("Cannot sign token")
			return "", err
		}

		return t, nil
	}

	return "", errors.New("password not correct")
}
