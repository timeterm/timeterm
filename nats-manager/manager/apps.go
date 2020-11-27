package manager

type appUser struct {
	accountName string
	userName    string
}

var apps = map[string]appUser{
	"backend": {accountName: "BACKEND", userName: "backend"},
}

var appsByUsers = makeAppsByUsers()

func makeAppsByUsers() map[appUser]string {
	m := make(map[appUser]string, len(apps))
	for app, user := range apps {
		m[user] = app
	}
	return m
}

func GetAppByUser(userName, accountName string) (string, bool) {
	user, ok := appsByUsers[appUser{
		accountName: accountName,
		userName:    userName,
	}]
	return user, ok
}
