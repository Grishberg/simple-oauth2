package data

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"strings"
	"log"
	"strconv"
	"fmt"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type Db struct {
	db        *sql.DB
	db_driver string
	dbName    string
	expire    int
}

func (t *Db) Init() {
	t.dbName = "./oauth2.db"
	t.expire = 60
	sql.Register(t.db_driver, &sqlite3.SQLiteDriver{})
	t.Connect()
	t.createTables()
	t.Close()

}
// for tests
func (t *Db) InitWithName(dbName string, expirePeriod int) {
	sql.Register(t.db_driver, &sqlite3.SQLiteDriver{})
	t.expire = expirePeriod
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
		refresh_token varchar(255),
	) ;`
	_, err = t.db.Exec(query)
	t.checkErr(err)

	// access tokens
	query = `CREATE TABLE IF NOT EXISTS tokens (
	accessToken string PRIMARY KEY,
	profileId INTEGER,
	expires_in INTEGER,
	FOREIGN KEY(profileId) REFERENCES profiles(id),
	UNIQUE (profileId) ON CONFLICT REPLACE
	) ;`
	_, err = t.db.Exec(query)
	t.checkErr(err)
}

// добавить профиль
func (t *Db) AddProfile(profileId int) string {

	// insert
	refreshToken := generateToken(8)

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

func (t*Db) UpdateRefreshToken(profileId int) string {
	// update
	refreshToken := generateToken(8)
	expiresIn := time.Now().Add(t.expire * time.Minute).Unix()
	tx, err := t.db.Begin()
	t.checkErr(err)
	query := `INSERT INTO tokens (accessToken, profileId, expires_in) values(?, ?, ?)`
	stmt, err := tx.Prepare(query)
	t.checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(refreshToken, profileId, expiresIn)
	tx.Commit()
	return refreshToken

}

func (t*Db) GenerateAccessToken(profileId int) int64 {

	// update
	accessToken := generateToken(8)

	tx, err := t.db.Begin()
	t.checkErr(err)
	query := `UPDATE profiles SET refresh_token = ? WHERE id = ?;`
	stmt, err := tx.Prepare(query)
	t.checkErr(err)
	defer stmt.Close()

	var res sql.Result
	res, err = stmt.Exec(accessToken, profileId)
	tx.Commit()
	id, err := res.LastInsertId()
	t.checkErr(err)
	return accessToken
}

func (t*Db) GetProfile(accessToken string) int64 {
	var profileId int64
	var expiresIn int64
	query := "SELECT profileId, expires_in FROM tokens WHERE accessToken =?;"
	err := t.db.QueryRow(query, accessToken).Scan(&profileId, &expiresIn)
	t.checkErr(err)

	if time.Now().Unix() > expiresIn {
		return profileId
	}

	return -1
}

func (t *Db)DeleteProfile(int profileId) {
	//tx, _ := t.db.Begin()
	query := "DELETE FROM tasks WHERE profileId = ?;"
	_, err := t.db.Exec(query)
	t.checkErr(err)
	//tx.Commit()
}

func (t *Db)DeleteWorkFlows() {
	//tx, _ := t.db.Begin()
	//defer tx.Commit()
	query := "DELETE FROM workflows;"
	_, err := t.db.Exec(query)
	t.checkErr(err)

}


// workflow
func (t *Db) AddWorkflow(work WorkFlowHtml) (int64, bool) {
	parentId := t.getTaskByKey(t.db, work.TaskKey)
	//log.Println("task key: ", work.TaskKey, parentId)
	userId := t.getUserByName(t.db, work.User)

	t.updateProjectByKey(t.db, work.ProjectKey, work.ProjectTitle, work.ProjectUrl)
	tx, err := t.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	query := `INSERT INTO workflows(
		project_id  ,
		assignee_id ,
		fact 		,
		due		,
		projectKey	,
		taskKey		,
		subTaskKey
	)
	 values(?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(query)
	t.checkErr(err)
	defer stmt.Close()

	var res sql.Result
	res, err = stmt.Exec(
		parentId,
		userId,
		work.Spent,
		work.Due.Unix(),
		work.ProjectKey,
		work.SubTaskTitle,
		work.TaskKey)
	tx.Commit()
	if t.checkErr(err) {
		return -1, parentId > 0
	}

	id, err := res.LastInsertId()
	t.checkErr(err)
	return id, parentId > 0
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
