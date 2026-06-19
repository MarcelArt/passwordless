/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/routes"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	webroutes "github.com/MarcelArt/passwordless/web/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.SetupENV()
		configs.ConnectDB()
		configs.SetupOauth2()

		if configs.Env.ServerENV == "prod" {
			gin.SetMode(gin.ReleaseMode)
		}

		uRepo := repositories.NewUserRepo(configs.DB)
		authService, err := services.NewAuthService(uRepo)
		if err != nil {
			log.Fatalf("failed constructing auth service: %s", err.Error())
		}

		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"POST, OPTIONS, GET, PUT, PATCH, DELETE"},
			AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		api := r.Group("/api")
		routes.SetupRoutes(api, authService)

		webroutes.SetupWebRoutes(r, authService)

		port := fmt.Sprintf(":%s", configs.Env.PORT)
		r.Run(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
