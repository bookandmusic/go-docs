package common

const (
	APIPrefix = "/api/v1"

	APIDashboard = "/dashboard"

	APILoginUrl    = "/users/login"
	APILogoutUrl   = "/users/logout"
	APIUserInfoUrl = "/users/info"

	APITagUrl       = "/tags"
	APITagDetailUrl = "/tags/:tagId"

	APICategoryUrl       = "/categories"
	APICategoryDetailUrl = "/categories/:categoryId"

	APICollectionUrl        = "/collections"
	APICollectionDetailUrl  = "/collections/:collectionId"
	APICollectionTocListUrl = "/collections/:collectionId/toclist"

	APIArticleUrl       = "/articles"
	APIArticleDetailUrl = "/articles/:articleId"
	APIArticImportleUrl = "/articles/import"
)
