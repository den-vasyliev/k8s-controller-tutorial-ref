package cmd

import (
	"fmt"
	"os"

	"github.com/buaazp/fasthttprouter"
	"github.com/go-logr/zerologr"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"github.com/yourusername/k8s-controller-tutorial/pkg/api"
	"github.com/yourusername/k8s-controller-tutorial/pkg/ctrl"
	"github.com/yourusername/k8s-controller-tutorial/pkg/informer"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	frontendv1alpha1 "github.com/yourusername/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var serverPort int
var serverKubeconfig string
var serverInCluster bool
var enableLeaderElection bool
var leaderElectionNamespace string
var metricsPort int
var enableMCP bool
var mcpPort int
var FrontendAPI *api.FrontendPageAPI
var jwtSecret string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server and deployment informer",
	Run: func(cmd *cobra.Command, args []string) {
		level := parseLogLevel(logLevel)
		configureLogger(level)

		logf.SetLogger(zap.New(zap.UseDevMode(true)))
		logf.SetLogger(zerologr.New(&log.Logger))

		scheme := runtime.NewScheme()
		if err := clientgoscheme.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Failed to add client-go scheme")
			os.Exit(1)
		}
		if err := frontendv1alpha1.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Failed to add FrontendPage scheme")
			os.Exit(1)
		}
		mgr, err := ctrlruntime.NewManager(ctrlruntime.GetConfigOrDie(), manager.Options{
			Scheme:                  scheme,
			LeaderElection:          enableLeaderElection,
			LeaderElectionID:        "k8s-controller-tutorial-leader-election",
			LeaderElectionNamespace: leaderElectionNamespace,
			Metrics:                 server.Options{BindAddress: fmt.Sprintf(":%d", metricsPort)},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create controller manager")
			os.Exit(1)
		}

		if err := ctrl.AddFrontendController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add frontend controller")
			os.Exit(1)
		}

		// --- API ROUTER SETUP ---
		router := fasthttprouter.New()
		router.POST("/api/token", api.TokenHandler)
		frontendAPI := &api.FrontendPageAPI{
			K8sClient: mgr.GetClient(),
			Namespace: "default", // or make configurable
		}
		api.FrontendAPI = frontendAPI
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

		// Set the JWT secret for the API package
		api.JWTSecret = jwtSecret

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
	serverCmd.Flags().StringVar(&leaderElectionNamespace, "leader-election-namespace", "default", "Namespace for leader election")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8081, "Port for controller manager metrics")
	serverCmd.Flags().BoolVar(&enableMCP, "enable-mcp", false, "Enable MCP server")
	serverCmd.Flags().IntVar(&mcpPort, "mcp-port", 9090, "Port for MCP server")
	serverCmd.Flags().StringVar(&jwtSecret, "jwt-secret", "", "Secret key for signing JWT tokens (required)")
}
