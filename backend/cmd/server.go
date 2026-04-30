package main

import (
	"log"
	"net/http"
	"time"

	"forum/backend/internal/database"
	"forum/backend/internal/handlers"
	"forum/backend/internal/middleware"
	"forum/backend/internal/services"
)

type Server struct {
	mux  *http.ServeMux
	addr string
}

func NewServer(db *database.SQLiteStore) *Server {
	// Services
	userService := &services.UserService{Db: db}
	sessionService := &services.SessionService{Db: db, SessionDuration: 24 * time.Hour}
	postService := &services.PostService{Db: db}
	commentService := &services.CommentService{Db: db}
	reactionService := &services.ReactionService{Db: db}

	// Auth middleware
	auth := middleware.AuthMiddlewareFactory(sessionService)

	// Handlers
	authHandler := &handlers.AuthHandler{
		UserService:    userService,
		SessionService: sessionService,
	}
	userHandler := &handlers.UserHandler{
		UserService: userService,
	}
	contentHandler := &handlers.ContentHandler{
		PostService:     postService,
		CommentService:  commentService,
		ReactionService: reactionService,
	}
	categoryHandler := &handlers.CategoryHandler{
		PostService: postService,
	}

	mux := http.NewServeMux()

	// Static file serving for user-uploaded images (outside /api)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Auth
	mux.HandleFunc("POST /api/auth/signup", handlers.ErrorHandlerAdapter(authHandler.HandleSignUp))
	mux.HandleFunc("POST /api/auth/login", handlers.ErrorHandlerAdapter(authHandler.HandleLogin))
	mux.HandleFunc("POST /api/auth/logout", auth(handlers.ErrorHandlerAdapter(authHandler.HandleLogout)))

	// Categories
	mux.HandleFunc("GET /api/categories", handlers.ErrorHandlerAdapter(categoryHandler.HandleGetCategories))
	mux.HandleFunc("GET /api/categories/{id}", handlers.ErrorHandlerAdapter(categoryHandler.HandleGetCategoryPosts))

	// Users – specific paths before wildcard
	mux.HandleFunc("GET /api/users", handlers.ErrorHandlerAdapter(userHandler.HandleGetUsers))
	mux.HandleFunc("GET /api/users/me", auth(handlers.ErrorHandlerAdapter(userHandler.HandleGetMe)))
	mux.HandleFunc("PATCH /api/users/me", auth(handlers.ErrorHandlerAdapter(userHandler.HandleUpdateMe)))
	mux.HandleFunc("DELETE /api/users/me", auth(handlers.ErrorHandlerAdapter(userHandler.HandleDeleteMe)))
	mux.HandleFunc("GET /api/users/me/posts", auth(handlers.ErrorHandlerAdapter(userHandler.HandleGetMePosts)))
	mux.HandleFunc("GET /api/users/me/comments", auth(handlers.ErrorHandlerAdapter(userHandler.HandleGetMeComments)))
	mux.HandleFunc("GET /api/users/me/likes", auth(handlers.ErrorHandlerAdapter(userHandler.HandleGetMeLikes)))
	mux.HandleFunc("GET /api/users/{id}", handlers.ErrorHandlerAdapter(userHandler.HandleGetUser))

	// Posts
	mux.HandleFunc("GET /api/posts", handlers.ErrorHandlerAdapter(contentHandler.HandleGetPosts))
	mux.HandleFunc("POST /api/posts", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleCreatePost)))
	mux.HandleFunc("GET /api/posts/{id}", handlers.ErrorHandlerAdapter(contentHandler.HandleGetPost))
	mux.HandleFunc("PATCH /api/posts/{id}", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleUpdatePost)))
	mux.HandleFunc("DELETE /api/posts/{id}", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleDeletePost)))

	// Post reactions
	mux.HandleFunc("PUT /api/posts/{id}/reaction", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleReactToPost)))
	mux.HandleFunc("DELETE /api/posts/{id}/reaction", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleRemovePostReaction)))

	// Comments
	mux.HandleFunc("GET /api/posts/{id}/comments", handlers.ErrorHandlerAdapter(contentHandler.HandleGetPostComments))
	mux.HandleFunc("POST /api/posts/{id}/comments", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleCreateComment)))
	mux.HandleFunc("PATCH /api/comments/{id}", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleUpdateComment)))
	mux.HandleFunc("DELETE /api/comments/{id}", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleDeleteComment)))

	// Comment reactions
	mux.HandleFunc("PUT /api/comments/{id}/reaction", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleReactToComment)))
	mux.HandleFunc("DELETE /api/comments/{id}/reaction", auth(handlers.ErrorHandlerAdapter(contentHandler.HandleRemoveCommentReaction)))

	return &Server{mux: mux, addr: ":8080"}
}

func (s *Server) ListenAndServe() error {
	log.Printf("Server listening on %s", s.addr)
	return http.ListenAndServe(s.addr, middleware.RequestLoggerMiddleware(s.mux))
}

func main() {
	db, err := database.NewSQLiteStore("sqlite3", "./forum.db", "./internal/database/schema.sql")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Fatal(NewServer(db).ListenAndServe())
}
