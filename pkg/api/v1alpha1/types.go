package v1alpha1

type Client struct {
	Name        string `json:"name"`
	Id          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURL string `json:"redirect_url"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
