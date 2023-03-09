package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/skmonir/mango-gui/backend/server/controllers"
	"github.com/skmonir/mango-gui/backend/socket"
)

func SetRoutes(app fiber.Router) {
	// Socket Connection Routes
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		socket.CreateNewSocketConnection(conn)
	}))

	// Parse Routes
	app.Get("/parse/:encoded_url", parse)

	// Problem Routes
	app.Get("/problem/:encoded_url", getProblemList)

	app.Get("/problem/:platform/:cid/:label", getProblem)

	// Config Routes
	app.Get("/config", getConfig)

	app.Put("/config", updateConfig)

	// Code Routes
	app.Get("/code/:platform/:cid/:label", controllers.GetCodeByProblemPath)

	app.Put("/code", controllers.GetCodeByFilePath)

	app.Put("/code/update", controllers.UpdateCodeByFilePath)

	app.Put("/code/update/:platform/:cid/:label", controllers.UpdateCodeByProblemPath)

	// Source Routes
	app.Get("/source/open/:platform/:cid/:label", openSourceByMetadata)

	app.Put("/source/open", openSourceByPath)

	// Testcase Routes
	app.Put("/testcase/custom", getCustomTestByPath)

	app.Post("/testcase/custom/add", addCustomTest)

	app.Put("/testcase/custom/update", updateCustomTest)

	app.Delete("/testcase/custom/delete", deleteCustomTest)

	// Test Routes
	app.Get("/test/:platform/:cid/:label", testProblem)

	app.Get("/execresult/:platform/:cid/:label", getExecutionResult)

	// Misc. Routes
	app.Get("/directories/:encoded_url", getInputOutputDirectoriesByUrl)
}
