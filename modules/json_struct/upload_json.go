package json_struct

type JSONCruise struct {
	Cruises
	JSONSailor []JSONSailor
}

type JSONSailor struct {
	ID       string
	Position string
	Distance float64
}

type JSON_NEW_GUEST struct {
	First_name string
	Last_name  string
}
