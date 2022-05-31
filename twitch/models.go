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
	ID              string `json:"id" bson:"id" example:"11897156"`
	Login           string `json:"login" bson:"login" example:"roflgator"`
	DisplayName     string `json:"display_name" bson:"display_name" example:"roflgator"`
	ProfileImageUrl string `json:"profile_image_url" bson:"profile_image_url" example:"https://static-cdn.jtvnw.net/jtv_user_pictures/2e24ced1-fec5-4cd8-858f-130f798fa023-profile_image-600x600.png"`
	TwitchURL       string `json:"twitch" bson:"twitch" example:"https://www.twitch.tv/roflgator"`
	RedditURL       string `json:"reddit" bson:"reddit" example:""`
	InstagramURL    string `json:"instagram" bson:"instagram" example:""`
	TwitterURL      string `json:"twitter" bson:"twitter" example:"https://twitter.com/roflgatorow"`
	DiscordURL      string `json:"discord" bson:"discord" example:""`
	YouTubeURL      string `json:"youtube" bson:"youtube" example:"https://www.youtube.com/channel/UCrIz-xXkVr0PSUdY4e7p-8w"`
	TikTokURL       string `json:"tiktok" bson:"tiktok" example:""`
}

type StreamerURLs struct {
	TwitchURL    string `json:"twitch" bson:"twitch" example:"https://www.twitch.tv/roflgator"`
	RedditURL    string `json:"reddit" bson:"reddit" example:""`
	InstagramURL string `json:"instagram" bson:"instagram" example:""`
	TwitterURL   string `json:"twitter" bson:"twitter" example:"https://twitter.com/roflgatorow"`
	DiscordURL   string `json:"discord" bson:"discord" example:""`
	YouTubeURL   string `json:"youtube" bson:"youtube" example:"https://www.youtube.com/channel/UCrIz-xXkVr0PSUdY4e7p-8w"`
	TikTokURL    string `json:"tiktok" bson:"tiktok" example:""`
}

type PublicStream struct {
	Status              string `json:"status" bson:"status" example:"online"`
	UserID              string `json:"user_id" bson:"user_id" example:"11897156"`
	UserLogin           string `json:"user_login" bson:"user_login" example:"roflgator"`
	UserDisplayName     string `json:"user_display_name" bson:"user_display_name" example:"roflgator"`
	UserProfileImageUrl string `json:"user_profile_image_url" bson:"user_profile_image_url" example:"https://static-cdn.jtvnw.net/jtv_user_pictures/f40e0bfe-f376-49b1-ad08-7b63f866dabb-profile_image-300x300.png"`
	StreamID            string `json:"stream_id" bson:"stream_id" example:"46365071629"`
	StreamTitle         string `json:"stream_title" bson:"stream_title" example:"GARLIC PHONE WITH THE MOST POGGERS ARTISTS, SPECIAL GUEST RUBBERROSS, OBAMA AND YOUR MOM!"`
	StreamGameID        string `json:"stream_game_id" bson:"stream_game_id" example:""`
	StreamGameName      string `json:"stream_game_name" bson:"stream_game_name" example:"Gartic Phone"`
	StreamViewerCount   string `json:"stream_viewer_count" bson:"stream_viewer_count" example:"4590"`
	StreamThumbnailUrl  string `json:"stream_thumbnail_url" bson:"stream_thumbnail_url" example:"https://static-cdn.jtvnw.net/previews-ttv/live_user_roflgator-{width}x{height}.jpg"`
	StreamStartedAt     string `json:"stream_started_at" bson:"stream_started_at" example:""`
	TwitchURL           string `json:"twitch" bson:"twitch" example:"https://www.twitch.tv/roflgator"`
	RedditURL           string `json:"reddit" bson:"reddit" example:""`
	InstagramURL        string `json:"instagram" bson:"instagram" example:""`
	TwitterURL          string `json:"twitter" bson:"twitter" example:"https://twitter.com/roflgatorow"`
	DiscordURL          string `json:"discord" bson:"discord" example:""`
	YouTubeURL          string `json:"youtube" bson:"youtube" example:"https://www.youtube.com/channel/UCrIz-xXkVr0PSUdY4e7p-8w"`
	TikTokURL           string `json:"tiktok" bson:"tiktok" example:""`
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
