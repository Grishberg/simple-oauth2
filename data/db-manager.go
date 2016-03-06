package data

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"strings"
	"math/rand"
	"time"
)

const TOKEN_LENGTH = 16

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type Db struct {
	db        *sql.DB
	db_driver string
	dbName    string
	expire    int32
}

func (t *Db) Init() {
	t.dbName = "./oauth2.db"
	t.expire = 3600
	//sql.Register(t.db_driver, &sqlite3.SQLiteDriver{})
	t.Connect()
	t.createTables()
	t.Close()

}
// for tests
func (t *Db) InitWithName(dbName string, expirePeriod int) {
	sql.Register(t.db_driver, &sqlite3.SQLiteDriver{})
	t.expire = int32(expirePeriod)
	t.dbName = dbName
	t.Connect()
	t.createTables()
	t.Close()
}

func (t*Db) Connect() {
	t.db = t.connect(t.dbName)
}

func (t *Db) connect(dbName string) *sql.DB {
	db, err := sql.Open(t.db_driver, dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func (t *Db) createTables() {
	var err error
	// profiles
	query := `CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY,
		refresh_token varchar(255)
	) ;`
	_, err = t.db.Exec(query)
	t.checkErr(err)

	// access tokens
	query = `CREATE TABLE IF NOT EXISTS tokens (
	access_token string PRIMARY KEY,
	profile_id INTEGER,
	expires_in INTEGER,
	FOREIGN KEY(profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
	UNIQUE (profile_id) ON CONFLICT REPLACE
	) ;`
	_, err = t.db.Exec(query)
	t.checkErr(err)
}

// добавить профиль
func (t *Db) AddProfile(profileId int64) string {

	// insert
	refreshToken := generateToken(TOKEN_LENGTH)

	tx, err := t.db.Begin()
	t.checkErr(err)
	query := `INSERT INTO profiles(id, refresh_token) values(?, ?)`
	stmt, err := tx.Prepare(query)
	t.checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(profileId, refreshToken)
	tx.Commit()

	return refreshToken

}


// добавить профиль
func (t *Db) RefreshToken(refreshToken string) (string, int) {
	profileId := t.GetProfileByRefreshToken(refreshToken)
	if profileId <= 0 {
		return "", 1
	}
	accessToken := t.UpdateAccessToken(profileId)
	return accessToken, 0
}

func (t*Db) UpdateAccessToken(profileId int64) string {
	// update
	refreshToken := generateToken(TOKEN_LENGTH)
	duration := time.Duration(t.expire) * time.Second
	expiresIn := time.Now().Add(duration).Unix()
	tx, err := t.db.Begin()
	t.checkErr(err)
	query := `INSERT INTO tokens (access_token, profile_id, expires_in) values(?, ?, ?)`
	stmt, err := tx.Prepare(query)
	t.checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(refreshToken, profileId, expiresIn)
	tx.Commit()
	return refreshToken

}

func (t*Db) GetProfile(accessToken string) int64 {
	var profileId int64
	var expiresIn int64
	query := "SELECT profile_id, expires_in FROM tokens WHERE access_token =?;"
	t.db.QueryRow(query, accessToken).Scan(&profileId, &expiresIn)
	if time.Now().Unix() < expiresIn {
		return profileId
	}
	return -1
}

func (t*Db) GeteRefreshToken(profileId int64) string {
	var refreshToken string
	query := "SELECT refresh_token FROM profiles WHERE id =?;"
	err := t.db.QueryRow(query, profileId).Scan(&refreshToken)
	t.checkErr(err)

	return refreshToken
}

func (t*Db) GetProfileByRefreshToken(refreshToken string) int64 {
	var profileId int64
	query := "SELECT id FROM profiles WHERE refresh_token =?;"
	err := t.db.QueryRow(query, refreshToken).Scan(&profileId)
	t.checkErr(err)

	return profileId
}

func (t *Db)DeleteProfile(profileId int64) {
	//tx, _ := t.db.Begin()
	query := "DELETE FROM profiles WHERE profile_id = ?;"
	_, err := t.db.Exec(query, profileId)
	t.checkErr(err)
	//tx.Commit()
}

func (t *Db) Close() {
	t.db.Close()
}

func (t *Db)checkErr(err error) bool {
	if err != nil {
		if ( strings.Contains(err.Error(), "UNIQUE constraint failed")) {
			return true
		}
		//log.Fatal(err)
		panic(err)
	}
	return false
}
