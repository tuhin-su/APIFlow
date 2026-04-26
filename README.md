# FlowGo 🚀

**FlowGo** is a lightweight and flexible API workflow runner written in Go.
It allows you to define a sequence of HTTP requests in a JSON configuration file and execute them with support for loops, delays, dynamic variables, and proxy settings.

---

## ✨ Features

* 🔁 Loop execution (supports infinite loop with `-1`)
* ⏱️ Configurable delay between iterations
* 🔗 Chain multiple API endpoints
* 🔐 Dynamic variable extraction (e.g., tokens)
* 📦 JSON-based configuration
* 🌐 Proxy support (HTTP + authentication)
* ⚡ Fast and minimal (pure Go, no heavy dependencies)

---

## 📁 Project Structure

```
flowgo/
├── main.go
├── Makefile
├── README.md
└── config.json
```

---

## 📦 Installation

### 🔧 Build locally

```
make build
```

Binary will be created in:

```
./build/flowgo
```

---

### ⚡ Install (Go bin)

```
make install
```

Make sure `$GOBIN` or `$GOPATH/bin` is in your `PATH`.

---

### 🌍 Install globally

```
sudo make install-global
```

---

## 🚀 Usage

### Run with config file

```
flowgo --load config.json
```

or shorthand:

```
flowgo -l config.json
```

---

### Show version

```
flowgo -v
```

---

## ⚙️ Configuration Example

```json
{
  "app_url": "http://app.test",
  "loop": 10,
  "delay": 1000,
  "headers": {
    "Accept": "application/json",
    "Content-Type": "application/json",
    "Authorization": "Bearer @token"
  },
  "proxy": {
    "url": "http://127.0.0.1:8080",
    "username": "user",
    "password": "pass"
  },
  "endpoints": [
    {
      "name": "Login",
      "url": "/api/auth/login",
      "method": "POST",
      "body": {
        "email": "test@example.com",
        "password": "password"
      },
      "set": {
        "@token": "token"
      }
    },
    {
      "name": "Get Profile",
      "url": "/api/auth/me",
      "method": "GET"
    },
    {
      "name": "Logout",
      "url": "/api/auth/logout",
      "method": "POST"
    }
  ]
}
```

---

## 🔄 How It Works

1. Loads the JSON config file
2. Executes endpoints sequentially
3. Extracts values from responses (e.g., tokens)
4. Stores variables for reuse (e.g., `@token`)
5. Repeats based on `loop` and `delay`

---

## 🔐 Variable System

You can define dynamic variables using the `set` field:

```
"set": {
  "@token": "token"
}
```

Then reuse it in headers or body:

```
"Authorization": "Bearer @token"
```

---

## 🌐 Proxy Support

Supports HTTP proxies with optional authentication:

```json
"proxy": {
  "url": "http://127.0.0.1:8080",
  "username": "user",
  "password": "pass"
}
```

---

## 🛠️ Development

### Run locally

```
make run
```

### Format code

```
make fmt
```

### Clean build

```
make clean
```

---

## 📦 Cross Platform Build

```
make build-all
```

---

## ⚠️ Disclaimer

This tool is intended for:

* API testing
* Automation
* Development workflows

Do **NOT** use it for:

* Unauthorized access
* API abuse or spamming
* Bypassing rate limits

---

## 🛣️ Roadmap

* [ ] CLI flags (threads, retries)
* [ ] Concurrent workers
* [ ] Retry & backoff support
* [ ] Cookie/session handling
* [ ] Proxy rotation
* [ ] Logging system

---

## 📄 License

MIT License

---

## 🤝 Contributing

Pull requests are welcome.
Feel free to open issues for bugs, ideas, or improvements.

---
