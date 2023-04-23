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

func main() {
	app := fiber.New()

	app.Post("/upload", func(c *fiber.Ctx) error {
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
		byt, err := webp.EncodeRGBA(img, 80)
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
		return c.JSON(fiber.Map{
			"file": fileName,
		})
	})

	app.Listen(":3000")
}
