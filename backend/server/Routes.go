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
	app.Get("/parse/:encoded_url", controllers.Parse)

	// Problem Routes
	app.Get("/problem/:encoded_url", controllers.GetProblemList)

	app.Get("/problem/:platform/:cid/:label", controllers.GetProblem)

	app.Post("/problem/custom/add", controllers.AddCustomProblem)

	// Config Routes
	app.Get("/config", controllers.GetConfig)

	app.Put("/config", controllers.UpdateConfig)

	// Code Routes
	app.Get("/code/:platform/:cid/:label", controllers.GetCodeByProblemPath)

	app.Put("/code", controllers.GetCodeByFilePath)

	app.Put("/code/update", controllers.UpdateCodeByFilePath)

	app.Put("/code/update/:platform/:cid/:label", controllers.UpdateCodeByProblemPath)

	// Source Routes
	app.Get("/source/open/:platform/:cid/:label", controllers.OpenSourceByMetadata)

	app.Get("/source/generate/:platform/:cid/:label", controllers.GenerateSourceByProblemPath)

	// Testcase Routes
	app.Put("/testcase/custom", controllers.GetCustomTestByPath)

	app.Post("/testcase/custom/add", controllers.AddCustomTest)

	app.Put("/testcase/custom/update", controllers.UpdateCustomTest)

	app.Delete("/testcase/custom/delete", controllers.DeleteCustomTest)

	app.Post("/testcase/random/input/generate", controllers.GenerateRandomTests)

	app.Post("/testcase/random/output/generate", controllers.GenerateOutputs)

	// Test Routes
	app.Get("/test/:platform/:cid/:label", controllers.TestProblem)

	app.Get("/execresult/:platform/:cid/:label", controllers.GetExecutionResult)

	// Misc. Routes
	app.Get("/directories/:encoded_url", controllers.GetInputOutputDirectoriesByUrl)

	app.Get("/misc/directory/check/:encoded_path", controllers.CheckDirectoryPathValidity)

	app.Get("/misc/filepath/check/:encoded_path", controllers.CheckFilePathValidity)

	app.Put("/misc/resource/open", controllers.OpenResource)
}
