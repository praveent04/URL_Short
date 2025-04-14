package routes

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/praveent04/url-short/database"
	"github.com/praveent04/url-short/helpers"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL         string 			`json:"url"`
	CustomShort string			`json:"short"`
	Expiry      time.Duration 	`json:"expiry"`
}

type response struct {
	URL         	string 			 `json:"url"`
	CustomShort 	string			 `json:"short"`
	Expiry      	time.Duration 	 `json:"expiry"`
	XRateRemaining 	int 			 `json:"rate_limit"`
	XRateLimitRest 	time.Duration 	 `json:"rate_limit_rest"`
}

func ShortenURL(c *fiber.Ctx) error{
	body := new(request)
	
	if err := c.BodyParser(&body); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse JSON"})
	}

	//rate limiter
	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(databaseCtx,c.IP()).Result()
	if err ==  redis.Nil{
		_ = r2.Set(database.Ctx, c.IP, os.Getenv("API_QUOTA "), 30*60*time.Second).Err()
	} else{
		val, _ = r2.Get(database.Ctx,c.IP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0{
			limit,_ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":"Rate limit exceeded","rate_limit_rest": limit/ time.Nanosecond/time.Minute})
		}
	}


	// check if the input sent by user is actual url or not

	if !govalidator.IsURL(body.URL){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}

	//check for domain error
	
	if !helpers.RemoveDomainError(body.URL){
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{})
	}
	//enforce https, SSL

	body.URL =  helpers.EnforceHTTP(body.URL)

	r2.Decr(database.Ctx, c.IP())
}