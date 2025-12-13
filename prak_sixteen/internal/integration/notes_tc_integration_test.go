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
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

	"pz16/internal/db"
	"pz16/internal/httpapi"
	"pz16/internal/repo"
	"pz16/internal/service"
)

// поднимаем Postgres на host-порту 5433 (чтобы не конфликтовать с локальным 5432)
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
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"5432/tcp": []nat.PortBinding{{
					HostIP:   "127.0.0.1",
					HostPort: "5433",
				}},
			}
		},
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

	return dsn, func() {
		_ = c.Terminate(ctx)
	}
}

func newTestServer(t *testing.T, dsn string) *httptest.Server {
	t.Helper()

	dbx, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbx.Close() })

	// ждём, пока Postgres окончательно поднимется
	deadline := time.Now().Add(20 * time.Second)
	for {
		if err := dbx.Ping(); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("db ping timeout")
		}
		time.Sleep(200 * time.Millisecond)
	}

	db.MustApplyMigrations(dbx)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	svc := service.New(repo.NoteRepo{DB: dbx})
	httpapi.Router{Svc: svc}.Register(r)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv
}

func Test_CreateAndGet_withTC(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	srv := newTestServer(t, dsn)

	// Create
	resp, err := http.Post(srv.URL+"/notes", "application/json",
		strings.NewReader(`{"title":"Hello","content":"World"}`))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var created map[string]any
	require.NoError(t, json.Unmarshal(body, &created))

	id := int64(created["id"].(float64))

	// Get
	resp2, err := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, id))
	require.NoError(t, err)
	defer resp2.Body.Close()

	require.Equal(t, http.StatusOK, resp2.StatusCode)
}

func Test_Get_NotFound_withTC(t *testing.T) {
	dsn, stop := withPostgres(t)
	defer stop()

	srv := newTestServer(t, dsn)

	resp, err := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, 999999))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
