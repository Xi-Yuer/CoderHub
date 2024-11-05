package sqlDB

import (
	"coderhub/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	DB   *gorm.DB
	once sync.Once
)

// NewGorm initializes the database connection and returns a singleton instance of *gorm.DB
func NewGorm() *gorm.DB {
	once.Do(func() {
		// Configure custom logger
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Warn, // Log level set to Warn for better debugging
				IgnoreRecordNotFoundError: false,       // Log ErrRecordNotFound errors
				ParameterizedQueries:      false,       // Include params in SQL log for easier debugging
				Colorful:                  true,        // Enable colored output
			},
		)

		// Get DSN from environment variable or configuration file
		dsn := os.Getenv("MYSQL_DSN")
		if dsn == "" {
			dsn = "root:2214380963Wx!!@tcp(mysql:3306)/coderhub?charset=utf8&parseTime=True&loc=Local"
		}

		// Connect to MySQL using GORM
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // Data source name
			DefaultStringSize:         256,   // Default string size
			DisableDatetimePrecision:  true,  // Disable datetime precision for MySQL < 5.6
			DontSupportRenameIndex:    true,  // Drop & create index for renaming in MySQL < 5.7
			DontSupportRenameColumn:   true,  // Use 'change' for renaming columns
			SkipInitializeWithVersion: false, // Automatically configure based on MySQL version
		}), &gorm.Config{
			Logger: newLogger,
		})

		// Handle connection errors
		if err != nil {
			log.Fatalf("数据库连接失败: %v", err)
		}

		// Automatically migrate the schema
		if err := db.AutoMigrate(&model.User{}); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}

		DB = db
	})

	return DB
}
