package contracts

// Container defines the interface for the application container.
type Container interface {
	Bind(name string, resolver func() any)
	BindInstance(name string, instance any)
	Singleton(name string, resolver func() any)
	Make(name string) (any, error)
	Has(name string) bool
}

// EventDispatcher defines the interface for event dispatching.
type EventDispatcher interface {
	Listen(event any, listener func(any))
	Dispatch(event any)
	HasListeners(event any) bool
}

// Mailer defines the interface for sending mail.
type Mailer interface {
	Send(mailable any) error
	Queue(mailable any) error
}

// Queue defines the interface for queue operations.
type Queue interface {
	Push(job any) error
	Size() int
}

// CacheStore defines the interface for cache operations.
type CacheStore interface {
	Get(key string) (any, bool)
	Put(key string, value any, ttl ...int) bool
	Forget(key string) bool
	Has(key string) bool
	Flush() bool
}

// Session defines the interface for session operations.
type Session interface {
	Get(key string) any
	Put(key string, value any)
	Forget(key string)
	Has(key string) bool
	All() map[string]any
	Flush()
}

// Translator defines the interface for translation operations.
type Translator interface {
	Translate(key string, replaces ...map[string]string) string
	GetLocale() string
	SetLocale(locale string)
}

// Dispatcher defines the interface for mail dispatching.
type Dispatcher interface {
	Send(mailable any) error
}

// Hasher defines the interface for password hashing.
type Hasher interface {
	Make(value string) (string, error)
	Check(value, hashed string) bool
	NeedsRehash(hashed string) bool
}

// Filesystem defines the interface for file storage operations.
type Filesystem interface {
	Put(path string, contents []byte) error
	Get(path string) ([]byte, error)
	Delete(path string) error
	Exists(path string) bool
	Size(path string) (int64, error)
	LastModified(path string) (int64, error)
}

// Validator defines the interface for validation operations.
type Validator interface {
	Validate(data map[string]any, rules map[string][]string) map[string][]error
}

// Database defines the interface for database operations.
type Database interface {
	Table(name string) QueryBuilder
	Raw() any
}

// QueryBuilder defines the interface for query building.
type QueryBuilder interface {
	Where(column, operator, value any) QueryBuilder
	Select(columns ...string) QueryBuilder
	Limit(limit int) QueryBuilder
	Get() (any, error)
	First() (any, error)
	Count(columns ...string) (int64, error)
	Insert(values map[string]any) (int64, error)
	Update(values map[string]any) (int64, error)
	Delete() (int64, error)
}

// Logger defines the interface for logging operations.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Config defines the interface for configuration operations.
type Config interface {
	Get(key string, defaultValue ...string) string
	GetInt(key string, defaultValue ...int) int
	GetBool(key string, defaultValue ...bool) bool
	GetFloat(key string, defaultValue ...float64) float64
	Set(key string, value string)
	Has(key string) bool
	All() map[string]string
}

// Router defines the interface for routing operations.
type Router interface {
	Get(path string, handler any)
	Post(path string, handler any)
	Put(path string, handler any)
	Patch(path string, handler any)
	Delete(path string, handler any)
	Group(prefix string, fn func())
}

// Response defines the interface for HTTP responses.
type Response interface {
	Status(code int) Response
	Header(key, value string) Response
	Json(data any) Response
	Send() error
}

// Request defines the interface for HTTP requests.
type Request interface {
	Get(key string, defaultValue ...string) string
	Has(key string) bool
	All() map[string]any
	Input(key string, defaultValue ...string) string
	Only(keys ...string) map[string]any
	Except(keys ...string) map[string]any
	HasFile(key string) bool
	File(key string) any
}