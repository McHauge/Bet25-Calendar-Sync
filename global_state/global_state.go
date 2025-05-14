package global_state

// Global Vars
var (
	Debug = false
)

type State struct {
	// Gets set above
	Debug bool

	// Gets set by .env and config file
	DEBUG              string
	UPDATE_INTERVAL    string
	GOOGLE_CALENDAR_ID string
	GOOGLE_CREDENTIALS string
	GOOGLE_API_KEY     string

	// Rebook Params
	REBOK_USERNAME string
	REBOK_PASSWORD string
	REBOK_GET_USER string
}

func NewState() State {
	return State{
		Debug: Debug,
	}
}
