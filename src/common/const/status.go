package status

type ErrorStatus struct {
	StatusCode int    `json:"statusCd"`
	StatusName string `json:"statusName"`
}

var ErrorStatusMap = map[string]ErrorStatus{
	"OK": {
		StatusCode: 200,
		StatusName: "OK",
	},
	"CREATED": {
		StatusCode: 201,
		StatusName: "Created",
	},
	"ACCEPTED": {
		StatusCode: 202,
		StatusName: "Accepted",
	},
	"BAD_REQUEST": {
		StatusCode: 400,
		StatusName: "Bad Request",
	},
	"ENABLE_CHECK_ERROR": {
		StatusCode: 400,
		StatusName: "Enable Check Error",
	},
	"UNAUTHORIZED": {
		StatusCode: 401,
		StatusName: "Unauthorized",
	},
	"FORBIDDEN": {
		StatusCode: 403,
		StatusName: "Forbidden",
	},
	"NOT_FOUND": {
		StatusCode: 404,
		StatusName: "Not Found",
	},
	"INTERNAL_SERVER_ERROR": {
		StatusCode: 500,
		StatusName: "Internal Server Error",
	},
}
