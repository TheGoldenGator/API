package twitch

import "time"

// Get array of users from endpoint
type ManyUsers struct {
	Users []User `json:"data"`
}

// User represents a Twitch User
type User struct {
	ID              string    `json:"id" example:"237509153"`
	Login           string    `json:"login" example:"mahcksimus"`
	DisplayName     string    `json:"display_name" example:"Mahcksimus"`
	Type            string    `json:"type" example:""`
	BroadcasterType string    `json:"broadcaster_type" example:""`
	Description     string    `json:"description" example:"I chat and program."`
	ProfileImageURL string    `json:"profile_image_url" example:"https://static-cdn.jtvnw.net/jtv_user_pictures/41236a31-635c-4bee-ba3e-dc791371a746-profile_image-300x300.png"`
	OfflineImageURL string    `json:"offline_image_url" example:""`
	ViewCount       int       `json:"view_count" example:"200"`
	Email           string    `json:"email" example:""`
	CreatedAt       time.Time `json:"created_at" example:"2018-07-10T02:16:03Z"`
}

type Streamer struct {
	ID              int    `json:"id" bson:"id"`
	Login           string `json:"login" bson:"login"`
	DisplayName     string `json:"display_name" bson:"display_name"`
	ProfileImageUrl string `json:"profile_image_url" bson:"profile_image_url"`
}

type PublicStream struct {
	Status              string `json:"status" bson:"status"`
	UserID              int    `json:"user_id" bson:"user_id"`
	UserLogin           string `json:"user_login" bson:"user_login"`
	UserDisplayName     string `json:"user_display_name" bson:"user_display_name"`
	UserProfileImageUrl string `json:"user_profile_image_url" bson:"user_profile_image_url"`
	StreamID            string `json:"stream_id" bson:"stream_id"`
	StreamTitle         string `json:"stream_title" bson:"stream_title"`
	StreamGameID        string `json:"stream_game_id" bson:"stream_game_id"`
	StreamGameName      string `json:"stream_game_name" bson:"stream_game_name"`
	StreamViewerCount   int    `json:"stream_viewer_count" bson:"stream_viewer_count"`
	StreamThumbnailUrl  string `json:"stream_thumbnail_url" bson:"stream_thumbnail_url"`
}

/* EventSub */
// Represents a subscription
type EventSubSubscription struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Status    string            `json:"status"`
	Condition EventSubCondition `json:"condition"`
	Transport EventSubTransport `json:"transport"`
	CreatedAt time.Time         `json:"created_at"`
	Cost      int               `json:"cost"`
}

type EventSubCondition struct {
	BroadcasterUserID     string `json:"broadcaster_user_id"`
	FromBroadcasterUserID string `json:"from_broadcaster_user_id"`
	ToBroadcasterUserID   string `json:"to_broadcaster_user_id"`
	RewardID              string `json:"reward_id"`
	ClientID              string `json:"client_id"`
	ExtensionClientID     string `json:"extension_client_id"`
	UserID                string `json:"user_id"`
}

// Transport for the subscription, currently the only supported Method is "webhook". Secret must be between 10 and 100 characters
type EventSubTransport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

// Data for a stream online notification
type EventSubStreamOnlineEvent struct {
	ID                   string    `json:"id"`
	BroadcasterUserID    string    `json:"broadcaster_user_id"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login"`
	BroadcasterUserName  string    `json:"broadcaster_user_name"`
	Type                 string    `json:"type"`
	StartedAt            time.Time `json:"started_at"`
}

// Data for a stream offline notification
type EventSubStreamOfflineEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}

// Data for a channel update notification
type EventSubChannelUpdateEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Title                string `json:"title"`
	Language             string `json:"language"`
	CategoryID           string `json:"category_id"`
	CategoryName         string `json:"category_name"`
	IsMature             bool   `json:"is_mature"`
}

/* Streams */
type Stream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	TagIDs       []string  `json:"tag_ids"`
	IsMature     bool      `json:"is_mature"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

type ManyStreams struct {
	Streams    []Stream   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

/* Pagination */
type Pagination struct {
	Cursor string `json:"cursor"`
}
