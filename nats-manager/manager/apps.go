package manager

type appUser struct {
	accountName string
	userName    string
}

func makeAppsByUsers(apps map[string]appUser) map[appUser]string {
	m := make(map[appUser]string, len(apps))
	for app, user := range apps {
		m[user] = app
	}
	return m
}

type appIndex struct {
	apps map[appUser]string
}

func newAppIndex() appIndex {
	return appIndex{
		apps: makeAppsByUsers(map[string]appUser{
			"backend":        {accountName: "BACKEND", userName: "backend"},
			"backend-emdevs": {accountName: "EMDEVS", userName: "backend"},
		}),
	}
}

func (i appIndex) getAppByUser(userName, accountName string) (string, bool) {
	user, ok := i.apps[appUser{
		accountName: accountName,
		userName:    userName,
	}]
	return user, ok
}
