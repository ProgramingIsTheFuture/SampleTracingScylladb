docker exec -i sampletracingscylladb_scylladb_1 cqlsh -ucassandra -pcassandra <<EOF
CREATE KEYSPACE IF NOT EXISTS messages_service WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }; 
CREATE TABLE IF NOT EXISTS messages_service.messages (
	id uuid,
	content text,
	user_id uuid,
	PRIMARY KEY (id, user_id),
);

CREATE KEYSPACE IF NOT EXISTS users_service WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }; 
CREATE TABLE IF NOT EXISTS users_service.users (
	id uuid,
	username text,
	PRIMARY KEY (id, username),
);
exit;
EOF

