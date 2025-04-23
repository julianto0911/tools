package tools

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManager struct {
	connections map[string]*gorm.DB
	mu          sync.RWMutex
	config      DBConfiguration
}

func NewDBManager(defaultConfig DBConfiguration) *DBManager {
	return &DBManager{
		connections: make(map[string]*gorm.DB),
		config:      defaultConfig,
	}
}

func (m *DBManager) GetConnection(dbname string) (*gorm.DB, error) {
	m.mu.RLock()
	db, exists := m.connections[dbname]
	m.mu.RUnlock()

	if exists {
		return db, nil
	}

	return m.createConnection(dbname)
}

func (m *DBManager) createConnection(dbname string) (*gorm.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check if connection was created while waiting for lock
	if db, exists := m.connections[dbname]; exists {
		return db, nil
	}

	// Clone base config and modify for tenant
	tenantConfig := m.config
	tenantConfig.DBName = fmt.Sprintf("%s_%s", m.config.DBName, dbname)

	// Create new connection
	sqlConn, err := ConnectDB(tenantConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tenant database: %w", err)
	}

	// Initialize GORM
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlConn,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM: %w", err)
	}

	m.connections[dbname] = db
	return db, nil
}

func (m *DBManager) CloseConnections() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for tenantID, db := range m.connections {
		sqlDB, err := db.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("tenant %s: %w", tenantID, err))
			continue
		}
		if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("tenant %s: %w", tenantID, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}
	return nil
}
