//go:build !linux && !android

package cmd

import (
	"github.com/hxhieu/b1-timetask-cli-go/console"
	"github.com/hxhieu/b1-timetask-cli-go/gui"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func (c *guiCmd) Run(ctx CLIContext) error {
	if !ctx.Experimental {
		console.ErrorLn("Coming soon (TM)")
		return nil
	}

	// DEBUG MODE ONLY

	// Create an instance of the app structure
	app := gui.NewApp()

	// Create application with options
	if err := wails.Run(&options.App{
		Title:  APP_TITLE,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: *ctx.GuiAssets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	}); err != nil {
		return err
	}

	return nil
}
