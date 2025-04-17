package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/praveent04/url-short/database"
	"github.com/redis/go-redis/v9"
)


func ResolveURL(c fiber.Ctx) error{
	url := c.Params("url")

	r:= database.CreateClient(0)

	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"short not found on database"})
	} else if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": " unable to connect to db"})
	}
	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect().Status(301).To(value)

}