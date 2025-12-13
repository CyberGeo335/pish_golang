//go:build integration

package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/CyberGeo335/prak_sixteen/internal/db"
	"github.com/CyberGeo335/prak_sixteen/internal/httpapi"
	"github.com/CyberGeo335/prak_sixteen/internal/repo"
	"github.com/CyberGeo335/prak_sixteen/internal/service"
)

func withPostgres(t *testing.T) (dsn string, term func()) {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "notes_test",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort("5432/tcp"),
		).WithStartupTimeout(30 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := c.Host(ctx)
	require.NoError(t, err)

	port, err := c.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	dsn = fmt.Sprintf("postgres://test:test@%s:%s/notes_test?sslmode=disable", host, port.Port())
	return dsn, func() { _ = c.Terminate(ctx) }
}

func newTestServer(t *testing.T, dsn string) (baseURL string, closeFn func()) {
	t.Helper()

	dbx, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	require.NoError(t, dbx.Ping())

	db.MustApplyMigrations(dbx)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	svc := service.New(repo.NoteRepo{DB: dbx})
	httpapi.Router{Svc: svc}.Register(r)

	srv := httptest.NewServer(r)
	return srv.URL, func() {
		srv.Close()
		_ = dbx.Close()
	}
}

type noteDTO struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func Test_CreateAndGet_withTC(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	url, closeSrv := newTestServer(t, dsn)
	defer closeSrv()

	// Create
	resp, err := http.Post(url+"/notes", "application/json",
		strings.NewReader(`{"title":"Hello","content":"World"}`))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	var created noteDTO
	require.NoError(t, json.Unmarshal(body, &created))
	require.True(t, created.ID > 0)

	// Get
	resp2, err := http.Get(fmt.Sprintf("%s/notes/%d", url, created.ID))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)
	_ = resp2.Body.Close()
}

func Test_Get_NotFound_withTC(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	url, closeSrv := newTestServer(t, dsn)
	defer closeSrv()

	resp, err := http.Get(url + "/notes/999999")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	_ = resp.Body.Close()
}

func Test_Create_BadJSON_withTC(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	url, closeSrv := newTestServer(t, dsn)
	defer closeSrv()

	resp, err := http.Post(url+"/notes", "application/json",
		strings.NewReader(`{"title":`)) // сломанный JSON
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	_ = resp.Body.Close()
}
