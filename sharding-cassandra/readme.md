# Setup
- Connect to the Cassandra cluster
```sh
docker exec -it cassandra1 cqlsh
```

- Create a keyspace and a table
```sql
CREATE KEYSPACE test_keyspace WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 2};

USE test_keyspace;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT,
    email TEXT
);
```

- Insert some data
```sql
INSERT INTO users (id, name, email) VALUES (uuid(), 'Alice', 'alice@example.com');
INSERT INTO users (id, name, email) VALUES (uuid(), 'Bob', 'bob@example.com');
INSERT INTO users (id, name, email) VALUES (uuid(), 'Charlie', 'charlie@example.com');
```