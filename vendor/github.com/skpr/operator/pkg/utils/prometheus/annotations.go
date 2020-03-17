package prometheus

const (
	// AnnotationScrape identifies that an object has Prometheus metrics data.
	AnnotationScrape = "prometheus.io/scrape"
	// AnnotationPort tells Prometheus which port to scrape metrics.
	AnnotationPort = "prometheus.io/port"
	// AnnotationScheme tells Prometheus which http scheme to use when scraping metrics.
	AnnotationScheme = "prometheus.io/scheme"
	// AnnotationPath tells Prometheus which path to scrape metrics.
	AnnotationPath = "prometheus.io/path"
	// AnnotationToken tells Prometheus which token to use when scraping.
	AnnotationToken = "prometheus.io/token"
)

const (
	// SchemeHTTPS uses secure connections when scraping.
	SchemeHTTPS = "https"
)

const (
	// ScrapeTrue enables Prometheus scraping.
	ScrapeTrue = "true"
	// ScrapeFalse disables Prometheus scraping.
	ScrapeFalse = "false"
)
