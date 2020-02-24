package models

// EmailNotifications type
type EmailNotifications struct {
	UnsubscribeFromAll bool `bson:"unsubscribe_from_all" json:"unsubscribe_from_all"`
}

// PushNotifications type
type PushNotifications struct {
	UnsubscribeFromAll bool `bson:"unsubscribe_from_all" json:"unsubscribe_from_all"`
}

// Preferences type
type Preferences struct {
	Language           string             `bson:"language" json:"language"`
	EmailNotifications EmailNotifications `bson:"email_notifications" json:"email_notifications"`
	PushNotifications  PushNotifications  `bson:"push_notification" json:"push_notification"`
}

// ProfileName type
type ProfileName struct {
	FamilyName string `bson:"family_name" json:"family_name"`
	MiddleName string `bson:"middle_name" json:"middle_name"`
	GivenName  string `bson:"given_name" json:"given_name"`
}
