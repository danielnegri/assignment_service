package web

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/surajjain36/assignment_service/infra"
	"github.com/surajjain36/assignment_service/misc"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
)

// Service HTTP server info
type Service struct {
	shutdownChan chan bool
	domain       string
	router       *gin.Engine
	wg           sync.WaitGroup
	mdb          *infra.Mongo
	AppName      string
	Version      string
	BuildTime    string
}

// NewService Creates a new web service
func NewService(conf *misc.Config) (*Service, error) {
	mdb, err := infra.NewMongo(&conf.Mongo)
	if err != nil {
		log.WithError(err).Error("Failed to connect to MongoDB")
		return nil, err
	}

	s := &Service{
		router:       gin.New(),
		mdb:          mdb,
		shutdownChan: make(chan bool),
		domain:       conf.HTTP.Domain,
	}

	s.router.Use(gin.Logger())
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS", "POST"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	s.router.GET("/", s.index)
	s.router.GET("/ping", s.ping)

	v1 := s.router.Group("/v1")
	{
		v1.POST("/assignment", s.createAssignment)
		v1.GET("/assignment/:id", s.getAssignment)
		v1.GET("/search/assignment", s.searchAssignmentByTags)
	}

	url := ginSwagger.URL(fmt.Sprintf("%s%s%s", "http://", conf.HTTP.Domain, "/docs/swagger.json"))
	s.router.StaticFile("docs/swagger.json", "./docs/swagger.json")
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return s, nil
}

// Start the web service
func (s *Service) Start(address string) error {
	return s.router.Run(address)
}

// Close all threads and free up resources
func (s *Service) Close() {
	close(s.shutdownChan)
	s.wg.Wait()
}
