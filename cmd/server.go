package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/buaazp/fasthttprouter"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"github.com/yourusername/k8s-controller-tutorial/pkg/api"
	"github.com/yourusername/k8s-controller-tutorial/pkg/ctrl"
	"github.com/yourusername/k8s-controller-tutorial/pkg/informer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var serverPort int
var serverKubeconfig string
var serverInCluster bool
var enableLeaderElection bool
var metricsPort int
var enableMCP bool
var mcpPort int

type rootFlagsStruct struct {
	MetricsBindAddress string
}

var rootFlags = rootFlagsStruct{}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server and deployment informer",
	Run: func(cmd *cobra.Command, args []string) {
		level := parseLogLevel(logLevel)
		configureLogger(level)
		clientset, err := getServerKubeClient(serverKubeconfig, serverInCluster)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}
		ctx := context.Background()
		mgr, err := ctrlruntime.NewManager(ctrlruntime.GetConfigOrDie(), manager.Options{
			LeaderElection:   enableLeaderElection,
			LeaderElectionID: "k8s-controller-tutorial-leader-election",
			Metrics:          server.Options{BindAddress: rootFlags.MetricsBindAddress},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create controller manager")
			os.Exit(1)
		}

		// Register FrontendPage controller
		if err := ctrl.SetupFrontendPageController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to set up FrontendPage controller")
			os.Exit(1)
		}

		// --- API ROUTER SETUP ---
		router := fasthttprouter.New()
		router.POST("/api/token", api.TokenHandler)
		frontendAPI := &api.FrontendPageAPI{
			K8sClient: mgr.GetClient(),
			Namespace: "default", // or make configurable
		}
		router.GET("/api/frontendpages", api.JWTMiddleware(frontendAPI.ListFrontendPages))
		router.POST("/api/frontendpages", api.JWTMiddleware(frontendAPI.CreateFrontendPage))
		router.GET("/api/frontendpages/:name", api.JWTMiddleware(frontendAPI.GetFrontendPage))
		router.PUT("/api/frontendpages/:name", api.JWTMiddleware(frontendAPI.UpdateFrontendPage))
		router.DELETE("/api/frontendpages/:name", api.JWTMiddleware(frontendAPI.DeleteFrontendPage))

		// Legacy endpoint for deployments
		router.GET("/deployments", func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.Set("Content-Type", "application/json")
			deployments := informer.GetDeploymentNames()
			ctx.SetStatusCode(200)
			ctx.Write([]byte("["))
			for i, name := range deployments {
				ctx.WriteString("\"")
				ctx.WriteString(name)
				ctx.WriteString("\"")
				if i < len(deployments)-1 {
					ctx.WriteString(",")
				}
			}
			ctx.Write([]byte("]"))
		})

		go informer.StartDeploymentInformer(ctx, clientset)
		go func() {
			log.Info().Msg("Starting controller-runtime manager...")
			if err := mgr.Start(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Manager exited with error")
				os.Exit(1)
			}
		}()

		if enableMCP {
			go func() {
				mcpServer := NewMCPServer("K8s Controller MCP", appVersion)
				sseServer := mcpserver.NewSSEServer(mcpServer,
					mcpserver.WithBaseURL(fmt.Sprintf("http://:%d", mcpPort)),
				)
				log.Info().Msgf("Starting MCP server in SSE mode on port %d", mcpPort)
				if err := sseServer.Start(fmt.Sprintf(":%d", mcpPort)); err != nil {
					log.Fatal().Err(err).Msg("MCP SSE server error")
				}
			}()
			log.Info().Msgf("MCP server ready on port %d", mcpPort)
		}

		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s", addr)
		if err := fasthttp.ListenAndServe(addr, router.Handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func getServerKubeClient(kubeconfigPath string, inCluster bool) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
	serverCmd.Flags().StringVar(&serverKubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	serverCmd.Flags().BoolVar(&serverInCluster, "in-cluster", false, "Use in-cluster Kubernetes config")
	serverCmd.Flags().BoolVar(&enableLeaderElection, "enable-leader-election", true, "Enable leader election for controller manager")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8081, "Port for controller manager metrics")
	serverCmd.Flags().BoolVar(&enableMCP, "enable-mcp", false, "Enable MCP server")
	serverCmd.Flags().IntVar(&mcpPort, "mcp-port", 9090, "Port for MCP server")
	rootFlags.MetricsBindAddress = fmt.Sprintf(":%d", metricsPort)
}
