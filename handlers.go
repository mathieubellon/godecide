package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func Homepage(ctx *fiber.Ctx) error {
	session, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	log.Println(session.ID())
	log.Println(session.Get("userEmail"))
	return ctx.Render("./index.html", fiber.Map{"Email": session.Get("userEmail")})
}

func ListIdeas(ctx *fiber.Ctx) error {
	store, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return ctx.JSON(&fiber.Map{
		"page":    "ideas",
		"session": store.ID(),
	})
}

func Me(ctx *fiber.Ctx) error {
	sess, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}

	return ctx.JSON(&fiber.Map{
		"page":          "me",
		"session":       sess.ID(),
		"user":          sess.Get("userEmail"),
		"provider":      sess.Get("provider"),
		"userID":        sess.Get("userID"),
		"workspaceID":   sess.Get("workspaceID"),
		"workspaceName": sess.Get("workspaceName"),
		"fresh":         sess.Fresh(),
		"keys":          sess.Keys(),
	})
}

func Logout(ctx *fiber.Ctx) error {
	if err := goth_fiber.Logout(ctx); err != nil {
		log.Fatal(err)
	}
	// Destroy session
	sess, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	if err != nil {
		panic(err)
	}
	if err := sess.Destroy(); err != nil {
		panic(err)
	}

	return ctx.Redirect("/")
}

func Callback(ctx *fiber.Ctx) error {
	oauthResponse, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	user, err := FindOrCreateUser(oauthResponse)
	if err != nil {
		return ctx.SendStatus(500)
	}

	sess, err := globalSession.Get(ctx) // Get session ( creates one if not exist )
	if err != nil {
		return err
	}
	sess.Fresh()
	sess.Set("userEmail", user.Email)
	sess.Set("provider", "github")
	sess.Set("userID", user.ID)
	// sess.Set("workspaceID", wp.ID)
	// sess.Set("workspaceName", wp.Name)
	sess.Save()

	log.Println(sess)

	return ctx.JSON(oauthResponse)
}