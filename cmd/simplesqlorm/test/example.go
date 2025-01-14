package test

//go:generate ../../../bin/simplesqlorm-gen --struct-name ClusterRow --pkg-name test --table-name=cluster
type ClusterRow struct {
	ID            string `db:"id" orm:"op=get key:primary filter=In"`
	Version       uint64 `db:"version" orm:"op_lock:true filter=Gte,Lte,Eq"`
	CreatedAt     int64  `db:"created_at"`
	LastUpdatedAt int64  `db:"last_updated_at" orm:"op=update"`
	DeletedAt     int64  `db:"deleted_at" orm:"soft_delete:true"`

	Name             string `db:"name" orm:"op=get filter=In"`
	ClusterManagerID string `db:"cluster_manager_id" orm:"key:primary filter=In"`
	State            string `db:"state" orm:"op=update filter=In,NotIn"`
	Message          string `db:"message" orm:"op=update"`
}
