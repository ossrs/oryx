package profile

type HttpProfile struct {
	ReqMethod  string
	ReqTimeout int
	Scheme     string
	RootDomain string
	Endpoint   string
	// Deprecated, use Scheme instead
	Protocol string
	Proxy    string
}

func NewHttpProfile() *HttpProfile {
	return &HttpProfile{
		ReqMethod:  "POST",
		ReqTimeout: 60,
		Scheme:     "HTTPS",
		RootDomain: "",
		Endpoint:   "",
	}
}
