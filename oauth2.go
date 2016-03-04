package simple_oauth2

type Oauth2 interface {
	Init()
	AddProfile(profileId int64) string
	DeleteProfile(profileId int64)
	GetProfile(accessToken string) int64
	RefreshToken(refreshToken string)  (string, int)
	GetRefreshToken(profileId int64) string
}