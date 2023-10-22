package types

type OperationUserInfo struct {
	EffectiveUserInfo *UserInfo
	RealUserInfo      *UserInfo
	// IsSameUser indicate that operation user is the same as real user
	IsSameUser bool
	// IsSuperUser indicate that operation user is "root" or a user that belong to a super-user group
	IsSuperUser bool
	// OperatingAsSuperUser indicates that operation user is "root" or a user runs command with "sudo"
	OperatingAsSuperUser bool
}

// GetDefaultWorkingUser returns:
//
// - Effective user if is not super-user and not operating as super-user
//
// - Real user otherwise
func (o *OperationUserInfo) GetDefaultWorkingUser() *UserInfo {
	if o.IsSameUser {
		return o.RealUserInfo
	} else if !o.IsSuperUser && !o.OperatingAsSuperUser {
		return o.EffectiveUserInfo
	} else {
		return o.RealUserInfo
	}
}
