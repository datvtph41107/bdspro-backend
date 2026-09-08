package constants

const (
	// chat
	MESSAGE_TYPE                  = "message"
	CREATE_ROOM_TYPE              = "create_room"
	JOIN_ROOM_TYPE                = "join_room"  // unimplemented
	LEAVE_ROOM_TYPE               = "leave_room" // unimplemented
	READ_TYPE                     = "read"
	RECALL_MESSAGE_TYPE           = "recall_message"
	EDIT_MESSAGE_TYPE             = "edit_message" // unimplemented
	PIN_MESSAGE_TYPE              = "pin_message"
	REACT_MESSAGE_TYPE            = "react_message"
	SEEN_TYPE                     = "seen"            // duplicate
	UPDATE_SETTINGS_TYPE          = "update_settings" // unimplemented
	DELETE_ROOM_TYPE              = "delete_room"     // unimplemented
	REMOVE_MEMBER_TYPE            = "remove_member"
	DELETE_VIOLATION_MESSAGE_TYPE = "delete_violation_message"
	UPDATE_ROOM_TYPE              = "update_room"
	REACT_MESSAGE                 = "react_message" // unimplemented
	SYSTEM_MESSAGE_TYPE           = "system"

	// notification
	NOTIFICATION_TYPE = "notification"

	// app state
	APP_STATE_TYPE = "app_state"

	// typing — không gửi push notification (xử lý riêng trong sendNotificationForMessage)
	TYPING_INDICATOR_TYPE = "typing"
)
