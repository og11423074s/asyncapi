package fixtures

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/og11423074s/asyncapi/config"
	"github.com/og11423074s/asyncapi/store"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

type TestEnv struct {
	Config *config.Config
	Db     *sql.DB
}

func NewTestEnv(t *testing.T) *TestEnv {
	os.Setenv("ENV", string(config.EnvTest))
	conf, err := config.New()
	require.NoError(t, err)

	db, err := store.NewPostgres(conf)
	require.NoError(t, err)

	return &TestEnv{
		Config: conf,
		Db:     db,
	}
}

func (te *TestEnv) SetupDb(t *testing.T) func(t *testing.T) {
	// run migrations
	m, err := migrate.New(
		fmt.Sprintf("file:///%s/migrations", te.Config.ProjectRoot),
		te.Config.DatabaseUrl())

	require.NoError(t, err)

	if err := m.Up(); err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	return te.TeardownDb
}

func (te *TestEnv) TeardownDb(t *testing.T) {
	_, err := te.Db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", strings.Join([]string{"users", "refresh_tokens", "reports"}, ", ")))
	require.NoError(t, err)
}
