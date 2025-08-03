# MongoDB Database Syncer

A lightweight Go application that synchronizes MongoDB databases by copying all collections from a source database to a destination database.

## Features

- Full database synchronization with all collections
- Concurrent processing for better performance
- Transaction support for data consistency
- Simple command-line interface

## Prerequisites

- Go 1.23.0 or later
- MongoDB instance (local or remote)

## Installation

```bash
go install github.com/princecee/mongo-db-syncer/cmd/dbsyncer@latest
```

## Usage

```bash
dbsyncer -mongo_uri=<connection_string> -source_db=<source_database> -dest_db=<destination_database>
```

### Required Flags

- `-mongo_uri`: MongoDB connection URI (without database name)
- `-source_db`: Name of the source database to copy from
- `-dest_db`: Name of the destination database to copy to

### Examples

```bash
# Local MongoDB
dbsyncer -mongo_uri="mongodb://localhost:27017" -source_db="production" -dest_db="backup"

# With Authentication
dbsyncer -mongo_uri="mongodb://username:password@localhost:27017" -source_db="myapp" -dest_db="myapp_backup"

# MongoDB Atlas
dbsyncer -mongo_uri="mongodb+srv://username:password@cluster.mongodb.net" -source_db="prod_db" -dest_db="staging_db"
```

## License

This project is licensed under the [MIT License](LICENSE).
