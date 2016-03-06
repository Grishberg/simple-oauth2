package simple_oauth2

import "github.com/grishberg/simple-oauth2/data"

func NewAuthenticater() Oauth2 {
	var db data.Db
	return &Authenticater{db}
}

type Authenticater struct {
	db data.Db
}

func (t *Authenticater) Init() {
	t.db.Init()
	t.db.Connect()
	return
}

func (t *Authenticater) AddProfile(profileId int64) string {
	return t.db.AddProfile(profileId)
}

func (t *Authenticater) DeleteProfile(profileId int64) {
	t.db.DeleteProfile(profileId)
}
func (t *Authenticater) GetProfile(accessToken string) int64 {
	return t.db.GetProfile(accessToken)
}
func (t *Authenticater) RefreshToken(refreshToken string) (string, int) {
	return t.db.RefreshToken(refreshToken)
}
func (t *Authenticater) GetRefreshToken(profileId int64) string {
	return t.db.GeteRefreshToken(profileId)
}
