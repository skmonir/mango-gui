package main

import (
	"context"
	"fmt"
	"github.com/skmonir/mango-ui/backend/server"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	//
	// Run the server for getting variable information.
	//
	go server.RunServer()
}
