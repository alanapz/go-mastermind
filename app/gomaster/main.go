package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	games := make(map[string]*Game)

	router := gin.Default()

	router.POST("/api/game/new", func(c *gin.Context) {
		var request StartGameRequest

		err := c.MustBindWith(&request, binding.JSON)

		if err != nil {
			log.Printf("Unable to parse request: %v", err)
			return
		}

		id := Guid()

		game, err := StartGame(GameConfig{NumberOfColours: request.NumberOfColours, NumberOfPositions: request.NumberOfPositions})

		if err != nil {
			log.Printf("Unable to start game: %v", c.AbortWithError(http.StatusBadRequest, err))
			return
		}

		games[id] = game
		c.JSON(200, StartGameResponse{GameId: id})
	})

	router.GET("/api/game/:game", func(c *gin.Context) {

		id := c.Param("game")

		game, ok := games[id]

		if !ok {
			log.Printf("Unable to retrieve game: %v", c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("game not found: %v", id))))
			return
		}

		c.JSON(200, QueryGameResponse{
			NumberOfPositions: game.NumberOfPositions,
			Colours:           game.Colours,
			Complete:          game.Complete,
			Attempts: Map(game.Attempts, func(attempt Attempt) QueryGameAttemptResponse {
				return QueryGameAttemptResponse{
					Positions:             attempt.Positions,
					Complete:              attempt.Complete,
					RightColourRightPlace: attempt.RightColourRightPlace,
					RightColourWrongPlace: attempt.RightColourWrongPlace,
					WrongColour:           attempt.WrongColour,
				}
			}),
		})
	})

	router.POST("/api/game/:game/guess", func(c *gin.Context) {

		id := c.Param("game")

		game, ok := games[id]

		if !ok {
			err := c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("game not found: %v", id)))
			log.Printf("Unable to retrieve game: %v", err)
			return
		}

		var request SubmitGuessRequest

		err := c.MustBindWith(&request, binding.JSON)

		if err != nil {
			log.Printf("Unable to parse request: %v", err)
			return
		}

		_, err = game.SubmitGuess(request.Positions)

		if err != nil {
			log.Printf("Unable to submit guess: %v", c.AbortWithError(http.StatusBadRequest, err))
			return
		}

		c.JSON(200, SubmitGuessResponse{})
	})

	router.StaticFile("/", "./assets/index.html")
	router.StaticFile("/assets/gomaster.css", "./assets/gomaster.css")

	// See: https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-without-context/server.go#L25
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Println("Preparing to start server ...")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
