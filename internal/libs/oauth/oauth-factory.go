package oauth

// import (
// )

// type OauthClientFactory struct {
// 	clientMap map[user.ProviderType]OauthClient
// }

// func NewFactory() *OauthClientFactory {
// 	return &OauthClientFactory{
// 		clientMap: map[user.ProviderType]OauthClient{
// 			user.ProviderKAKAO:  newKakaoClient(),
// 			user.ProviderNAVER:  newNaverClient(),
// 			user.ProviderGOOGLE: newGoogleClient(),
// 		},
// 	}
// }

// func (factory *OauthClientFactory) GetClient(providerType user.ProviderType) OauthClient {
// 	return factory.clientMap[providerType]
// }
