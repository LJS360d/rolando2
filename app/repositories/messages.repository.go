package repositories

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema (creates the tables if they don't exist)
	if err := db.AutoMigrate(&Message{}); err != nil {
		return nil, err
	}

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

// ReadAllMessages fetches all messages for a specific guild
func (repo *MessagesRepository) ReadAllMessages(guildID string) ([]Message, error) {
	var messages []Message
	// Query messages for a specific guild, ordered by timestamp (default order)
	if err := repo.DB.Where("guild_id = ?", guildID).Order("created_at asc").Find(&messages).Error; err != nil {
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
