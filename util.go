package mango

import "net/http"

func (u *Service) mapkey(value string) (key string, ok bool) {
	for _, v := range u.Routes {
		if v.Func == value {
			key = v.Pattern
			ok = true
			return
		}
	}
	return
}

func (u *Service) findRinfo(mount string, method func(http.ResponseWriter, *http.Request), key string, mName string) {
	for i, v := range u.Routes {
		if v.MountAt == mount && v.Func == mName {
			u.Routes[i].Handler = method

		}
	}
}
