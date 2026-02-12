<div dir="rtl">

### 🚀 سیستم Enterprise Logger for Golang

یک سیستم Logging حرفه‌ای، ساخت‌یافته و Production-Grade برای Go
مناسب برای **Backend Systems، Microservices، APIها، Workerها و CLIها**

این پروژه بر پایه‌ی **zap** طراحی شده و از ابتدا با رویکرد:
- performance بالا
- observability
- testability
-  عدم وابستگی به framework

پیاده‌سازی شده است.

---

## ✨ ویژگی‌های کلیدی

- Structured logging کاملاً سازگار با ELK / Datadog / Loki
- Context-aware logger (request-scoped)
- Middleware حرفه‌ای برای HTTP
- پشتیبانی از Gin (و قابل توسعه برای سایر فریم‌ورک‌ها)
- Panic recovery با stacktrace
- Log level policy هوشمند (info / warn / error)
- Log sampling برای high-traffic systems
- قابل استفاده بدون middleware (CLI / Worker)
- Test-friendly و injectable (بدون وابستگی اجباری به singleton)


---

## 🧠 فلسفه طراحی

این logger بر اساس اصول زیر طراحی شده است:

- **Context First**
  هر request یک logger مخصوص به خود دارد.

- **Framework Agnostic Core**
  منطق اصلی روی `net/http` پیاده‌سازی شده است.

- **No Hidden Globals**
  singleton فقط برای CLI / Worker استفاده می‌شود.

- **Structured by Default**
  هیچ log متنی (unstructured) تولید نمی‌شود.

---

## ⚙️ نصب

```bash
go get github.com/your-org/enterprise-logger
```

---

## 🔧 پیکربندی Logger
<div dir="ltr">

```
log, err := logger.New(logger.Config{
    Env:          "production",
    Service:      "user-api",
    Level:        zap.InfoLevel, #Production default
    JSON:         true, #Production Level
    Sampling:     true, #Production Level
    EnableCaller: false, #Production Level
})
if err != nil {
    panic(err)
}
```
<div dir="rtl">

| فیلد | توضیح |
|---------|---------|
| Env | production یا development |
| Service | نام سرویس یا میکروسرویس که log تولید میکند، این مقدار به تمام log ها اضافه میشود |
| Level | حداقل سطح لاگی که ثبت می‌شود (Debug, Info, Warn, Error) |
| JSON | فرمت خروجی log (در صورت تمایل به نمایش مقادیر به‌صورت ستونی و رنگی، این فیلد را روی false تنظیم کنید) |
| Sampling | جلوگیری از Log Storm |
| EnableCaller | اضافه کردن File:Line (مثال: "caller": "user/handler.go:42") |

---

## 🌐 استفاده در Gin (با Middleware)

### ثبت Middleware

<div dir="ltr">

```
import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/askari/gpm/logger"
	"github.com/askari/gpm/logger/adapters"
)
	r := gin.New()
	r.Use(adapters.Logger(log))
	r.Use(gin.Recovery())
```

<div dir="rtl">

### استفاده در Handler

<div dir="ltr">

```
r.GET("/users/:id", func(c *gin.Context) {
    l := logger.From(c.Request.Context())

    l.Info("fetching user",
        zap.String("user_id", c.Param("id")),
    )

    c.JSON(200, gin.H{"status": "ok"})
})
```
<div dir="rtl">

### نمونه خروجی Log

<div dir="ltr">

```
{
  "timestamp": "2026-02-02T14:21:00Z",
  "level": "info",
  "service": "user-api",
  "request_id": "f1c2...",
  "method": "GET",
  "path": "/users/42",
  "status": 200,
  "latency": "12ms",
  "msg": "http_request_completed"
}
```
<div dir="rtl">

## 🧵 استفاده بدون Middleware
### تنظیم Global Logger

<div dir="ltr">

```
log, _ := logger.New(logger.Config{
    Env:     "production",
    Service: "email-worker",
    Level:   zap.InfoLevel,
    JSON:    true,
})

logger.SetGlobal(log)
```

<div dir="rtl">

#### استفاده در هرجای برنامه

<div dir="ltr">

```
logger.L().Info("worker started")
```

<div dir="rtl">

## 🧪 استفاده در Service Layer

<div dir="ltr">

```
func Process(ctx context.Context) {
    log := logger.From(ctx)

    log.Info("processing started",
        zap.String("job_id", "123"),
    )
}
```

<div dir="rtl">
