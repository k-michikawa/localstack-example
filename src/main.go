package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var svc *sqs.SQS

type Config struct {
	ListenPort   string `split_words:"true"`
	AwsRegion    string `split_words:"true"`
	SqsQueueName string `split_words:"true"`
	SqsEndpoint  string `split_words:"true"`
}

type EnhancedContext struct {
	echo.Context
	Config
}

type SendMessageInput struct {
	Message string `json:"message"`
}

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world")
}

func SendMessage(c echo.Context) error {
	ec := c.(*EnhancedContext)
	req := new(SendMessageInput)
	if err := c.Bind(req); err != nil {
		log.Print(err)
		return c.String(http.StatusInternalServerError, "")
	}

	url := strings.Join([]string{ec.Config.SqsEndpoint, "localstack", ec.Config.SqsQueueName}, "/")
	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(req.Message),
		QueueUrl:     aws.String(url),
		DelaySeconds: aws.Int64(1),
	}

	res, err := svc.SendMessage(params)
	if err != nil {
		log.Print(err)
		return c.String(http.StatusInternalServerError, "")
	}
	return c.String(http.StatusOK, *res.MessageId)
}

// 手抜き
func enhanceContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ec := &EnhancedContext{Context: c}
		return next(ec)
	}
}

func configMiddleware(config Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ec := c.(*EnhancedContext)
			ec.Config = config
			return next(ec)
		}
	}
}

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err.Error())
	}
	s := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(config.SqsEndpoint),
		Region:   aws.String(config.AwsRegion),
	}))

	fmt.Print(config)

	svc = sqs.New(s)
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(enhanceContextMiddleware)
	e.Use(configMiddleware(config))

	// Routes
	e.GET("/", Hello)
	e.POST("/send-message", SendMessage)

	e.Logger.Fatal(e.Start(config.ListenPort))
}
