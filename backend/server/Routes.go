package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/skmonir/mango-ui/backend/socket"
)

func SetRoutes(app fiber.Router) {
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		socket.CreateNewSocketConnection(conn)
	}))

	app.Get("/parse/:encoded_url", parse)

	app.Get("/problem/:encoded_url", getProblemList)

	app.Get("/problem/:platform/:cid/:label", getProblem)

	app.Get("/config", getConfig)

	app.Put("/config", updateConfig)

	app.Get("/code/:platform/:cid/:label", getCodeByMetadata)

	app.Put("/code", getCodeByPath)

	app.Get("/source/open/:platform/:cid/:label", openSourceByMetadata)

	app.Put("/source/open", openSourceByPath)

	app.Put("/testcase/custom", getCustomTestByPath)

	app.Post("/testcase/custom/add", addCustomTest)

	app.Put("/testcase/custom/update", updateCustomTest)

	app.Delete("/testcase/custom/delete", deleteCustomTest)

	app.Get("/test/:platform/:cid/:label", testProblem)

	app.Get("/execresult/:platform/:cid/:label", getExecutionResult)
}
