package simple_oauth2

type Oauth2 interface {
	AddProfile(profileId int) string
	RemoveProfile(profileId int)
	CheckAccessToken(accessToken string) bool
	RefreshToken(clientId string, clientSecret string, refreshToken string)
	GetRefreshToken(profileId int) string
}