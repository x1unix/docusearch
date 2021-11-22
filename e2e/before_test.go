package e2e

import (
	"context"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"github.com/x1unix/docusearch/internal/config"
	"github.com/x1unix/docusearch/internal/web"
	"github.com/x1unix/docusearch/pkg/api"
	"go.uber.org/zap"
)

var (
	client      *api.Client
	redisClient redis.Cmdable
	storageDir  string
)

func TestMain(m *testing.M) {
	log.Println("Starting e2e test. Config file can be provided via E2E_CONFIG_FILE env variable.")
	cfgFile, ok := os.LookupEnv("E2E_CONFIG_FILE")
	if !ok {
		log.Println("Please set config file path in E2E_CONFIG_FILE environment variable:",
			"\n\t$ export E2E_CONFIG_FILE=config.dev.yml\n", "")
		cfgFile = "../config.dev.yml"
	}

	cfg, err := config.FromFile(cfgFile)
	if err != nil {
		log.Fatalln("failed to load config for e2e test:", err)
	}

	redisConn, err := cfg.RedisClient()
	if err != nil {
		log.Fatalln(err)
	}

	defer redisConn.Close()
	if err := redisConn.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("failed to connect to Redis:", err)
	}

	redisClient = redisConn
	storageDir = cfg.Storage.UploadsDirectory

	log.Println("cleaning up Redis...")
	if err := redisConn.Del(context.Background(), "*").Err(); err != nil {
		log.Fatalln("failed to clean redis data:", err)
	}

	log.Println("cleaning storage directory...")
	if err := os.RemoveAll(cfg.Storage.UploadsDirectory); err != nil {
		log.Fatalln("failed to remove uploads directory:", err)
	}

	svc := web.NewService(zap.NewNop(), cfg, redisConn)
	srv := httptest.NewServer(svc)
	defer srv.Close()

	log.Println("started HTTP server at:", srv.URL)
	client = api.NewClient(srv.URL)
	os.Exit(m.Run())
}

func readTestData(t *testing.T, fname string) []byte {
	d, err := ioutil.ReadFile(filepath.Join("testdata", fname))
	require.NoError(t, err, "failed to open testdata")
	return d
}

func cleanData(t *testing.T) {
	t.Log("cleaning up Redis...")
	if err := redisClient.Del(context.Background(), "*").Err(); err != nil {
		t.Fatal("failed to clean redis data:", err)
	}

	t.Log("cleaning storage directory...")
	if err := os.RemoveAll(storageDir); err != nil {
		t.Fatal("failed to remove uploads directory:", err)
	}
}

func assertResponseError(t *testing.T, err error, want api.ErrorResponse) {
	t.Helper()
	require.Error(t, err)
	rspErr, ok := err.(*api.ErrorResponse)
	if !ok {
		t.Fatalf("error is not %[1]T, got %[2]T %[2]s", want, err)
	}

	require.Equal(t, want.StatusCode, rspErr.StatusCode, "status code mismatch")
	require.Equal(t, want.Message, rspErr.Message, "error message mismatch")
}
