package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/anond0rf/vecchioclient/client"
	"github.com/labstack/echo/v4"
)

// APIHandler handles the requests to vecchiochan.
type APIHandler struct {
	vc *client.VecchioClient
}

// NewAPIHandler creates a new handler with a VecchioClient.
func NewAPIHandler(userAgent string, verbose bool) *APIHandler {
	conf := client.DefaultConfig
	if strings.TrimSpace(userAgent) != "" {
		conf.UserAgent = userAgent
	}
	conf.Verbose = verbose
	return &APIHandler{
		vc: client.NewVecchioClientWithConfig(conf),
	}
}

// NewThread Create a new thread on a specific board
func (h *APIHandler) NewThread(ctx echo.Context) error {
	log.Println("New thread request received")

	var body Thread
	if err := ctx.Bind(&body); err != nil {
		log.Printf("Invalid request body: %v\n", err)
		resp := ErrorResponse{
			Error: "Invalid request body",
		}
		return echo.NewHTTPError(http.StatusBadRequest, resp)
	}

	thread := toClientThread(body)

	log.Printf("Thread to be posted: %+v\n", thread)

	id, err := h.vc.NewThread(thread)
	if err != nil {
		log.Printf("Failed to create thread: %v\n", err)
		resp := ErrorResponse{
			Error: "Failed to create thread",
		}
		return echo.NewHTTPError(http.StatusInternalServerError, resp)
	}

	log.Printf("Thread #%d posted successfully\n", id)

	resp := SuccessResponse{
		Id: id,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// PostReply Post a reply to an existing thread
func (h *APIHandler) PostReply(ctx echo.Context) error {
	log.Println("Post reply request received")

	var body Reply
	if err := ctx.Bind(&body); err != nil {
		log.Printf("Invalid request body: %v\n", err)
		resp := ErrorResponse{
			Error: "Invalid request body",
		}
		return echo.NewHTTPError(http.StatusBadRequest, resp)
	}

	reply := toClientReply(body)

	log.Printf("Reply to be posted: %+v\n", reply)

	id, err := h.vc.PostReply(reply)
	if err != nil {
		log.Printf("Failed to post reply: %v\n", err)
		resp := ErrorResponse{
			Error: "Failed to post reply",
		}
		return echo.NewHTTPError(http.StatusInternalServerError, resp)
	}

	log.Printf("Reply #%d posted successfully\n", id)

	resp := SuccessResponse{
		Id: id,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func toClientThread(body Thread) client.Thread {
	thread := client.Thread{
		Board: body.Board,
	}
	if body.Name != nil {
		thread.Name = *body.Name
	}
	if body.Email != nil {
		thread.Email = *body.Email
	}
	if body.Subject != nil {
		thread.Subject = *body.Subject
	}
	if body.Spoiler != nil {
		thread.Spoiler = *body.Spoiler
	}
	if body.Body != nil {
		thread.Body = *body.Body
	}
	if body.Embed != nil {
		thread.Embed = *body.Embed
	}
	if body.Password != nil {
		thread.Password = *body.Password
	}
	if body.Sage != nil {
		thread.Sage = *body.Sage
	}
	if body.Files != nil {
		thread.Files = *body.Files
	}
	return thread
}

func toClientReply(body Reply) client.Reply {
	thread := client.Reply{
		Thread: body.Thread,
		Board:  body.Board,
	}
	if body.Name != nil {
		thread.Name = *body.Name
	}
	if body.Email != nil {
		thread.Email = *body.Email
	}
	if body.Spoiler != nil {
		thread.Spoiler = *body.Spoiler
	}
	if body.Body != nil {
		thread.Body = *body.Body
	}
	if body.Embed != nil {
		thread.Embed = *body.Embed
	}
	if body.Password != nil {
		thread.Password = *body.Password
	}
	if body.Sage != nil {
		thread.Sage = *body.Sage
	}
	if body.Files != nil {
		thread.Files = *body.Files
	}
	return thread
}
