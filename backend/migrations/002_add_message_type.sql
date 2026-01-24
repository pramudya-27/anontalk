ALTER TABLE messages ADD COLUMN type ENUM('text', 'image', 'audio', 'file') DEFAULT 'text' AFTER content;
