package analytics

var (
	Client AnalyticsService
)

type AnalyticsService interface {
	Close() error
}

func Init() error {
	client := NewPosthogAnalyticsService()
	err := client.Init()
	if err != nil {
		return err
	}
	Client = client
	return nil
}

func Close() error {
	return Client.Close()
}
