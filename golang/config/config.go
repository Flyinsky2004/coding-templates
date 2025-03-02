package config

// 项目全局参数
var (
	// Gin配置
	GinPort          = ":8080"
	CorsAllowOrigins = []string{"http://localhost:5173"}
	// mysql配置
	MysqlHost     = "localhost"
	MysqlPort     = "3306"
	MysqlUser     = "root"
	MysqlPassword = "root"
	MysqlDatabase = "ideacosmos"

	// redis配置
	RedisHost     = "localhost"
	RedisPort     = "6379"
	RedisPassword = "root"
	RedisDatabase = 0

	// JWT配置
	JwtSecretKey  = "5771DF8B-A018-4576-92E2-988E1AFF2909"
	JwtExpireTime = 7 * 24 // 过期时间（小时）

	// 邮件配置
	MailSenderName = "Flyinsky"
	MailSmtpHost   = "smtp.office365.com"
	MailUsername   = "test@test.onmicrosoft.com"
	MailPassword   = "password"
	MailSmtpPort   = 587

	// Azure TTS 配置
	AzureTTSBaseURL      = "https://%s.tts.speech.microsoft.com/cognitiveservices/v1"
	AzureTTSSSMLTemplate = `
<speak version='1.0' xml:lang='%s'>
    <voice xml:lang='%s' xml:gender='%s' name='%s'>
        %s
    </voice>
</speak>`
	AzureTTSSubscriptionKey = "your-subscription-key" // 请替换为实际的订阅密钥
	AzureTTSRegion          = "eastasia"              // 请替换为实际的区域

	// OpenAI 配置
	OpenAIKey         = "sk-myOpenAiKey"
	OpenAIBaseURL     = "https://api.openai.com/v1"
	GlobalTemperature = 0.5
	ThinkModelName    = "deepseek-r1-250120"
	AgentModelName    = "deepseek-r1-250120"
	UseModelName      = "deepseek-v3-241226"
)
