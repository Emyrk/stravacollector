package strava

//type AppClient struct {
//	ClientID     string
//	ClientSecret string
//}
//
//func NewAppClient() *AppClient {
//
//}
//
//func (c *AppClient) Subscription(ctx context.Context) (any, error) {
//	formData := url.Values{}
//
//	resp, err := c.Request(ctx, http.MethodPost, "https://www.strava.com/api/v3/push_subscriptions", formData.Encode(), nil)
//	if err != nil {
//		return Athlete{}, fmt.Errorf("request: %w", err)
//	}
//
//	var athlete Athlete
//	return athlete, c.DecodeResponse(resp, &athlete, http.StatusOK)
//}
