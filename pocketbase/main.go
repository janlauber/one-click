package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/janlauber/one-click/hooks"
	"github.com/janlauber/one-click/pkg/controller"
	"github.com/janlauber/one-click/pkg/env"
	"github.com/janlauber/one-click/pkg/k8s"
	"github.com/janlauber/one-click/pkg/watcher"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/cron"
)

func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_public"
	}

	return filepath.Join(os.Args[0], "../pb_public")
}

func init() {
	env.Init()
	k8s.Init()
}

func main() {
	app := pocketbase.New()

	var publicDirFlag string

	// add "--publicDir" option flag
	app.RootCmd.PersistentFlags().StringVar(
		&publicDirFlag,
		"publicDir",
		defaultPublicDir(),
		"the directory to serve static files",
	)
	migrationsDir := "" // default to "pb_migrations" (for js) and "migrations" (for go)

	// load js files to allow loading external JavaScript migrations
	jsvm.MustRegister(app, jsvm.Config{
		// Dir: migrationsDir,
		MigrationsDir: migrationsDir,
	})

	// register the `migrate` command
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangJS, // or migratecmd.TemplateLangGo (default)
		Dir:          migrationsDir,
		Automigrate:  true,
	})

	// call this only if you want to use the configurable "hooks" functionality
	hooks.PocketBaseInit(app)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDirFlag), true))

		return nil
	})

	app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		switch e.Collection.Name {
		case "rollouts":
			return controller.HandleRolloutCreate(e, app)
		}
		return nil
	})

	app.OnRecordBeforeUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		switch e.Collection.Name {
		case "rollouts":
			return controller.HandleRolloutUpdate(e, app)
		}
		return nil
	})

	app.OnRecordBeforeDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		switch e.Collection.Name {
		case "rollouts":
			return controller.HandleRolloutDelete(e, app)
		case "deployments":
			return controller.HandleDeploymentDelete(e, app)
		case "projects":
			return controller.HandleProjectDelete(e, app)
		}
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// get status of a specific rollout
		e.Router.GET("/rollouts/:projectId/:deploymentId/status", func(c echo.Context) error {
			projectId := c.PathParam("projectId")
			deploymentId := c.PathParam("deploymentId")

			return controller.HandleRolloutStatus(c, app, projectId, deploymentId)
		}, apis.RequireRecordAuth("users"))

		e.Router.GET("/rollouts/:projectId/:deploymentId/metrics", func(c echo.Context) error {
			projectId := c.PathParam("projectId")
			deploymentId := c.PathParam("deploymentId")

			return controller.HandleRolloutMetrics(c, app, projectId, deploymentId)
		}, apis.RequireRecordAuth("users"))

		e.Router.GET("/rollouts/:projectId/:deploymentId/events", func(c echo.Context) error {
			projectId := c.PathParam("projectId")
			deploymentId := c.PathParam("deploymentId")

			return controller.HandleRolloutEvents(c, app, projectId, deploymentId)
		}, apis.RequireRecordAuth("users"))

		e.Router.GET("/rollouts/:projectId/:podName/logs", func(c echo.Context) error {
			projectId := c.PathParam("projectId")
			podName := c.PathParam("podName")

			return k8s.GetRolloutLogs(c.Response().Writer, projectId, podName)
		}, apis.RequireRecordAuth("users"))

		e.Router.GET("/pb/blueprints/:blueprintId", func(c echo.Context) error {
			blueprintId := c.PathParam("blueprintId")

			return controller.HandleBlueprint(c, app, blueprintId)
		}, apis.RequireRecordAuth("users"))

		e.Router.POST("/pb/blueprints/shared/:blueprintId", func(c echo.Context) error {
			blueprintId := c.PathParam("blueprintId")

			return controller.HandleBlueprintAdd(c, app, blueprintId)
		}, apis.RequireRecordAuth("users"))

		e.Router.POST("/auto-update/:autoUpdateId", func(c echo.Context) error {
			// TODO: change the auth to be a token generated by the user
			autoUpdateId := c.PathParam("autoUpdateId")

			return controller.HandleAutoUpdate(c, app, autoUpdateId)
		})

		e.Router.GET("/cluster-info", func(c echo.Context) error {
			return controller.HandleClusterInfo(c, app)
			// }, apis.RequireRecordAuth("users"))
		})

		// delete a pod of a rollout by pod name
		e.Router.DELETE("/rollouts/:projectId/:podName", func(c echo.Context) error {
			projectId := c.PathParam("projectId")
			podName := c.PathParam("podName")

			return controller.HandlePodDelete(c, app, projectId, podName)
		}, apis.RequireRecordAuth("users"))

		// websocket for rollout status
		e.Router.GET("/ws/k8s/rollouts", watcher.WsK8sRolloutsHandler)

		// websocket for pod logs
		e.Router.GET("/ws/k8s/logs", watcher.WsK8sRolloutLogsHandler)

		// websocket for rollout events
		e.Router.GET("/ws/k8s/events", watcher.WsK8sRolloutEventsHandler)

		return nil
	})

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		scheduler := cron.New()

		// Run image update every minute
		scheduler.MustAdd("autoUpdate", env.Config.CronTick, func() {
			err := controller.AutoUpdateController(app)
			if err != nil {
				log.Printf("Error updating image: %v\n", err)
			}
		})

		scheduler.Start()
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
