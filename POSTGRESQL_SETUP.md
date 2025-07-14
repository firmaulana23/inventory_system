# PostgreSQL Setup Guide

The Inventory Management System now uses PostgreSQL as the primary database. This guide will help you set up PostgreSQL for the application.

## 🐘 PostgreSQL Setup Options

### Option 1: Docker Compose (Recommended)

The easiest way to set up PostgreSQL is using Docker Compose:

```bash
# Start PostgreSQL and Adminer
docker-compose up -d

# Check if containers are running
docker-compose ps

# View logs
docker-compose logs postgres
```

This will start:
- **PostgreSQL** on port `5432`
- **Adminer** (web-based database admin) on port `8081`

**Access Adminer:** http://localhost:8081
- System: PostgreSQL
- Server: postgres
- Username: postgres
- Password: postgres
- Database: inventory_system

### Option 2: Local PostgreSQL Installation

#### macOS (Homebrew)
```bash
# Install PostgreSQL
brew install postgresql

# Start PostgreSQL service
brew services start postgresql

# Create database
createdb inventory_system
```

#### Ubuntu/Debian
```bash
# Install PostgreSQL
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database
sudo -u postgres createdb inventory_system
```

#### CentOS/RHEL
```bash
# Install PostgreSQL
sudo yum install postgresql-server postgresql-contrib

# Initialize database
sudo postgresql-setup initdb

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database
sudo -u postgres createdb inventory_system
```

### Option 3: Automated Setup Script

Use the provided setup script:

```bash
# Make script executable
chmod +x scripts/setup_db.sh

# Run setup script
./scripts/setup_db.sh
```

## 🔧 Configuration

Update your `.env` file with PostgreSQL settings:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=inventory_system
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `inventory_system` |
| `DB_SSLMODE` | SSL mode (disable/require/verify-full) | `disable` |
| `DB_TIMEZONE` | Database timezone | `UTC` |

## 🚀 Running the Application

1. **Start PostgreSQL** (using one of the methods above)

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Add sample data:**
```bash
go run cmd/seed/main.go
```

4. **Start the application:**
```bash
go run main.go
```

The application will automatically:
- Connect to PostgreSQL
- Create all required tables
- Set up the admin user

## 🗄️ Database Schema

The application uses GORM for database operations and will automatically create these tables:

- `users` - User accounts and authentication
- `products` - Product catalog
- `sales` - Sales transactions
- `sale_items` - Individual sale items
- `stock_movements` - Inventory tracking
- `suppliers` - Supplier information
- `purchase_orders` - Purchase order management
- `purchase_order_items` - Purchase order line items
- `activity_logs` - System activity tracking

## 🔍 Database Management

### Using Adminer (Docker Compose)
- Access: http://localhost:8081
- Full web-based database administration
- Query editor, table browser, export/import

### Using psql Command Line
```bash
# Connect to database
psql -h localhost -U postgres -d inventory_system

# Common commands
\dt                    # List tables
\d table_name         # Describe table
SELECT * FROM users;  # Query example
\q                    # Quit
```

### Backup and Restore
```bash
# Create backup
pg_dump -h localhost -U postgres inventory_system > backup.sql

# Restore from backup
psql -h localhost -U postgres inventory_system < backup.sql
```

## 🐳 Docker Management

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f postgres

# Remove volumes (WARNING: This deletes all data)
docker-compose down -v

# Update images
docker-compose pull
docker-compose up -d
```

## 🔧 Troubleshooting

### Connection Issues
1. **Check if PostgreSQL is running:**
```bash
pg_isready -h localhost -p 5432 -U postgres
```

2. **Test database connection:**
```bash
psql -h localhost -U postgres -d inventory_system -c "SELECT version();"
```

3. **Check Docker containers:**
```bash
docker-compose ps
docker-compose logs postgres
```

### Common Errors

**"database does not exist"**
```bash
# Create the database
createdb -h localhost -U postgres inventory_system
```

**"role does not exist"**
```bash
# Create user (if needed)
createuser -h localhost -U postgres --createdb --pwprompt postgres
```

**"connection refused"**
- Check if PostgreSQL is running
- Verify host and port settings
- Check firewall settings

## 📊 Performance Tips

1. **Connection Pooling** - GORM automatically handles connection pooling
2. **Indexes** - The application creates necessary indexes automatically
3. **Query Optimization** - Use EXPLAIN ANALYZE for slow queries
4. **Regular Maintenance** - Run VACUUM and ANALYZE periodically

## 🔒 Security Considerations

For production deployments:

1. **Change default password:**
```env
DB_PASSWORD=your-secure-password-here
```

2. **Enable SSL:**
```env
DB_SSLMODE=require
```

3. **Restrict network access** - Configure PostgreSQL `pg_hba.conf`
4. **Use connection limits** - Configure `max_connections`
5. **Regular backups** - Set up automated backup schedules

## 🔄 Migration from SQLite

If you're migrating from the SQLite version:

1. **Export SQLite data** (if needed)
2. **Set up PostgreSQL** using this guide
3. **Update configuration** to use PostgreSQL
4. **Run the application** - tables will be created automatically
5. **Add sample data** using the seed command

The application will work identically with PostgreSQL, but with better performance and scalability.
