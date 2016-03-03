package simple_oauth2

type Authenticater struct{
}

func (t *Authenticater)AddProfile(profileId int) string{

}
func (t *Authenticater) RemoveProfile(profileId int){

}
func (t *Authenticater) CheckAccessToken(accessToken string) bool{

}
func (t *Authenticater) RefreshToken(clientId string, clientSecret string, refreshToken string){

}
func (t *Authenticater) GetRefreshToken(profileId int) string{

}
