package boilersControlInstantView

type Boiler struct {
	IsOk          int    `json:"isOk"`
	TPod          string `json:"tPod"`
	PPod          string `json:"pPod"`
	TUlica        string `json:"tUlica"`
	TPlan         string `json:"tPlan"`
	TAlarm        string `json:"tAlarm"`
	ImageResId    int    `json:"imageResId"`
	PPodLowFixed  string `json:"pPodLowFixed"`
	PPodHighFixed string `json:"pPodHighFixed"`
	TPodFixed     string `json:"tPodFixed"`
	ID            int    `json:"id"`
	Version       int64  `json:"version"`
	LastUpdated   int64  `json:"lastUpdated"`
	ImageURL      string `json:"-"`
}
