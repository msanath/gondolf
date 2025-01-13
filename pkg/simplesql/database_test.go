package simplesql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/msanath/gondolf/pkg/simplesql"
	"github.com/msanath/gondolf/pkg/simplesql/test"
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

type ClusterRow struct {
	ID            string `db:"id" orm:"op=get key:primary"`
	Version       uint64 `db:"version" orm:"op_lock:true"`
	CreatedAt     int64  `db:"created_at"`
	LastUpdatedAt int64  `db:"last_updated_at" orm:"op=update"`
	DeletedAt     int64  `db:"deleted_at" orm:"soft_delete:true"`

	Name             string `db:"name" orm:"op=get filter=In"`
	ClusterManagerID string `db:"cluster_manager_id" orm:"key:primary filter=In"`
	State            string `db:"state" orm:"op=update filter=In,NotIn"`
	Message          string `db:"message" orm:"op=update"`
}

type ClusterTableGetKeys struct {
	ID        *string `db:"id"`
	Name      *string `db:"name"`
	DeletedAt int64   `db:"deleted_at"`
}

type ClusterTableUpdateKey struct {
	ID               string `db:"id"`
	Version          uint64 `db:"version"`
	ClusterManagerID string `db:"cluster_manager_id"`
	DeleteAt         int64  `db:"deleted_at"`
}

type ClusterTableUpdateFields struct {
	State         *string `db:"state"`
	Message       *string `db:"message"`
	LastUpdatedAt *int64  `db:"last_updated_at"`
	DeletedAt     *int64  `db:"deleted_at"`
}

type ClusterTableSelectFilters struct {
	IDIn        []string `db:"id:in"`
	NameIn      []string `db:"name:in"`
	StateIn     []string `db:"state:in"`
	StateNotIn  []string `db:"state:not_in"`
	VersionGte  *uint64  `db:"version:gte"`
	VersionLte  *uint64  `db:"version:lte"`
	VersionEq   *uint64  `db:"version:eq"`
	DeletedAtEq *int64   `db:"deleted_at:eq"`
	Limit       uint32   `db:"limit"`
}

const clusterTableName = "cluster"

type ClusterTable struct {
	simplesql.Database
	tableName string
}

func NewClusterTable(db simplesql.Database) *ClusterTable {
	return &ClusterTable{
		Database:  db,
		tableName: clusterTableName,
	}
}

func (s *ClusterTable) Insert(ctx context.Context, execer sqlx.ExecerContext, row ClusterRow) error {
	return s.Database.Insert(ctx, execer, s.tableName, row)
}

func (s *ClusterTable) Get(ctx context.Context, keys ClusterTableGetKeys) (ClusterRow, error) {
	var row ClusterRow
	err := s.Database.Get(ctx, s.tableName, keys, &row)
	if err != nil {
		return ClusterRow{}, err
	}
	return row, nil
}

func (s *ClusterTable) Update(
	ctx context.Context, execer sqlx.ExecerContext, updateKey ClusterTableUpdateKey, updateFields ClusterTableUpdateFields,
) error {
	return s.Database.Update(ctx, execer, s.tableName, updateKey, updateFields)
}

func (s *ClusterTable) Delete(ctx context.Context, execer sqlx.ExecerContext, updateKey ClusterTableUpdateKey) error {
	return s.Database.Delete(ctx, s.tableName, updateKey)
}

func (s *ClusterTable) List(ctx context.Context, filters ClusterTableSelectFilters) ([]ClusterRow, error) {
	var rows []ClusterRow
	err := s.Database.List(ctx, s.tableName, filters, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
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
