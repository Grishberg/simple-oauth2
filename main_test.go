package simple_oauth2

import (
	"testing"
	"os"
	"github.com/grishberg/simple-oauth2/data"
	"log"
	"time"
)

const (
	TEST_BASE_NAME = "./test.db"
	TEST_PROFILE = 1
)

func TestDb(t *testing.T) {
	os.Remove(TEST_BASE_NAME)
	var db data.Db
	log.Println("init db")
	db.InitWithName(TEST_BASE_NAME, 1)
	db.Connect()
	refreshToken := db.AddProfile(TEST_PROFILE)
	accessToken, err := db.RefreshToken(refreshToken)
	if ( err > 0) {
		t.Error("err != 0: ")
	}
	profileId := db.GetProfile(accessToken)

	if ( profileId != TEST_PROFILE) {
		t.Error("profile != TEST_PROFILE: ", profileId)
	}

	db.Close()
}
// тестирование устаревания токена
func TestCreatingProfile(t *testing.T) {
	os.Remove(TEST_BASE_NAME)
	var db data.Db
	log.Println("init db")
	db.InitWithName(TEST_BASE_NAME, 1)
	db.Connect()

	var auth Oauth2
	auth = NewAuthenticater(db)
	refreshToken := auth.AddProfile(TEST_PROFILE)
	accessToken, err := auth.RefreshToken(refreshToken)
	if ( err > 0) {
		t.Error("err != 0: ")
	}
	time.Sleep(2 * time.Second)
	profileId := auth.GetProfile(accessToken)

	if ( profileId != -1) {
		t.Error("profile != TEST_PROFILE: ", profileId)
	}

	db.Close()
}


// тестирование устаревания токена и успешное обновление
func TestRefreshToken(t *testing.T) {
	os.Remove(TEST_BASE_NAME)
	var db data.Db
	log.Println("init db")
	db.InitWithName(TEST_BASE_NAME, 1)
	db.Connect()

	var auth Oauth2
	auth = NewAuthenticater(db)
	refreshToken := auth.AddProfile(TEST_PROFILE)
	accessToken, err := auth.RefreshToken(refreshToken)
	if ( err > 0) {
		t.Error("err != 0: ")
	}
	time.Sleep(2 * time.Second)
	profileId := auth.GetProfile(accessToken)

	if ( profileId != -1) {
		t.Error("profile != TEST_PROFILE: ", profileId)
	}
	accessToken, err = auth.RefreshToken(refreshToken)

	if (err != 0 || auth.GetProfile(accessToken) != TEST_PROFILE) {
		t.Error("Fail, wrong access token")
	}
	db.Close()
}

