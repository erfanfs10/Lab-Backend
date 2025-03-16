package handlers

import (
	"errors"
	"fmt"

	"github.com/erfanfs10/Lab-Backend/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var allowedFormats = [3]string{"png", "jpeg", "webp"}

var wsConnections = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ProgressMessage struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Step     int    `json:"step"`
	URL      string `json:"url"`
}

func Home(c echo.Context) error {
	return c.String(200, "Lab Backend")
}

func WebsocketHandler(c echo.Context) error {
	// Upgrade to WebSocket
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError,
			err, "websocket connection error")
	}
	// print a log into terminal
	fmt.Println(c.RealIP(), " ws connected")
	// store ws to cns
	wsConnections[c.RealIP()] = ws
	// close connection after function ends
	defer ws.Close()
	// delete connection from map after function ends
	defer delete(wsConnections, c.RealIP())
	// listen and write to the connection
	for {
		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return utils.HandleError(c, http.StatusInternalServerError,
				err, "websocket connection error")
		}
		// print the message to terminal
		fmt.Printf("%s from %s\n", msg, c.RealIP())
	}
}

func Convert(c echo.Context) error {
	// get websocket connection
	ws := wsConnections[c.RealIP()]
	if ws == nil {
		return utils.HandleError(c, http.StatusBadRequest,
			errors.New("can not get ws from wsConnections"),
			"There is no connection with server")
	}
	// send uploading status to client
	ws.WriteJSON(ProgressMessage{Status: "Uploading", Progress: 33, Step: 1})
	// get format from query params
	fileFormat := c.QueryParam("format")
	if fileFormat == "" {
		return utils.HandleError(c, http.StatusBadRequest,
			errors.New("no file format sent"), "No file format sent")
	}
	// Get image from form
	image, err := c.FormFile("image")
	if err != nil {
		return utils.HandleError(c, http.StatusBadRequest,
			err, "No file sent")
	}
	// Check if file is image
	fileType := strings.Split(image.Header["Content-Type"][0], "/")[0] // like image
	if fileType != "image" {
		return utils.HandleError(c, http.StatusBadRequest,
			errors.New("the file is not an image"), "The File Is Not An Image")
	}
	// Check if image is .jpeg or .png or .webp
	imageType := strings.Split(image.Header["Content-Type"][0], "/")[1] // like png
	if !utils.CheckFileFormat(allowedFormats, imageType) {
		return utils.HandleError(c, http.StatusBadRequest,
			errors.New("the file format not allowed"), "The File Format Not Allowed")
	}
	// Check image size: up to 5MB not allowed
	if image.Size > 5000000 {
		return utils.HandleError(c, http.StatusBadRequest,
			errors.New("the image is too large"), "The Image Is Too Large")
	}
	// open image
	src, err := image.Open()
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError,
			err, "Can not open the file")
	}
	// close it when func is endded
	defer src.Close()
	// store the file extention
	fileExt := filepath.Ext(image.Filename) // like .png or .jpeg
	// create directory name for the image
	dir := "../images/"
	// current time
	currentTime := strings.Join(strings.Split(time.Now().Format(time.DateTime), " "), "_")
	// create the image name
	fileName := fmt.Sprintf("original_%s%s", currentTime, fileExt) // like original_2025-02-23_23:38:05.png
	// create destination path
	dstPath := filepath.Join(dir, fileName) // like ../images/original_2025-02-23_23:38:05.png
	// create the destination
	dst, err := os.Create(dstPath)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError,
			err, "Can not create the dst")
	}
	// close the dst after func is endded
	defer dst.Close()
	// store the file in the dst from src
	_, err = io.Copy(dst, src)
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError,
			err, "Can not copy the content")
	}
	// send converting status to client
	ws.WriteJSON(ProgressMessage{Status: "Converting", Progress: 66, Step: 2})
	// convert the file
	convertedFile := "compressed_" + currentTime + fileFormat // like compressed_2025-02-25_20:47:08.webp
	// create convert command from os
	cmd := exec.Command("ffmpeg", "-i", dstPath, dir+convertedFile)
	// run the convert command
	err = cmd.Run()
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError,
			err, "Can not convert")
	}
	// create delete the original image command
	cmd = exec.Command("rm", dstPath)
	// run the command
	err = cmd.Run()
	if err != nil {
		return utils.HandleError(c, http.StatusInternalServerError,
			err, "Can not convert")
	}
	// send done status and converted file path to client
	ws.WriteJSON(ProgressMessage{Status: "Done", Progress: 100, Step: 3,
		URL: "http://127.0.0.1:8000/static/images/" + convertedFile})
	// return the response
	return c.JSON(http.StatusOK, echo.Map{"message": "image converted"})
}
