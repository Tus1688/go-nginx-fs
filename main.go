package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var authorization string

const webpQuality = 80

func main() {
	app := fiber.New()

	authorization = os.Getenv("AUTHORIZATION")
	if authorization == "" {
		panic("AUTHORIZATION is not set")
	}
	app.Use(AuthMiddleware)
	app.Post("/handler", func(c *fiber.Ctx) error {
		file, err := c.FormFile("picture")
		if err != nil {
			return c.SendStatus(400)
		}
		// open the uploaded file
		src, err := file.Open()
		if err != nil {
			return c.SendStatus(400)
		}
		defer src.Close()

		var img image.Image
		// decode the image
		ext := filepath.Ext(file.Filename)
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(src)
		case ".png":
			img, err = png.Decode(src)
		default:
			return c.SendStatus(400)
		}
		if err != nil {
			return c.SendStatus(400)
		}
		fileName := uuid.New().String() + ".webp"

		// Encode the image as a WebP image
		byt, err := webp.EncodeRGBA(img, webpQuality)
		if err != nil {
			return c.SendStatus(500)
		}

		// create new file with the name of fileName and fill it with the encoded image
		dst, err := os.Create(fileName)
		if err != nil {
			return c.SendStatus(500)
		}
		defer dst.Close()

		// write the encoded image to the new file
		_, err = dst.Write(byt)
		if err != nil {
			return c.SendStatus(500)
		}
		// send json of the new file name
		return c.Status(201).JSON(fiber.Map{
			"file": fileName,
		})
	})

	app.Delete("/handler", func(c *fiber.Ctx) error {
		fileName := c.Query("file")
		if fileName == "" {
			return c.SendStatus(400)
		}
		// delete the file
		err := os.Remove(fileName)
		if err != nil {
			if os.IsNotExist(err) {
				return c.SendStatus(404)
			}
			return c.SendStatus(500)
		}
		return c.SendStatus(200)
	})

	err := os.Chdir("/usr/share/nginx/html")
	if err != nil {
		panic(err)
	}
	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func AuthMiddleware(c *fiber.Ctx) error {
	authValue := c.Get("Authorization")
	if authValue != authorization {
		return c.SendStatus(401)
	}
	return c.Next()
}
