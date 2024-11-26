package auth

import (
	"AABBCCDD/app/db"
	"AABBCCDD/app/db/sqlc"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ignoxx/ohmyskit/kit"
	v "github.com/ignoxx/ohmyskit/validate"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	userSessionName = "user-session"
)

var authSchema = v.Schema{
	"email":    v.Rules(v.Email),
	"password": v.Rules(v.Required),
}

func HandleLoginIndex(kit *kit.Kit) error {
	if kit.Auth().Check() {
		redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")
		return kit.Redirect(http.StatusSeeOther, redirectURL)
	}
	return kit.Render(LoginIndex(LoginIndexPageData{}))
}

func HandleLoginCreate(kit *kit.Kit) error {
	var values LoginFormValues
	errors, ok := v.Request(kit.Request, &values, authSchema)
	if !ok {
		return kit.Render(LoginForm(values, errors))
	}

	user, err := db.Get().FindUserByEmail(kit.Request.Context(), values.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.Add("credentials", "invalid credentials")
			return kit.Render(LoginForm(values, errors))
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(values.Password))
	if err != nil {
		errors.Add("credentials", "invalid credentials")
		return kit.Render(LoginForm(values, errors))
	}

	skipVerify := kit.Getenv("SUPERKIT_AUTH_SKIP_VERIFY", "false")
	if skipVerify != "true" {
		if !user.EmailVerifiedAt.Valid {
			errors.Add("verified", "please verify your email")
			return kit.Render(LoginForm(values, errors))
		}
	}

	sessionExpiryStr := kit.Getenv("SUPERKIT_AUTH_SESSION_EXPIRY_IN_HOURS", "48")
	sessionExpiry, err := strconv.Atoi(sessionExpiryStr)
	if err != nil {
		sessionExpiry = 48
	}

	session, err := db.Get().CreateSession(kit.Request.Context(), sqlc.CreateSessionParams{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(sessionExpiry)),
	})
	if err != nil {
		return err
	}

	sess := kit.GetSession(userSessionName)
	sess.Values["sessionToken"] = session.Token
	sess.Save(kit.Request, kit.Response)
	redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")

	return kit.Redirect(http.StatusSeeOther, redirectURL)
}

func HandleLoginDelete(kit *kit.Kit) error {
	sess := kit.GetSession(userSessionName)
	defer func() {
		sess.Values = map[any]any{}
		sess.Save(kit.Request, kit.Response)
	}()
	err := db.Get().DeleteSessionByToken(kit.Request.Context(), sess.Values["sessionToken"])
	if err != nil {
		return err
	}
	return kit.Redirect(http.StatusSeeOther, "/")
}

func HandleEmailVerify(kit *kit.Kit) error {
	tokenStr := kit.Request.URL.Query().Get("token")
	if len(tokenStr) == 0 {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	token, err := jwt.ParseWithClaims(
		tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SUPERKIT_SECRET")), nil
		}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if !token.Valid {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	user, err := db.Get().FindUserByID(kit.Request.Context(), userID)
	if err != nil {
		return err
	}

	if user.EmailVerifiedAt.Time.After(time.Time{}) {
		return kit.Render(EmailVerificationError("Email already verified"))
	}

	err = db.Get().UpdateUserEmailVerifiedAt(kit.Request.Context(), userID)
	if err != nil {
		return err
	}

	return kit.Redirect(http.StatusSeeOther, "/login")
}

func AuthenticateUser(kit *kit.Kit) (kit.Auth, error) {
	auth := Auth{}
	sess := kit.GetSession(userSessionName)
	token, ok := sess.Values["sessionToken"]
	if !ok {
		return auth, nil
	}

	session, err := db.Get().FindSessionByTokenAndExpiration(kit.Request.Context(), token)
	if err != nil || session.ID == 0 {
		return auth, nil
	}

	return Auth{
		LoggedIn: true,
		UserID:   uint(session.UserID),
		Email:    session.UserEmail,
	}, nil
}
