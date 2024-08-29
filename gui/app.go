package gui

import (
	"context"

	"github.com/hxhieu/b1-timetask-cli-go/api/intervals"
	"github.com/hxhieu/b1-timetask-cli-go/common"
)

// App struct
type App struct {
	ctx        context.Context
	taskParser *common.TaskCsvParser
	apiClient  *intervals.Client
	debugMode  bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) FetchTaskInputs() ([]*common.TimeTaskInput, error) {
	if a.taskParser == nil {
		parse, err := common.NewTaskParser(nil)
		if err != nil {
			return nil, err
		}
		a.taskParser = parse
	}
	return a.taskParser.Tasks, nil
}

func (a *App) InitUser() (*string, error) {
	token, err := common.GetUserToken()
	if err != nil {
		return nil, err
	}

	a.apiClient = intervals.New(token, a.debugMode)

	me, err := a.apiClient.Me()
	if err != nil {
		return nil, err
	}

	return &me.Email, nil
}
