package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
	Email    EmailConfig
	SMS      SMSConfig
	Payment  PaymentConfig
	GST      GSTConfig
	OAuth    OAuthConfig
	CORS     CORSConfig
	RateLimit RateLimitConfig
	Security SecurityConfig
	Business BusinessConfig
	Logging  LoggingConfig
	Features FeatureFlags
	Admin    AdminConfig
}

type ServerConfig struct {
	Port    string
	Env     string
	Name    string
	Version string
}

type DatabaseConfig struct {
	URL            string
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	SSLMode        string
	MaxConnections int
	MaxIdle        int
	MaxLifetime    time.Duration
}

type RedisConfig struct {
	URL      string
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret              string
	Expiry              time.Duration
	RefreshTokenExpiry  time.Duration
	Algorithm           string
}

type StorageConfig struct {
	UploadDir         string
	MaxUploadSize     int64
	AllowedFileTypes  string
	UseS3             bool
	S3Bucket          string
	S3Region          string
	AWSAccessKey      string
	AWSSecretKey      string
	S3Endpoint        string
	S3UsePathStyle    bool
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type SMSConfig struct {
	TwilioAccountSID       string
	TwilioAuthToken        string
	TwilioPhoneNumber      string
	TwilioMessagingService string
}

type PaymentConfig struct {
	RazorpayKeyID        string
	RazorpayKeySecret    string
	RazorpayWebhookSecret string
	RazorpayMode         string
	StripeKey            string
	StripeSecret         string
	StripeWebhookSecret  string
}

type GSTConfig struct {
	APIKey     string
	APIURL     string
	APIEnabled bool
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

type CORSConfig struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

type RateLimitConfig struct {
	Enabled         bool
	Requests        int
	Duration        time.Duration
	AuthRequests    int
	AuthDuration    time.Duration
}

type SecurityConfig struct {
	BcryptCost         int
	OTPLength          int
	OTPExpiry          time.Duration
	PasswordResetExpiry time.Duration
	EmailVerifyExpiry   time.Duration
}

type BusinessConfig struct {
	MaxCartItems          int
	MaxQuantityPerItem    int
	CartExpiryDays        int
	OrderAutoCancelMinutes int
	VendorApprovalDays    int
	ReturnWindowDays      int
	CODMaxAmount          int
	CODFee                int
}

type LoggingConfig struct {
	Level  string
	Format string
	Output string
}

type FeatureFlags struct {
	GoogleAuth       bool
	COD              bool
	Wallet           bool
	BulkUpload       bool
	GSTVerification  bool
}

type AdminConfig struct {
	Email    string
	Password string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		// If .env doesn't exist, try to continue with environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port:    viper.GetString("PORT"),
			Env:     viper.GetString("ENV"),
			Name:    viper.GetString("APP_NAME"),
			Version: viper.GetString("APP_VERSION"),
		},
		Database: DatabaseConfig{
			URL:            viper.GetString("DATABASE_URL"),
			Host:           viper.GetString("DB_HOST"),
			Port:           viper.GetString("DB_PORT"),
			User:           viper.GetString("DB_USER"),
			Password:       viper.GetString("DB_PASSWORD"),
			Name:           viper.GetString("DB_NAME"),
			SSLMode:        viper.GetString("DB_SSLMODE"),
			MaxConnections: viper.GetInt("DB_MAX_CONNECTIONS"),
			MaxIdle:        viper.GetInt("DB_MAX_IDLE"),
			MaxLifetime:    viper.GetDuration("DB_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			URL:      viper.GetString("REDIS_URL"),
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("JWT_SECRET"),
			Expiry:             viper.GetDuration("JWT_EXPIRY"),
			RefreshTokenExpiry: viper.GetDuration("REFRESH_TOKEN_EXPIRY"),
			Algorithm:          viper.GetString("JWT_ALGORITHM"),
		},
		Storage: StorageConfig{
			UploadDir:        viper.GetString("UPLOAD_DIR"),
			MaxUploadSize:    viper.GetInt64("MAX_UPLOAD_SIZE"),
			AllowedFileTypes: viper.GetString("ALLOWED_FILE_TYPES"),
			UseS3:            viper.GetBool("USE_S3"),
			S3Bucket:         viper.GetString("S3_BUCKET"),
			S3Region:         viper.GetString("S3_REGION"),
			AWSAccessKey:     viper.GetString("AWS_ACCESS_KEY_ID"),
			AWSSecretKey:     viper.GetString("AWS_SECRET_ACCESS_KEY"),
			S3Endpoint:       viper.GetString("S3_ENDPOINT"),
			S3UsePathStyle:   viper.GetBool("S3_USE_PATH_STYLE"),
		},
		Email: EmailConfig{
			SMTPHost:     viper.GetString("SMTP_HOST"),
			SMTPPort:     viper.GetInt("SMTP_PORT"),
			SMTPUser:     viper.GetString("SMTP_USER"),
			SMTPPassword: viper.GetString("SMTP_PASSWORD"),
			FromEmail:    viper.GetString("EMAIL_FROM"),
			FromName:     viper.GetString("EMAIL_FROM_NAME"),
		},
		SMS: SMSConfig{
			TwilioAccountSID:       viper.GetString("TWILIO_ACCOUNT_SID"),
			TwilioAuthToken:        viper.GetString("TWILIO_AUTH_TOKEN"),
			TwilioPhoneNumber:      viper.GetString("TWILIO_PHONE_NUMBER"),
			TwilioMessagingService: viper.GetString("TWILIO_MESSAGING_SERVICE_SID"),
		},
		Payment: PaymentConfig{
			RazorpayKeyID:         viper.GetString("RAZORPAY_KEY_ID"),
			RazorpayKeySecret:     viper.GetString("RAZORPAY_KEY_SECRET"),
			RazorpayWebhookSecret: viper.GetString("RAZORPAY_WEBHOOK_SECRET"),
			RazorpayMode:          viper.GetString("RAZORPAY_MODE"),
			StripeKey:             viper.GetString("STRIPE_KEY"),
			StripeSecret:          viper.GetString("STRIPE_SECRET"),
			StripeWebhookSecret:   viper.GetString("STRIPE_WEBHOOK_SECRET"),
		},
		GST: GSTConfig{
			APIKey:     viper.GetString("GST_API_KEY"),
			APIURL:     viper.GetString("GST_API_URL"),
			APIEnabled: viper.GetBool("GST_API_ENABLED"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
			GoogleRedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetString("ALLOWED_ORIGINS"),
			AllowedMethods: viper.GetString("ALLOWED_METHODS"),
			AllowedHeaders: viper.GetString("ALLOWED_HEADERS"),
		},
		RateLimit: RateLimitConfig{
			Enabled:      viper.GetBool("RATE_LIMIT_ENABLED"),
			Requests:     viper.GetInt("RATE_LIMIT_REQUESTS"),
			Duration:     viper.GetDuration("RATE_LIMIT_DURATION"),
			AuthRequests: viper.GetInt("RATE_LIMIT_AUTH_REQUESTS"),
			AuthDuration: viper.GetDuration("RATE_LIMIT_AUTH_DURATION"),
		},
		Security: SecurityConfig{
			BcryptCost:          viper.GetInt("BCRYPT_COST"),
			OTPLength:           viper.GetInt("OTP_LENGTH"),
			OTPExpiry:           viper.GetDuration("OTP_EXPIRY"),
			PasswordResetExpiry: viper.GetDuration("PASSWORD_RESET_EXPIRY"),
			EmailVerifyExpiry:   viper.GetDuration("EMAIL_VERIFY_EXPIRY"),
		},
		Business: BusinessConfig{
			MaxCartItems:           viper.GetInt("MAX_CART_ITEMS"),
			MaxQuantityPerItem:     viper.GetInt("MAX_QUANTITY_PER_ITEM"),
			CartExpiryDays:         viper.GetInt("CART_EXPIRY_DAYS"),
			OrderAutoCancelMinutes: viper.GetInt("ORDER_AUTO_CANCEL_MINUTES"),
			VendorApprovalDays:     viper.GetInt("VENDOR_APPROVAL_DAYS"),
			ReturnWindowDays:       viper.GetInt("RETURN_WINDOW_DAYS"),
			CODMaxAmount:           viper.GetInt("COD_MAX_AMOUNT"),
			CODFee:                 viper.GetInt("COD_FEE"),
		},
		Logging: LoggingConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
			Output: viper.GetString("LOG_OUTPUT"),
		},
		Features: FeatureFlags{
			GoogleAuth:      viper.GetBool("FEATURE_GOOGLE_AUTH"),
			COD:             viper.GetBool("FEATURE_COD"),
			Wallet:          viper.GetBool("FEATURE_WALLET"),
			BulkUpload:      viper.GetBool("FEATURE_BULK_UPLOAD"),
			GSTVerification: viper.GetBool("FEATURE_GST_VERIFICATION"),
		},
		Admin: AdminConfig{
			Email:    viper.GetString("ADMIN_EMAIL"),
			Password: viper.GetString("ADMIN_PASSWORD"),
		},
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// Server
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("APP_NAME", "OQart Backend")
	viper.SetDefault("APP_VERSION", "1.0.0")

	// Database
	viper.SetDefault("DB_MAX_CONNECTIONS", 25)
	viper.SetDefault("DB_MAX_IDLE", 5)
	viper.SetDefault("DB_MAX_LIFETIME", "5m")

	// JWT
	viper.SetDefault("JWT_EXPIRY", "1h")
	viper.SetDefault("REFRESH_TOKEN_EXPIRY", "720h")
	viper.SetDefault("JWT_ALGORITHM", "HS256")

	// Storage
	viper.SetDefault("UPLOAD_DIR", "./uploads")
	viper.SetDefault("MAX_UPLOAD_SIZE", 5242880)
	viper.SetDefault("USE_S3", false)

	// Rate Limit
	viper.SetDefault("RATE_LIMIT_ENABLED", true)
	viper.SetDefault("RATE_LIMIT_REQUESTS", 1000)
	viper.SetDefault("RATE_LIMIT_DURATION", "1h")
	viper.SetDefault("RATE_LIMIT_AUTH_REQUESTS", 10)
	viper.SetDefault("RATE_LIMIT_AUTH_DURATION", "1m")

	// Security
	viper.SetDefault("BCRYPT_COST", 12)
	viper.SetDefault("OTP_LENGTH", 6)
	viper.SetDefault("OTP_EXPIRY", "5m")
	viper.SetDefault("PASSWORD_RESET_EXPIRY", "30m")
	viper.SetDefault("EMAIL_VERIFY_EXPIRY", "24h")

	// Business
	viper.SetDefault("MAX_CART_ITEMS", 20)
	viper.SetDefault("MAX_QUANTITY_PER_ITEM", 10)
	viper.SetDefault("CART_EXPIRY_DAYS", 7)
	viper.SetDefault("ORDER_AUTO_CANCEL_MINUTES", 30)
	viper.SetDefault("VENDOR_APPROVAL_DAYS", 3)
	viper.SetDefault("RETURN_WINDOW_DAYS", 15)
	viper.SetDefault("COD_MAX_AMOUNT", 5000)
	viper.SetDefault("COD_FEE", 50)

	// Logging
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("LOG_OUTPUT", "stdout")

	// Features
	viper.SetDefault("FEATURE_GOOGLE_AUTH", true)
	viper.SetDefault("FEATURE_COD", true)
	viper.SetDefault("FEATURE_WALLET", false)
	viper.SetDefault("FEATURE_BULK_UPLOAD", true)
	viper.SetDefault("FEATURE_GST_VERIFICATION", false)
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	if c.Database.URL == "" && c.Database.Host == "" {
		return fmt.Errorf("DATABASE_URL or DB_HOST is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}
