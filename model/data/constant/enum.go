package constant

type Authority string

const (
	ROLE_STUDENT         Authority = "ROLE_STUDENT"
	ROLE_STUDENT_COUNCIL Authority = "ROLE_STUDENT_COUNCIL"
)

type Gender string

const (
	MAN   Gender = "MAN"
	WOMAN Gender = "WOMAN"
)

type Major string

const (
	SW_DEVELOP Major = "SW_DEVELOP"
	SMART_IOT  Major = "SMART_IOT"
	AI         Major = "AI"
)
