package users

import (
	"encoding/json"
)

func (users Users) Serialize(authType AuthType) []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Serialize(authType)
	}
	return result
}

func (user *User) Serialize(authType AuthType) interface{} {
	if authType == Admin {
		return UserPayload[PrivateUser]{Id: user.ID, Type: "users", Attributes: user.AsPrivate()}
	} else {
		return UserPayload[PublicUser]{Id: user.ID, Type: "users", Attributes: user.AsPublic()}
	}
}

func (user *User) AsPublic() PublicUser {
	userJson, _ := json.Marshal(user)
	var publicUser PublicUser
	json.Unmarshal(userJson, &publicUser)
	return publicUser
}

func (user *User) AsPrivate() PrivateUser {
	userJson, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(userJson, &privateUser)
	return privateUser
}
