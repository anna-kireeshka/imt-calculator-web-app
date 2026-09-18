package fitness

import "time"

func age(dob, now time.Time) int {
	y, m, d := now.Date()
	by, bm, bd := dob.Date()

	age := y - by
	if m < bm || (m == bm && d < bd) {
		age--
	}
	return age
}

func AgeAt(dob time.Time, loc *time.Location) int {
	return age(dob, time.Now().In(loc))
}
