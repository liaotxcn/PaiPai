package xerr

// User err
var (
	PhoneNotFound = New(USER_ERROR, "phone not found")
	IdNotFound    = New(USER_ERROR, "id not found")
	UserPwdErr    = New(USER_ERROR, "password is wrong")
	ParamError    = New(REQUEST_PARAM_ERROR, "params error")
)

// Friend Err
var (
	FriendAlreadyExists    = New(FRIEND_ERROR, "friend already exists")
	FriendRequestOnPending = New(FRIEND_ERROR, "friend request on pending")
	FriendRequestRefused   = New(FRIEND_ERROR, "friend request refused")
	FriendListNotFound     = New(FRIEND_ERROR, "friend list not found")
	FriendReqListNotFound  = New(FRIEND_ERROR, "friend request list not found")
	FindFriendByIdErr      = New(FRIEND_ERROR, "find friend by id error")
)

// Group Err
var (
	GroupNotFound        = New(GROUP_ERROR, "group not found")
	GroupPutInNotFound   = New(GROUP_ERROR, "group put in request not found")
	GroupInviterNotFound = New(GROUP_ERROR, "group inviter not found")
	FindGroupByIdErr     = New(GROUP_ERROR, "find group by id error, user haven't attend in any group")
)
