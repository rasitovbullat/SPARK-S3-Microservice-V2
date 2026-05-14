// Package domain contains business entities and error definitions.
package domain

import (
	"errors"
	"strings"
)

// FileType represents the type of file being uploaded.
type FileType string

const (
	FileTypeAvatar  FileType = "avatar"
	FileTypeMessage FileType = "message"
	FileTypeSticker FileType = "sticker"
)

// String returns the string representation of FileType.
func (ft FileType) String() string {
	return string(ft)
}

// AllowedExtensions defines valid extensions for each file type.
var AllowedExtensions = map[FileType][]string{
	FileTypeAvatar:  {"jpg", "jpeg", "png", "webp", "gif"},
	FileTypeMessage: {"jpg", "jpeg", "png", "webp", "gif", "pdf", "mp3", "wav", "ogg", "m4a", "mp4", "webm"},
	FileTypeSticker: {"png", "gif", "webp"},
}

// ContentTypes maps file extensions to MIME types.
var ContentTypes = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"pdf":  "application/pdf",
	"mp3":  "audio/mpeg",
	"wav":  "audio/wav",
	"ogg":  "audio/ogg",
	"m4a":  "audio/mp4",
	"mp4":  "video/mp4",
	"webm": "video/webm",
}

// IsMediaAudio checks if extension is audio format.
func IsMediaAudio(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	return ext == "mp3" || ext == "wav" || ext == "ogg" || ext == "m4a"
}

// IsMediaVideo checks if extension is video format.
func IsMediaVideo(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	return ext == "mp4" || ext == "webm"
}

// Domain errors.
var (
	ErrInvalidFileType    = errors.New("invalid file type")
	ErrInvalidExtension   = errors.New("invalid file extension")
	ErrInvalidFileID      = errors.New("invalid file ID format")
	ErrForbidden          = errors.New("access denied")
	ErrFileNotFound       = errors.New("file not found")
	ErrInternalError      = errors.New("internal server error")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// IsValidExtension checks if the extension is allowed for the file type.
func IsValidExtension(fileType FileType, ext string) bool {
	allowed, ok := AllowedExtensions[fileType]
	if !ok {
		return false
	}

	ext = strings.ToLower(strings.TrimPrefix(ext, "."))

	for _, allowedExt := range allowed {
		if allowedExt == ext {
			return true
		}
	}

	return false
}

// GetContentType returns the MIME type for a given extension.
func GetContentType(ext string) string {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	if ct, ok := ContentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

// ParseFileType converts a string to FileType.
func ParseFileType(s string) (FileType, error) {
	switch strings.ToLower(s) {
	case "avatar":
		return FileTypeAvatar, nil
	case "message":
		return FileTypeMessage, nil
	case "sticker":
		return FileTypeSticker, nil
	default:
		return "", ErrInvalidFileType
	}
}

// ExtractOwnerFromFileID extracts the user ID from a file_id.
// file_id format: {type}_{uid}_{timestamp}_{uuid}.{ext}
// Example: avatar_user123_1700000000_a1b2c3d4.jpg
func ExtractOwnerFromFileID(fileID string) (string, error) {
	parts := strings.Split(fileID, "_")
	// Minimum: type_uid_timestamp_uuid.ext = 4 parts
	if len(parts) < 4 {
		return "", ErrInvalidFileID
	}
	// Owner is the second part (index 1)
	return parts[1], nil
}

// ExtractFileTypeFromFileID extracts the file type from a file_id.
func ExtractFileTypeFromFileID(fileID string) (FileType, error) {
	parts := strings.Split(fileID, "_")
	if len(parts) < 4 {
		return "", ErrInvalidFileID
	}
	return ParseFileType(parts[0])
}

// ExtractExtensionFromFileID extracts the extension from a file_id.
func ExtractExtensionFromFileID(fileID string) (string, error) {
	parts := strings.Split(fileID, ".")
	if len(parts) < 2 {
		return "", ErrInvalidFileID
	}
	return parts[len(parts)-1], nil
}

// File represents a file entity with metadata.
type File struct {
	ID     string
	Type   FileType
	UserID string
	Ext    string
}

// BuildS3Key constructs the S3 object key from file metadata.
func (f *File) BuildS3Key() string {
	return f.Type.String() + "/" + f.ID
}

// FileIDToS3Key converts a file_id to S3 object key.
// file_id format: {type}_{uid}_{timestamp}_{uuid}.{ext}
func FileIDToS3Key(fileID string) (string, error) {
	fileType, err := ExtractFileTypeFromFileID(fileID)
	if err != nil {
		return "", err
	}
	return fileType.String() + "/" + fileID, nil
}

// IsFileID checks if the given string is a valid file_id format.
func IsFileID(url string) bool {
	if url == "" {
		return false
	}
	// file_id format: {type}_{uid}_{timestamp}_{short_uuid}.{ext}
	// Examples: avatar_user123_1700000000_abc123.png, message_user456_1700000000_def456.jpg
	parts := strings.Split(url, "_")
	if len(parts) < 3 {
		return false
	}
	// Check last part has extension
	lastPart := parts[len(parts)-1]
	if !strings.Contains(lastPart, ".") {
		return false
	}
	ext := strings.Split(lastPart, ".")[1]
	if len(ext) < 2 || len(ext) > 4 {
		return false
	}
	return true
}

// IsCloudinaryURL checks if URL is from Cloudinary (old format).
func IsCloudinaryURL(url string) bool {
	return strings.Contains(url, "cloudinary.com") || strings.Contains(url, "res.cloudinary.com")
}