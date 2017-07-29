package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/go-siris/siris"
	"github.com/go-siris/siris/context"
	"github.com/go-siris/siris/view"
)

func main() {
	app := siris.New()

	app.AttachView(view.HTML("./templates", ".html"))

	// Serve the form.html to the user
	app.Get("/upload", func(ctx context.Context) {
		//create a token (optionally)

		now := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(now, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		// render the form with the token for any use you like
		ctx.ViewData("", token)
		ctx.View("upload_form.html")
	})

	// Handle the post request from the upload_form.html to the server
	app.Post("/upload", context.LimitRequestBodySize(10<<20),
		func(ctx context.Context) {
			// or use ctx.SetMaxRequestBodySize(10 << 20)
			//to limit the uploaded file(s) size.

			// Get the file from the request
			file, info, err := ctx.FormFile("uploadfile")

			if err != nil {
				ctx.StatusCode(siris.StatusInternalServerError)
				ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
				return
			}

			fname := info.Filename

			// Create a file with the same name
			// assuming that you have a folder named 'uploads'
			out, err := os.OpenFile("./uploads/"+fname,
				os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				ctx.StatusCode(siris.StatusInternalServerError)
				ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
				file.Close()
				return
			}

			io.Copy(out, file)

			out.Close()
			file.Close()
		})

	// start the server at http://localhost:8080
	app.Run(siris.Addr(":8080"))
}
