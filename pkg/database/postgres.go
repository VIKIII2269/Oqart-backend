package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/oqart/backend/config"
	"github.com/oqart/backend/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Database struct {
	DB     *gorm.DB
	SqlDB  *sql.DB
	Config *config.DatabaseConfig
	Logger *logger.Logger
}

func NewPostgres(cfg *config.DatabaseConfig, log *logger.Logger) (*Database, error) {
	// Build DSN if not provided
	dsn := cfg.URL
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	}

	// Configure GORM logger
	var gormLogLevel gormlogger.LogLevel
	switch log.GetZerolog().GetLevel() {
	case -1: // trace/debug
		gormLogLevel = gormlogger.Info
	default:
		gormLogLevel = gormlogger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	// Ping database to verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL database",
		"host", cfg.Host,
		"database", cfg.Name,
		"max_connections", cfg.MaxConnections)

	return &Database{
		DB:     db,
		SqlDB:  sqlDB,
		Config: cfg,
		Logger: log,
	}, nil
}

func (d *Database) Close() error {
	d.Logger.Info("Closing database connection")
	return d.SqlDB.Close()
}

func (d *Database) Ping() error {
	return d.SqlDB.Ping()
}

func (d *Database) Stats() sql.DBStats {
	return d.SqlDB.Stats()
}

func (d *Database) Begin() *gorm.DB {
	return d.DB.Begin()
}

func (d *Database) Transaction(fn func(*gorm.DB) error) error {
	return d.DB.Transaction(fn)
}

// Health check
func (d *Database) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.SqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}
