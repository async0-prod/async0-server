package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/store"
	"github.com/grvbrk/async0_server/internal/utils"
	"github.com/rbcervilla/redisstore/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AdminOAuth interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Callback(w http.ResponseWriter, r *http.Request)
}

type AdminGoogleOauth struct {
	Logger    *log.Logger
	Config    *oauth2.Config
	Store     *redisstore.RedisStore
	UserStore *store.PostgresUserStore
}

func NewAdminGoogleOauth(logger *log.Logger, adminStore *redisstore.RedisStore, userStore *store.PostgresUserStore) (*AdminGoogleOauth, error) {
	return &AdminGoogleOauth{
		Logger: logger,
		Config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID_ADMIN"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET_ADMIN"),
			RedirectURL:  fmt.Sprintf("%s/auth/admin/google/callback", os.Getenv("FRONTEND_URL")),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		Store:     adminStore,
		UserStore: userStore,
	}, nil
}

func (g *AdminGoogleOauth) Login(w http.ResponseWriter, r *http.Request) {
	url := g.Config.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (g *AdminGoogleOauth) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		g.Logger.Println("Error exchanging admin token", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"Error": "Internal Server Error"})
		return
	}

	client := g.Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		g.Logger.Println("Error getting user info", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"Error": "Internal Server Error"})
		return
	}

	defer resp.Body.Close()

	var userInfo struct {
		GoogleID string `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Image    string `json:"picture"`
	}

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		g.Logger.Println("Error decoding user info", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"Error": "Internal Server Error"})
		return
	}

	var userId string
	user, err := g.UserStore.GetUserByGoogleID(userInfo.GoogleID)
	if user == nil || err == sql.ErrNoRows {
		g.Logger.Println("User not found")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"Error": "Unauthorized"})
		return
	}

	if user.Role != "ADMIN" {
		g.Logger.Println("User not admin")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"Error": "Unauthorized"})
		return
	}

	userId = user.ID.String()

	session, _ := g.Store.Get(r, "admin_session")
	session.Values["admin_email"] = userInfo.Email
	session.Values["admin_id"] = userId
	session.Values["admin_image"] = userInfo.Image
	session.Values["admin_name"] = userInfo.Name

	err = session.Save(r, w)
	if err != nil {
		g.Logger.Println("Error saving admin session", err)
	}

	http.Redirect(w, r, fmt.Sprintf("%s/dashboard", os.Getenv("FRONTEND_URL")), http.StatusSeeOther)
}

func (g *AdminGoogleOauth) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := g.Store.Get(r, "session")
	delete(session.Values, "admin_email")
	err := session.Save(r, w)
	if err != nil {
		g.Logger.Println("Error saving admin session", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (g *AdminGoogleOauth) AuthAdmin(w http.ResponseWriter, r *http.Request) {
	session, err := g.Store.Get(r, "admin_session")
	if err != nil || session.IsNew {
		g.Logger.Println("Error getting admin session", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Not Authenticated"})
		return
	}

	adminEmail, emailOk := session.Values["admin_email"].(string)
	adminIDStr, idOk := session.Values["admin_id"].(string)
	adminName, nameOk := session.Values["admin_name"].(string)
	adminImage, imageOk := session.Values["admin_image"].(string)

	if !emailOk || !idOk || !nameOk || !imageOk || adminEmail == "" || adminIDStr == "" || adminName == "" || adminImage == "" {
		g.Logger.Println("Invalid or missing admin data in session")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Not Authenticated"})
		return
	}

	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		g.Logger.Println("Invalid admin ID format in session:", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Not Authenticated"})
		return
	}

	adminInfo := map[string]interface{}{
		"id":    adminID,
		"email": adminEmail,
		"name":  adminName,
		"image": adminImage,
		"role":  "ADMIN",
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": adminInfo})
}
