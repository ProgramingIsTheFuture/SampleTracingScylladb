docker exec -i sampletracingscylladb_scylladb_1 cqlsh -ucassandra -pcassandra <<EOF
CREATE KEYSPACE IF NOT EXISTS test WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }; 
use test;
CREATE TABLE IF NOT EXISTS users (
	id uuid,
	username text,
	PRIMARY KEY (id, username),
);
exit;
EOF

