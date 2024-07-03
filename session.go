package tools

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New()

type Session struct {
	store    *session.Store
	isLocal  bool
	mockCode string
}

func (c *Session) Initiate(local bool) {
	if local {
		c.isLocal = true
		c.store = session.New()
		c.mockCode = "M11ZZZZ00003"

	} else {
		c.store = session.New(session.Config{
			CookiePath:     "/",
			CookieSecure:   true,
			CookieHTTPOnly: true,
		})
	}
}

func (c *Session) SetID(context *fiber.Ctx, id string) error {
	sess, err := store.Get(context)
	if err != nil {
		return err
	}

	sess.Set("id", id)
	err = sess.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *Session) Get(context *fiber.Ctx) (string, error) {
	if c.isLocal {
		return c.mockCode, nil
	}
	sess, err := store.Get(context)
	if err != nil {
		return "", fmt.Errorf("fail get session : %w", err)
	}

	id := sess.Get("id")
	err = sess.Save()
	if err != nil {
		return "", fmt.Errorf("fail save session : %w", err)
	}

	if id == nil {
		return "", errors.New("id is nil")
	}

	return fmt.Sprint(id), nil
}

func (c *Session) Remove(context *fiber.Ctx) error {
	sess, err := store.Get(context)
	if err != nil {
		return err
	}

	sess.Delete("id")
	err = sess.Destroy()
	if err != nil {
		return err
	}

	return nil
}
