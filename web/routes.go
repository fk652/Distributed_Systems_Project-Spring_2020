// routes.go

package main

func initializeRoutes() {

	router.Use(setUserStatus())

	router.GET("/", showIndexPage)

	userRoutes := router.Group("/u")
	{
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)

		userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)

		userRoutes.GET("/logout", ensureLoggedIn(), logout)

		userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)

		userRoutes.POST("/register", ensureNotLoggedIn(), register)

		userRoutes.GET("/followed", ensureLoggedIn(), getFollowPage)

		userRoutes.POST("/follow", ensureLoggedIn(), performFollow)

		userRoutes.POST("/unfollow", ensureLoggedIn(), performUnfollow)
	}

	articleRoutes := router.Group("/article")
	{
		articleRoutes.GET("/welcome", ensureLoggedIn(), getWelcomePage)

		articleRoutes.GET("/home", ensureLoggedIn(), getHomePage)

		articleRoutes.GET("/search", ensureLoggedIn(), showSearchArticlesPage)

		articleRoutes.POST("/search", ensureLoggedIn(), performSearchArticles)

		articleRoutes.GET("/articleID/:article_id", ensureLoggedIn(), getArticle)

		articleRoutes.POST("/articleID/:article_id", ensureLoggedIn(), getArticle)

		articleRoutes.GET("/userArticles/:username", ensureLoggedIn(), getUserArticle)

		articleRoutes.POST("/userArticles/:username", ensureLoggedIn(), getUserArticle)

		articleRoutes.GET("/myArticles", ensureLoggedIn(), getMyArticles)

		articleRoutes.GET("/create", ensureLoggedIn(), showArticleCreationPage)

		articleRoutes.POST("/create", ensureLoggedIn(), createArticle)
	}
}
