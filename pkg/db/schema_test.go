package db_test

import (
	"testing"

	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	migrations, err := db.Migrations()
	require.NoError(t, err, "should have been able to load migrations")
	require.GreaterOrEqual(t, len(migrations), 2, "wrong number of migrations, has a migration been added?")

	// The first three migrations should match our fixtures
	expected := []*db.Migration{
		{
			ID:   0,
			Name: "Migrations",
			Path: "0000_migrations.sql",
		},
		{
			ID:   1,
			Name: "Users Authn",
			Path: "0001_users_authn.sql",
		},
	}

	for i, migration := range migrations {
		if i > len(expected) {
			break
		}

		require.Equal(t, expected[i].ID, migration.ID)
		require.Equal(t, expected[i].Name, migration.Name)
		require.Equal(t, expected[i].Path, migration.Path)

		query, err := migration.SQL()
		require.NoError(t, err, "could not load SQL from the migration")
		require.NotEmpty(t, query, "no SQL was returned for the migration")
	}
}