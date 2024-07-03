package tools

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestSessionOnLocal(t *testing.T) {
	//declare fiber context
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	//declare session
	sess := Session{}

	//set session for local use
	sess.Initiate(true)

	//set session id
	err := sess.SetID(ctx, "test")
	assert.Nil(t, err, "error should be nil")

	//get session id
	val, err := sess.Get(ctx)
	assert.Nil(t, err, "error should nil")
	assert.Equal(t, val, "M11ZZZZ00003", "session id should be same as assigned")

	err = sess.Remove(ctx)
	assert.Nil(t, err, "error should nil")

}

func TestSessionOnProduction(t *testing.T) {
	//declare fiber context
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	//declare session
	sess := Session{}

	//set session for local use
	sess.Initiate(false)

	//set session id
	err := sess.SetID(ctx, "test")
	assert.Nil(t, err, "error should be nil")

	//get session id
	val, err := sess.Get(ctx)
	assert.Nil(t, err, "error should nil")
	assert.Equal(t, val, "test", "session id should be same as assigned")

	err = sess.Remove(ctx)
	assert.Nil(t, err, "error should nil")

}
