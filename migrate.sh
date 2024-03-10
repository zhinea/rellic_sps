echo "Migrating..."
echo "Creating the database..."
sudo mariadb -e "CREATE DATABASE rellic_m_elysia;"

echo "Database created, now migrating the schema..."
sudo mariadb rellic_m_elysia < db_migration.sql

echo "Schema migrated, now migrating the data..."