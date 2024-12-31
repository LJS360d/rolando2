package repositories

import (
	stdlog "log"
	"os"
	"rolando/cmd/log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Message struct {
	ID      uint   `gorm:"primaryKey"`
	GuildID string `gorm:"index"`
	Content string `gorm:"type:text"`
}

type MessagesRepository struct {
	DB *gorm.DB
}

func NewMessagesRepository(dbPath string) (*MessagesRepository, error) {
	// Open SQLite database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.New(
			stdlog.New(os.Stdout, "\r\n", stdlog.Flags()),
			logger.Config{
				SlowThreshold: time.Second,   // Set threshold to 1 second to suppress normal slow queries
				LogLevel:      logger.Silent, // Show Info level logs (optional)
				Colorful:      true,          // Disable colored output
			},
		),
	})
	if err != nil {
		return nil, err
	}

	// Migrate the schema (creates the tables if they don't exist)
	if err := db.AutoMigrate(&Message{}); err != nil {
		return nil, err
	}

	// Set up database session optimizations
	db = db.Session(&gorm.Session{
		// Enable WAL mode for better concurrency (especially in write-heavy workloads)
		// SQLite WAL mode is more performant in multi-threaded scenarios
		NowFunc: time.Now, // Set the `Now` function to get the correct time on queries
	})

	// Ensure indexes are created for performance
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_guild_id ON messages(guild_id);").Error; err != nil {
		return nil, err
	}

	// Return the repository with the configured database connection
	return &MessagesRepository{DB: db}, nil
}

// AppendMessage inserts a new message for the guild
func (repo *MessagesRepository) AppendMessage(guildID, content string) error {
	message := Message{
		GuildID: guildID,
		Content: content,
	}

	// Use GORM to insert the message (GORM will handle the INSERT statement)
	if err := repo.DB.Create(&message).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MessagesRepository) AddMessagesToGuild(guildID string, messages []string) error {
	// Prepare a slice of Message objects
	var messageRecords []Message
	for _, content := range messages {
		messageRecords = append(messageRecords, Message{
			GuildID: guildID,
			Content: content,
		})
	}

	// Perform batch insert using CreateInBatches
	if err := repo.DB.CreateInBatches(messageRecords, 100).Error; err != nil {
		log.Log.Errorf("Error inserting messages: %v", err)
		return err
	}

	return nil
}

// GetAllGuildMessages fetches all messages for a specific guild
func (repo *MessagesRepository) GetAllGuildMessages(guildID string) ([]Message, error) {
	var messages []Message
	// Query messages for a specific guild, ordered by timestamp (default order)
	if err := repo.DB.Where("guild_id = ?", guildID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// DeleteAllGuildMessages removes all messages for a specific guild
func (repo *MessagesRepository) DeleteAllGuildMessages(guildID string) error {
	if err := repo.DB.Where("guild_id = ?", guildID).Delete(&Message{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteGuildMessage removes a message for a specific guild
func (repo *MessagesRepository) DeleteGuildMessage(guildID, content string) error {
	if err := repo.DB.Where("guild_id = ? AND content = ?", guildID, content).Delete(&Message{}).Error; err != nil {
		return err
	}
	return nil
}

// CountMessages counts the number of messages for a specific guild
func (repo *MessagesRepository) CountMessages(guildID string) (int64, error) {
	var count int64
	if err := repo.DB.Model(&Message{}).Where("guild_id = ?", guildID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
