package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/gondolf/pkg/simplesql/test"
	"github.com/stretchr/testify/require"
)

var clusterTableMigrations = []simplesql.Migration{
	{
		Version: 1,
		Up: `
			CREATE TABLE cluster (
				id VARCHAR(255) NOT NULL PRIMARY KEY,
				version BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				state VARCHAR(255) NOT NULL,
				message TEXT NOT NULL,
				cluster_manager_id VARCHAR(255) NOT NULL,
				created_at BIGINT NOT NULL,
				last_updated_at BIGINT NOT NULL,
				deleted_at BIGINT NOT NULL DEFAULT 0,
				UNIQUE (name, deleted_at)
			);
		`,
		Down: `
				DROP TABLE IF EXISTS cluster;
			`,
	},
}

func TestSimpleSqlDB(t *testing.T) {

	db, err := test.NewTestSQLiteDB()
	require.NoError(t, err)
	defer db.Close()

	simplesqlDb := simplesql.NewDatabase(db)
	err = simplesqlDb.ApplyMigrations(clusterTableMigrations)
	require.NoError(t, err)

	clusterTable := NewClusterTable(simplesqlDb)

	t.Run("Insert", func(t *testing.T) {

		for i := 0; i < 5; i++ {
			cluster := ClusterRow{
				ID:               fmt.Sprintf("cluster%d", i),
				Version:          1,
				Name:             fmt.Sprintf("cluster%d", i),
				ClusterManagerID: fmt.Sprintf("cluster_manager%d", i),
				State:            "active",
				Message:          fmt.Sprintf("cluster%d is active", i),
			}
			err := clusterTable.Insert(context.Background(), db, cluster)
			require.NoError(t, err)
		}
	})

	t.Run("Get by ID", func(t *testing.T) {
		cluster, err := clusterTable.Get(context.Background(), ClusterTableGetKeys{ID: StringPtr("cluster0")})
		require.NoError(t, err)
		require.Equal(t, "cluster0", cluster.Name)
		require.Equal(t, "active", cluster.State)
		require.Equal(t, "cluster0 is active", cluster.Message)
		require.Equal(t, uint64(1), cluster.Version)
		require.Equal(t, int64(0), cluster.DeletedAt)
	})

	t.Run("Update", func(t *testing.T) {
		err := clusterTable.Update(
			context.Background(), db,
			ClusterTableUpdateKey{
				ID:               "cluster0",
				Version:          1,
				ClusterManagerID: "cluster_manager0",
			},
			ClusterTableUpdateFields{
				State:   StringPtr("inactive"),
				Message: StringPtr("cluster0 is inactive"),
			})
		require.NoError(t, err)

		cluster, err := clusterTable.Get(context.Background(), ClusterTableGetKeys{ID: StringPtr("cluster0")})
		require.NoError(t, err)
		require.Equal(t, "cluster0", cluster.Name)
		require.Equal(t, "inactive", cluster.State)
		require.Equal(t, "cluster0 is inactive", cluster.Message)
		require.Equal(t, uint64(2), cluster.Version)

		for i := 1; i < 5; i++ {
			cluster, err := clusterTable.Get(context.Background(), ClusterTableGetKeys{ID: StringPtr(fmt.Sprintf("cluster%d", i))})
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("cluster%d", i), cluster.Name)
			require.Equal(t, "active", cluster.State)
			require.Equal(t, fmt.Sprintf("cluster%d is active", i), cluster.Message)
			require.Equal(t, uint64(1), cluster.Version)
		}
	})

	t.Run("List", func(t *testing.T) {
		clusters, err := clusterTable.List(context.Background(), ClusterTableSelectFilters{
			StateIn: []string{"active", "inactive"},
		})
		require.NoError(t, err)
		require.Len(t, clusters, 5)
	})

	t.Run("Soft Delete", func(t *testing.T) {
		err := clusterTable.Update(
			context.Background(), db,
			ClusterTableUpdateKey{
				ID:               "cluster0",
				Version:          2,
				ClusterManagerID: "cluster_manager0",
			},
			ClusterTableUpdateFields{
				State:     StringPtr("inactive"),
				Message:   StringPtr("cluster0 is inactive"),
				DeletedAt: Int64Ptr(time.Now().Unix()),
			},
		)
		require.NoError(t, err)

		cluster, err := clusterTable.Get(context.Background(), ClusterTableGetKeys{ID: StringPtr("cluster0")})
		require.ErrorAs(t, err, &simplesql.ErrRecordNotFound)
		require.Equal(t, ClusterRow{}, cluster)

		clusters, err := clusterTable.List(context.Background(), ClusterTableSelectFilters{
			StateIn:     []string{"active", "inactive"},
			DeletedAtEq: Int64Ptr(0),
		})
		require.NoError(t, err)
		require.Len(t, clusters, 4)
	})

	t.Run("List without deleted", func(t *testing.T) {
		clusters, err := clusterTable.List(context.Background(), ClusterTableSelectFilters{
			StateIn:     []string{"active", "inactive"},
			DeletedAtEq: Int64Ptr(0),
		})
		require.NoError(t, err)
		require.Len(t, clusters, 4)
	})

	t.Run("List with deleted", func(t *testing.T) {
		clusters, err := clusterTable.List(context.Background(), ClusterTableSelectFilters{
			StateIn: []string{"active", "inactive"},
		})
		require.NoError(t, err)
		require.Len(t, clusters, 5)
	})

	t.Run("Delete", func(t *testing.T) {
		err := clusterTable.Delete(context.Background(), db, ClusterTableUpdateKey{
			ID:               "cluster1",
			Version:          1,
			ClusterManagerID: "cluster_manager1",
		})
		require.NoError(t, err)

		cluster, err := clusterTable.Get(context.Background(), ClusterTableGetKeys{ID: StringPtr("cluster1")})
		require.ErrorAs(t, err, &simplesql.ErrRecordNotFound)
		require.Equal(t, ClusterRow{}, cluster)

		clusters, err := clusterTable.List(context.Background(), ClusterTableSelectFilters{
			StateIn: []string{"active", "inactive"},
		})
		require.NoError(t, err)
		require.Len(t, clusters, 4)

		clusterNames := []string{}
		for _, cluster := range clusters {
			clusterNames = append(clusterNames, cluster.Name)
		}
		require.NotContains(t, clusterNames, "cluster1")
	})
}

func StringPtr(s string) *string {
	return &s
}

func Int64Ptr(i int64) *int64 {
	return &i
}
