// ─── backend/.env.example ──────────────────────────────
# 必填配置
DATABASE_URL=postgres://postgres:postgres@localhost:5432/todoapp?sslmode=disable
JWT_SECRET=your-256-bit-secret-minimum-32-characters!!
JWT_EXPIRATION=24h

# 可选配置
SERVER_PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:3000
LOG_LEVEL=info
DB_MAX_CONNS=25
DB_MIN_CONNS=5


// ─── backend/go.mod ──────────────────────────────
module todo-api

go 1.25.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/go-playground/validator/v10 v10.30.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.10.0
	golang.org/x/crypto v0.53.0
)

require (
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.15.0 // indirect
	github.com/bytedance/sonic/loader v0.5.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.59.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	go.mongodb.org/mongo-driver/v2 v2.5.0 // indirect
	golang.org/x/arch v0.22.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)


// ─── backend/go.sum ──────────────────────────────
github.com/bytedance/gopkg v0.1.3 h1:TPBSwH8RsouGCBcMBktLt1AymVo2TVsBVCY4b6TnZ/M=
github.com/bytedance/gopkg v0.1.3/go.mod h1:576VvJ+eJgyCzdjS+c4+77QF3p7ubbtiKARP3TxducM=
github.com/bytedance/sonic v1.15.0 h1:/PXeWFaR5ElNcVE84U0dOHjiMHQOwNIx3K4ymzh/uSE=
github.com/bytedance/sonic v1.15.0/go.mod h1:tFkWrPz0/CUCLEF4ri4UkHekCIcdnkqXw9VduqpJh0k=
github.com/bytedance/sonic/loader v0.5.0 h1:gXH3KVnatgY7loH5/TkeVyXPfESoqSBSBEiDd5VjlgE=
github.com/bytedance/sonic/loader v0.5.0/go.mod h1:AR4NYCk5DdzZizZ5djGqQ92eEhCCcdf5x77udYiSJRo=
github.com/cloudwego/base64x v0.1.6 h1:t11wG9AECkCDk5fMSoxmufanudBtJ+/HemLstXDLI2M=
github.com/cloudwego/base64x v0.1.6/go.mod h1:OFcloc187FXDaYHvrNIjxSe8ncn0OOM8gEHfghB2IPU=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/gabriel-vasile/mimetype v1.4.12 h1:e9hWvmLYvtp846tLHam2o++qitpguFiYCKbn0w9jyqw=
github.com/gabriel-vasile/mimetype v1.4.12/go.mod h1:d+9Oxyo1wTzWdyVUPMmXFvp4F9tea18J8ufA774AB3s=
github.com/gin-contrib/sse v1.1.0 h1:n0w2GMuUpWDVp7qSpvze6fAu9iRxJY4Hmj6AmBOU05w=
github.com/gin-contrib/sse v1.1.0/go.mod h1:hxRZ5gVpWMT7Z0B0gSNYqqsSCNIJMjzvm6fqCz9vjwM=
github.com/gin-gonic/gin v1.12.0 h1:b3YAbrZtnf8N//yjKeU2+MQsh2mY5htkZidOM7O0wG8=
github.com/gin-gonic/gin v1.12.0/go.mod h1:VxccKfsSllpKshkBWgVgRniFFAzFb9csfngsqANjnLc=
github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
github.com/go-playground/validator/v10 v10.30.1 h1:f3zDSN/zOma+w6+1Wswgd9fLkdwy06ntQJp0BBvFG0w=
github.com/go-playground/validator/v10 v10.30.1/go.mod h1:oSuBIQzuJxL//3MelwSLD5hc2Tu889bF0Idm9Dg26cM=
github.com/goccy/go-json v0.10.5 h1:Fq85nIqj+gXn/S5ahsiTlK3TmC85qgirsdTP/+DeaC4=
github.com/goccy/go-json v0.10.5/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
github.com/goccy/go-yaml v1.19.2 h1:PmFC1S6h8ljIz6gMRBopkjP1TVT7xuwrButHID66PoM=
github.com/goccy/go-yaml v1.19.2/go.mod h1:XBurs7gK8ATbW4ZPGKgcbrY1Br56PdM69F7LkFRi1kA=
github.com/golang-jwt/jwt/v5 v5.3.1 h1:kYf81DTWFe7t+1VvL7eS+jKFVWaUnK9cB1qbwn63YCY=
github.com/golang-jwt/jwt/v5 v5.3.1/go.mod h1:fxCRLWMO43lRc8nhHWY6LGqRcf+1gQWArsqaEUEa5bE=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/jackc/pgpassfile v1.0.0 h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=
github.com/jackc/pgpassfile v1.0.0/go.mod h1:CEx0iS5ambNFdcRtxPj5JhEz+xB6uRky5eyVu/W2HEg=
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 h1:iCEnooe7UlwOQYpKFhBabPMi4aNAfoODPEFNiAnClxo=
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761/go.mod h1:5TJZWKEWniPve33vlWYSoGYefn3gLQRzjfDlhSJ9ZKM=
github.com/jackc/pgx/v5 v5.10.0 h1:VhSvgU2jSli8o3AqIEOTJr7rZwAEUVo4E4XhR94Zfr0=
github.com/jackc/pgx/v5 v5.10.0/go.mod h1:mal1tBGAFfLHvZzaYh77YS/eC6IX9OWbRV1QIIM0Jn4=
github.com/jackc/puddle/v2 v2.2.2 h1:PR8nw+E/1w0GLuRFSmiioY6UooMp6KJv0/61nB7icHo=
github.com/jackc/puddle/v2 v2.2.2/go.mod h1:vriiEXHvEE654aYKXXjOvZM39qJ0q+azkZFrfEOc3H4=
github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
github.com/klauspost/cpuid/v2 v2.3.0 h1:S4CRMLnYUhGeDFDqkGriYKdfoFlDnMtqTiI/sFzhA9Y=
github.com/klauspost/cpuid/v2 v2.3.0/go.mod h1:hqwkgyIinND0mEev00jJYCxPNVRVXFQeu1XKlok6oO0=
github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
github.com/pelletier/go-toml/v2 v2.2.4 h1:mye9XuhQ6gvn5h28+VilKrrPoQVanw5PMw/TB0t5Ec4=
github.com/pelletier/go-toml/v2 v2.2.4/go.mod h1:2gIqNv+qfxSVS7cM2xJQKtLSTLUE9V8t9Stt+h56mCY=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/quic-go/qpack v0.6.0 h1:g7W+BMYynC1LbYLSqRt8PBg5Tgwxn214ZZR34VIOjz8=
github.com/quic-go/qpack v0.6.0/go.mod h1:lUpLKChi8njB4ty2bFLX2x4gzDqXwUpaO1DP9qMDZII=
github.com/quic-go/quic-go v0.59.0 h1:OLJkp1Mlm/aS7dpKgTc6cnpynnD2Xg7C1pwL6vy/SAw=
github.com/quic-go/quic-go v0.59.0/go.mod h1:upnsH4Ju1YkqpLXC305eW3yDZ4NfnNbmQRCMWS58IKU=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
github.com/stretchr/testify v1.10.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
github.com/ugorji/go/codec v1.3.1 h1:waO7eEiFDwidsBN6agj1vJQ4AG7lh2yqXyOXqhgQuyY=
github.com/ugorji/go/codec v1.3.1/go.mod h1:pRBVtBSKl77K30Bv8R2P+cLSGaTtex6fsA2Wjqmfxj4=
go.mongodb.org/mongo-driver/v2 v2.5.0 h1:yXUhImUjjAInNcpTcAlPHiT7bIXhshCTL3jVBkF3xaE=
go.mongodb.org/mongo-driver/v2 v2.5.0/go.mod h1:yOI9kBsufol30iFsl1slpdq1I0eHPzybRWdyYUs8K/0=
go.uber.org/mock v0.6.0 h1:hyF9dfmbgIX5EfOdasqLsWD6xqpNZlXblLB/Dbnwv3Y=
go.uber.org/mock v0.6.0/go.mod h1:KiVJ4BqZJaMj4svdfmHM0AUx4NJYO8ZNpPnZn1Z+BBU=
golang.org/x/arch v0.22.0 h1:c/Zle32i5ttqRXjdLyyHZESLD/bB90DCU1g9l/0YBDI=
golang.org/x/arch v0.22.0/go.mod h1:dNHoOeKiyja7GTvF9NJS1l3Z2yntpQNzgrjh1cU103A=
golang.org/x/crypto v0.53.0 h1:QZ4Muo8THX6CizN2vPPd5fBGHyogrdK9fG4wLPFUsto=
golang.org/x/crypto v0.53.0/go.mod h1:DNLU434OwVakk9PzuwV8w62mAJpRJL3vsgcfp4Qnsio=
golang.org/x/net v0.55.0 h1:bcvxaJn3e1U6InsFWt1JUq1aSjnRxLzT2rtD2KfkDF8=
golang.org/x/net v0.55.0/go.mod h1:L5U2KuzuOe1lY7Z+aWVIKK6qEeJXnXV9yzGA+WCHJww=
golang.org/x/sync v0.21.0 h1:HLII4xRRTtCRkxYp4HNFF0Js/Og6q2i++KXbg0gHCwM=
golang.org/x/sync v0.21.0/go.mod h1:9xrNwdLfx4jkKbNva9FpL6vEN7evnE43NNNJQ2LF3+0=
golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.46.0 h1:noSf2Fq6F8DBgS+LysIkx7rIExoNHJsxOAtPp4rthXw=
golang.org/x/sys v0.46.0/go.mod h1:4GL1E5IUh+htKOUEOaiffhrAeqysfVGipDYzABqnCmw=
golang.org/x/text v0.38.0 h1:sXmwo9DwP3OK9EZ7PqAdaooSGozfl/3a6/xJcbzPRhE=
golang.org/x/text v0.38.0/go.mod h1:YXZt3QhHUKYT53r2lLKFIVi6Ao1jdzrTR/KQ09qyxF4=
google.golang.org/protobuf v1.36.10 h1:AYd7cD/uASjIL6Q9LiTjz8JLcrh/88q5UObnmY3aOOE=
google.golang.org/protobuf v1.36.10/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=


// ─── backend/cmd/server/main.go ──────────────────────────────
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"todo-api/internal/auth"
	"todo-api/internal/config"
	"todo-api/internal/database"
	"todo-api/internal/middleware"
	"todo-api/internal/todo"
)

func main() {
	cfg := config.Load()

	// Initialize database connection pool
	pool, err := database.NewPool(cfg.DatabaseURL, cfg.DBMaxConns, cfg.DBMinConns)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Run migrations (development mode)
	if err := database.RunMigrations(pool, "migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize repositories
	authRepo := auth.NewRepository(pool)
	todoRepo := todo.NewRepository(pool)

	// Initialize services
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpiration)
	todoService := todo.NewService(todoRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	todoHandler := todo.NewHandler(todoService)

	// Setup Gin router
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		dbStatus := "healthy"
		if err := pool.Ping(ctx); err != nil {
			dbStatus = "unhealthy"
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "error",
				"database":  dbStatus,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"database":  dbStatus,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (no JWT required)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Protected routes (JWT required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/auth/me", authHandler.Me)
		}

		// Todo routes (JWT required)
		todoGroup := v1.Group("/todos")
		todoGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			todoGroup.POST("", todoHandler.CreateTodo)
			todoGroup.GET("", todoHandler.ListTodos)
			todoGroup.GET("/:todo_id", todoHandler.GetTodo)
			todoGroup.PATCH("/:todo_id", todoHandler.UpdateTodo)
			todoGroup.DELETE("/:todo_id", todoHandler.DeleteTodo)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}


// ─── backend/cmd/migrate/main.go ──────────────────────────────
package main

import (
	"log"
	"os"

	"todo-api/internal/config"
	"todo-api/internal/database"
)

func main() {
	cfg := config.Load()

	pool, err := database.NewPool(cfg.DatabaseURL, cfg.DBMaxConns, cfg.DBMinConns)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	migrationsDir := "migrations"
	if len(os.Args) > 1 {
		migrationsDir = os.Args[1]
	}

	if err := database.RunMigrations(pool, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}


// ─── backend/migrations/20250301000001_create_users.up.sql ──────────────────────────────
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_users_username_lower') THEN
        CREATE UNIQUE INDEX idx_users_username_lower ON users (LOWER(username));
    END IF;
END $$;


// ─── backend/migrations/20250301000002_create_todos.up.sql ──────────────────────────────
CREATE TABLE IF NOT EXISTS todos (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    due_date    DATE,
    completed   BOOLEAN NOT NULL DEFAULT FALSE,
    version     INTEGER NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_todos_user_id') THEN
        CREATE INDEX idx_todos_user_id ON todos (user_id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_todos_user_created_desc') THEN
        CREATE INDEX idx_todos_user_created_desc ON todos (user_id, created_at DESC);
    END IF;
END $$;


// ─── backend/internal/middleware/auth.go ──────────────────────────────
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	apperrors "todo-api/internal/errors"
)

// AuthMiddleware creates a JWT authentication middleware.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(apperrors.NewUnauthorizedError("unauthorized"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Error(apperrors.NewUnauthorizedError("unauthorized"))
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.Error(apperrors.NewUnauthorizedError("unauthorized"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Error(apperrors.NewUnauthorizedError("unauthorized"))
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.Error(apperrors.NewUnauthorizedError("unauthorized"))
			c.Abort()
			return
		}

		c.Set("user_id", int64(userIDFloat))
		c.Next()
	}
}

// getRequestID retrieves the request_id from the Gin context.
func getRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		if s, ok := id.(string); ok {
			return s
		}
	}
	return ""
}


// ─── backend/internal/middleware/cors.go ──────────────────────────────
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates a CORS middleware with configurable allowed origins.
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")
	// Trim spaces
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed
		allowed := false
		for _, o := range origins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			c.Header("Access-Control-Max-Age", "300")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}


// ─── backend/internal/middleware/error.go ──────────────────────────────
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
)

// ErrorMiddleware is a Gin middleware that catches errors added via c.Error()
// and formats them into the standard error response format.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only process if there are errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error (most specific)
		err := c.Errors.Last()
		if err == nil {
			return
		}

		// Determine the request ID
		requestID := getRequestID(c)

		// If it's an AppError, use its fields
		if appErr, ok := err.Err.(*apperrors.AppError); ok {
			appErr.RequestID = requestID
			c.JSON(appErr.Code, appErr)
			return
		}

		// Generic fallback
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":       500,
			"error_code": apperrors.ErrorCodeInternal,
			"message":    "internal server error",
			"request_id": requestID,
		})
	}
}


// ─── backend/internal/middleware/logger.go ──────────────────────────────
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggerMiddleware creates a request logging middleware.
// It generates a unique request_id for each request and logs method, path, status, and duration.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate and set request_id
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Log after request is processed
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[%s] %s %s %d %v", requestID, method, path, status, duration)
	}
}


// ─── backend/internal/middleware/security.go ──────────────────────────────
package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds security-related HTTP headers to every response.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}


// ─── backend/internal/database/migrate.go ──────────────────────────────
package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes SQL migration files from the migrations directory.
func RunMigrations(pool *pgxpool.Pool, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory %s: %w", migrationsDir, err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, fname := range upFiles {
		path := filepath.Join(migrationsDir, fname)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fname, err)
		}

		sql := string(content)
		if _, err := pool.Exec(context.Background(), sql); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fname, err)
		}
		log.Printf("migration applied: %s", fname)
	}

	return nil
}


// ─── backend/internal/database/postgres.go ──────────────────────────────
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a new database connection pool.
func NewPool(databaseURL string, maxConns, minConns int) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	cfg.MaxConns = int32(maxConns)
	cfg.MinConns = int32(minConns)
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}


// ─── backend/internal/config/config.go ──────────────────────────────
package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	JWTExpiration      time.Duration
	ServerPort         string
	CORSAllowedOrigins string
	LogLevel           string
	DBMaxConns         int
	DBMinConns         int
}

func Load() *Config {
	jwtSecret := getEnv("JWT_SECRET", "change-me-in-production-at-least-32-chars!!")
	if jwtSecret == "change-me-in-production-at-least-32-chars!!" {
		log.Println("WARNING: Using default JWT secret. This is insecure for production. Set JWT_SECRET environment variable.")
	}

	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todoapp?sslmode=disable"),
		JWTSecret:          jwtSecret,
		JWTExpiration:      getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		DBMaxConns:         getIntEnv("DB_MAX_CONNS", 25),
		DBMinConns:         getIntEnv("DB_MIN_CONNS", 5),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i
		}
	}
	return fallback
}


// ─── backend/internal/auth/dto.go ──────────────────────────────
package auth

// RegisterRequest represents the registration request body.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse represents the registration success response.
type RegisterResponse struct {
	Token string     `json:"token"`
	User  UserPublic `json:"user"`
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login success response.
type LoginResponse struct {
	Token string     `json:"token"`
	User  UserPublic `json:"user"`
}

// UserPublic represents the public user information.
type UserPublic struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}


// ─── backend/internal/auth/handler.go ──────────────────────────────
package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		return
	}

	resp, appErr := h.service.Register(c.Request.Context(), &req)
	if appErr != nil {
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		return
	}

	resp, appErr := h.service.Login(c.Request.Context(), &req)
	if appErr != nil {
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Me handles GET /api/v1/auth/me
func (h *Handler) Me(c *gin.Context) {
	userID := c.GetInt64("user_id")

	user, appErr := h.service.GetMe(c.Request.Context(), userID)
	if appErr != nil {
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}


// ─── backend/internal/auth/repository.go ──────────────────────────────
package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// User represents a user record from the database.
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

// CreateUser inserts a new user and returns the created user.
func (r *Repository) CreateUser(ctx context.Context, username, passwordHash string) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)
		 RETURNING id, username, password_hash`,
		username, passwordHash,
	).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username (case-insensitive via LOWER).
func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash FROM users WHERE LOWER(username) = LOWER($1)`,
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by their ID.
func (r *Repository) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}


// ─── backend/internal/auth/service.go ──────────────────────────────
package auth

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	apperrors "todo-api/internal/errors"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	// At least one letter and one digit
	hasLetter = regexp.MustCompile(`[a-zA-Z]`)
	hasDigit  = regexp.MustCompile(`\d`)
)

type Service struct {
	repo          *Repository
	jwtSecret     []byte
	jwtExpiration time.Duration
}

func NewService(repo *Repository, jwtSecret string, jwtExpiration time.Duration) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
	}
}

// Register creates a new user account and returns a JWT token.
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, *apperrors.AppError) {
	// Validate username
	if !usernameRegex.MatchString(req.Username) {
		return nil, apperrors.NewValidationError(
			"username must be 3-20 characters, allowing letters, digits and underscores",
		)
	}

	// Validate password
	if err := validatePassword(req.Password); err != nil {
		return nil, apperrors.NewValidationError(err.Error())
	}

	// Normalize username to lowercase
	username := strings.ToLower(req.Username)

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}

	// Create user
	user, err := s.repo.CreateUser(ctx, username, string(hash))
	if err != nil {
		// Check for unique constraint violation (PostgreSQL error code 23505)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperrors.NewConflictError(apperrors.ErrorCodeUsernameTaken, "username already exists")
		}
		return nil, apperrors.NewInternalError()
	}

	// Generate JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}

	return &RegisterResponse{
		Token: token,
		User: UserPublic{
			ID:       user.ID,
			Username: user.Username,
		},
	}, nil
}

// Login authenticates a user and returns a JWT token.
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, *apperrors.AppError) {
	username := strings.ToLower(req.Username)

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}
	if user == nil {
		return nil, apperrors.NewUnauthorizedError("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid username or password")
	}

	// Generate JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}

	return &LoginResponse{
		Token: token,
		User: UserPublic{
			ID:       user.ID,
			Username: user.Username,
		},
	}, nil
}

// GetMe returns the current user's public info based on user ID.
func (s *Service) GetMe(ctx context.Context, userID int64) (*UserPublic, *apperrors.AppError) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}
	if user == nil {
		return nil, apperrors.NewNotFoundError("user not found")
	}
	return &UserPublic{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

// generateToken creates a JWT token for the given user ID.
func (s *Service) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.jwtExpiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// validatePassword checks password strength rules.
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	letterOk := hasLetter.MatchString(password)
	digitOk := hasDigit.MatchString(password)
	if !letterOk && !digitOk {
		return errors.New("password must contain at least one letter and one digit")
	}
	if !letterOk {
		return errors.New("password must contain at least one letter")
	}
	if !digitOk {
		return errors.New("password must contain at least one digit")
	}
	return nil
}


// ─── backend/internal/todo/dto.go ──────────────────────────────
package todo

import "time"

// CreateTodoRequest represents the request body for creating a todo.
type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"` // format: YYYY-MM-DD
}

// UpdateTodoRequest represents the request body for updating a todo.
type UpdateTodoRequest struct {
	Version     int64   `json:"version" binding:"required"`
	Title       *string `json:"title"`
	Description *string `json:"description"` // empty string ("") clears the field (sets to NULL)
	DueDate     *string `json:"due_date"`    // empty string ("") clears the field (sets to NULL), format: YYYY-MM-DD
	Completed   *bool   `json:"completed"`
}

// DeleteTodoRequest represents the request body for deleting a todo.
type DeleteTodoRequest struct {
	Version int64 `json:"version" binding:"required"`
}

// TodoResponse represents the full todo object returned by the API.
type TodoResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *string    `json:"due_date"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Version     int64      `json:"version"`
	UserID      int64      `json:"user_id"`
}

// TodoListResponse represents the paginated list response.
type TodoListResponse struct {
	Items      []TodoResponse `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}


// ─── backend/internal/todo/handler.go ──────────────────────────────
package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateTodo handles POST /api/v1/todos
func (h *Handler) CreateTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		return
	}

	resp, appErr := h.service.Create(c.Request.Context(), userID, &req)
	if appErr != nil {
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListTodos handles GET /api/v1/todos
func (h *Handler) ListTodos(c *gin.Context) {
	userID := c.GetInt64("user_id")

	page, err := parseQueryInt(c, "page", 1)
	if err != nil {
		c.Error(apperrors.NewValidationError("page must be a valid integer"))
		return
	}
	if page < 1 {
		c.Error(apperrors.NewValidationError("page must be >= 1"))
		return
	}

	pageSize, err := parseQueryInt(c, "page_size", 20)
	if err != nil {
		c.Error(apperrors.NewValidationError("page_size must be a valid integer"))
		return
	}
	if pageSize < 1 {
		c.Error(apperrors.NewValidationError("page_size must be >= 1"))
		return
	}
	if pageSize > 100 {
		c.Error(apperrors.NewValidationError("page_size must not exceed 100"))
		return
	}

	resp, appErr := h.service.List(c.Request.Context(), userID, page, pageSize)
	if appErr != nil {
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTodo handles GET /api/v1/todos/:todo_id
func (h *Handler) GetTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.Error(apperrors.NewValidationError("invalid todo_id format"))
		return
	}

	resp, appErr := h.service.GetByID(c.Request.Context(), userID, todoID)
	if appErr != nil {
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateTodo handles PATCH /api/v1/todos/:todo_id
func (h *Handler) UpdateTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.Error(apperrors.NewValidationError("invalid todo_id format"))
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		return
	}

	resp, appErr := h.service.Update(c.Request.Context(), userID, todoID, &req)
	if appErr != nil {
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTodo handles DELETE /api/v1/todos/:todo_id
func (h *Handler) DeleteTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.Error(apperrors.NewValidationError("invalid todo_id format"))
		return
	}

	var req DeleteTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		return
	}

	if req.Version < 1 {
		c.Error(apperrors.NewValidationError("version must be >= 1"))
		return
	}

	appErr := h.service.Delete(c.Request.Context(), userID, todoID, req.Version)
	if appErr != nil {
		c.Error(appErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// parseQueryInt parses an integer query parameter with a default value.
func parseQueryInt(c *gin.Context, key string, defaultVal int) (int, error) {
	val := c.DefaultQuery(key, "")
	if val == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(val)
}


// ─── backend/internal/todo/repository.go ──────────────────────────────
package todo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Todo represents a todo record from the database.
type Todo struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Completed   bool       `json:"completed"`
	Version     int64      `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UpdateFields represents the fields that can be updated on a todo.
// For each pointer field, nil means the field is not updated.
// UpdateDescription/UpdateDueDate flags indicate whether to update the field.
// If the flag is true and the Val pointer is nil, the column is set to NULL.
type UpdateFields struct {
	Title             *string
	DescriptionVal    *string    // value to set (nil means set to NULL)
	UpdateDescription bool       // whether to update the description field
	DueDateVal        *time.Time // value to set (nil means set to NULL)
	UpdateDueDate     bool       // whether to update the due_date field
	Completed         *bool
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new todo record.
func (r *Repository) Create(ctx context.Context, todo *Todo) (*Todo, error) {
	var t Todo
	err := r.pool.QueryRow(ctx,
		`INSERT INTO todos (user_id, title, description, due_date)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, title, description, due_date, completed, version, created_at, updated_at`,
		todo.UserID, todo.Title, todo.Description, todo.DueDate,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.Version, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FindByIDAndUser finds a todo by id and user_id (ensures data isolation).
func (r *Repository) FindByIDAndUser(ctx context.Context, id, userID int64) (*Todo, error) {
	var t Todo
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, description, due_date, completed, version, created_at, updated_at
		 FROM todos WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.Version, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// ListByUser returns paginated todos for a user, ordered by created_at DESC.
func (r *Repository) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]*Todo, int, error) {
	// Count total
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM todos WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch page
	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, description, due_date, completed, version, created_at, updated_at
		 FROM todos WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.Version, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		todos = append(todos, &t)
	}

	return todos, total, nil
}

// UpdateVersioned updates a todo with dynamic fields and optimistic locking.
// Returns the updated todo, nil if not found, or errVersionConflict if version mismatch.
func (r *Repository) UpdateVersioned(ctx context.Context, todoID, userID, version int64, fields *UpdateFields) (*Todo, error) {
	// Build dynamic SET clause
	var setClauses []string
	args := []interface{}{}
	argIdx := 1

	if fields.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *fields.Title)
		argIdx++
	}
	if fields.UpdateDescription {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, fields.DescriptionVal) // nil maps to SQL NULL
		argIdx++
	}
	if fields.UpdateDueDate {
		setClauses = append(setClauses, fmt.Sprintf("due_date = $%d", argIdx))
		args = append(args, fields.DueDateVal) // nil maps to SQL NULL
		argIdx++
	}
	if fields.Completed != nil {
		setClauses = append(setClauses, fmt.Sprintf("completed = $%d", argIdx))
		args = append(args, *fields.Completed)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	// Add WHERE conditions
	setClauses = append(setClauses, "version = version + 1", "updated_at = NOW()")

	sql := fmt.Sprintf(
		`UPDATE todos SET %s WHERE id = $%d AND user_id = $%d AND version = $%d
		 RETURNING id, user_id, title, description, due_date, completed, version, created_at, updated_at`,
		strings.Join(setClauses, ", "),
		argIdx, argIdx+1, argIdx+2,
	)
	args = append(args, todoID, userID, version)

	var t Todo
	err := r.pool.QueryRow(ctx, sql, args...).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.Version, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Check if record exists at all
			existing, checkErr := r.FindByIDAndUser(ctx, todoID, userID)
			if checkErr != nil {
				return nil, checkErr
			}
			if existing == nil {
				return nil, nil // not found
			}
			// Record exists but version mismatch
			return nil, errVersionConflict
		}
		return nil, err
	}
	return &t, nil
}

// DeleteVersioned deletes a todo with optimistic locking.
// Returns true if deleted, false if not found, errVersionConflict if version mismatch.
func (r *Repository) DeleteVersioned(ctx context.Context, id, userID, version int64) (bool, error) {
	ct, err := r.pool.Exec(ctx,
		`DELETE FROM todos WHERE id = $1 AND user_id = $2 AND version = $3`,
		id, userID, version,
	)
	if err != nil {
		return false, err
	}

	if ct.RowsAffected() == 1 {
		return true, nil
	}

	// Check if record exists
	existing, checkErr := r.FindByIDAndUser(ctx, id, userID)
	if checkErr != nil {
		return false, checkErr
	}
	if existing == nil {
		return false, nil // not found
	}

	// Record exists but version mismatch
	return false, errVersionConflict
}

// errVersionConflict is a sentinel error for version conflicts.
var errVersionConflict = errors.New("version conflict")


// ─── backend/internal/todo/service.go ──────────────────────────────
package todo

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	apperrors "todo-api/internal/errors"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var errInvalidDateFormat = errors.New("due_date must be in YYYY-MM-DD format")

const (
	maxTitleLen       = 255
	maxDescriptionLen = 1000
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new todo for the given user.
func (s *Service) Create(ctx context.Context, userID int64, req *CreateTodoRequest) (*TodoResponse, *apperrors.AppError) {
	// Validate title
	if req.Title == "" {
		return nil, apperrors.NewValidationError("title is required")
	}
	if utf8.RuneCountInString(req.Title) > maxTitleLen {
		return nil, apperrors.NewValidationError("title must not exceed 255 characters")
	}

	// Validate description length
	if req.Description != nil && utf8.RuneCountInString(*req.Description) > maxDescriptionLen {
		return nil, apperrors.NewValidationError("description must not exceed 1000 characters")
	}

	// Parse and validate due_date if provided
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := parseDate(*req.DueDate)
		if err != nil {
			return nil, apperrors.NewValidationError(err.Error())
		}
		dueDate = &parsed
	}

	todo := &Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
		Completed:   false,
		Version:     1,
	}

	created, err := s.repo.Create(ctx, todo)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}

	return toResponse(created), nil
}

// GetByID retrieves a single todo by ID, ensuring it belongs to the user.
func (s *Service) GetByID(ctx context.Context, userID, todoID int64) (*TodoResponse, *apperrors.AppError) {
	todo, err := s.repo.FindByIDAndUser(ctx, todoID, userID)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}
	if todo == nil {
		return nil, apperrors.NewNotFoundError("todo not found")
	}
	return toResponse(todo), nil
}

// List returns a paginated list of todos for the user.
func (s *Service) List(ctx context.Context, userID int64, page, pageSize int) (*TodoListResponse, *apperrors.AppError) {
	todos, total, err := s.repo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	items := make([]TodoResponse, 0, len(todos))
	for _, t := range todos {
		items = append(items, *toResponse(t))
	}

	return &TodoListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates a todo with optimistic locking.
func (s *Service) Update(ctx context.Context, userID, todoID int64, req *UpdateTodoRequest) (*TodoResponse, *apperrors.AppError) {
	// Validate that at least one field (other than version) is provided
	if req.Title == nil && req.Description == nil && req.DueDate == nil && req.Completed == nil {
		return nil, apperrors.NewValidationError("at least one field to update (other than version) must be provided")
	}

	// Validate title if provided
	if req.Title != nil {
		if *req.Title == "" {
			return nil, apperrors.NewValidationError("title cannot be empty")
		}
		if utf8.RuneCountInString(*req.Title) > maxTitleLen {
			return nil, apperrors.NewValidationError("title must not exceed 255 characters")
		}
	}

	// Validate description if provided
	if req.Description != nil && utf8.RuneCountInString(*req.Description) > maxDescriptionLen {
		return nil, apperrors.NewValidationError("description must not exceed 1000 characters")
	}

	// Parse and validate due_date if provided
	var updateDueDate bool
	var dueDateVal *time.Time
	if req.DueDate != nil {
		updateDueDate = true
		if *req.DueDate != "" {
			parsed, err := parseDate(*req.DueDate)
			if err != nil {
				return nil, apperrors.NewValidationError(err.Error())
			}
			dueDateVal = &parsed
		}
		// If empty string, dueDateVal stays nil which maps to SQL NULL
	}

	// Check idempotent case: if only version + completed is provided
	// and the completed value matches the current state, return without modification.
	if req.Title == nil && req.Description == nil && req.DueDate == nil && req.Completed != nil {
		existing, err := s.repo.FindByIDAndUser(ctx, todoID, userID)
		if err != nil {
			return nil, apperrors.NewInternalError()
		}
		if existing == nil {
			return nil, apperrors.NewNotFoundError("todo not found")
		}
		if existing.Version != req.Version {
			details := formatVersionDetail(existing.Version)
			return nil, apperrors.NewVersionConflictError(details)
		}
		if existing.Completed == *req.Completed {
			// Idempotent: return current state without modification
			return toResponse(existing), nil
		}
	}

	// Build update fields
	fields := &UpdateFields{
		Title:         req.Title,
		DueDateVal:    dueDateVal,
		UpdateDueDate: updateDueDate,
		Completed:     req.Completed,
	}
	if req.Description != nil {
		fields.UpdateDescription = true
		if *req.Description == "" {
			fields.DescriptionVal = nil // clear to NULL
		} else {
			fields.DescriptionVal = req.Description
		}
	}

	updated, err := s.repo.UpdateVersioned(ctx, todoID, userID, req.Version, fields)
	if err != nil {
		if err == errVersionConflict {
			// Get current version for details message
			existing, _ := s.repo.FindByIDAndUser(ctx, todoID, userID)
			if existing != nil {
				details := formatVersionDetail(existing.Version)
				return nil, apperrors.NewVersionConflictError(details)
			}
			return nil, apperrors.NewVersionConflictError("")
		}
		return nil, apperrors.NewInternalError()
	}
	if updated == nil {
		return nil, apperrors.NewNotFoundError("todo not found")
	}

	return toResponse(updated), nil
}

// Delete deletes a todo with optimistic locking.
func (s *Service) Delete(ctx context.Context, userID, todoID, version int64) *apperrors.AppError {
	deleted, err := s.repo.DeleteVersioned(ctx, todoID, userID, version)
	if err != nil {
		if err == errVersionConflict {
			existing, _ := s.repo.FindByIDAndUser(ctx, todoID, userID)
			if existing != nil {
				details := formatVersionDetail(existing.Version)
				return apperrors.NewVersionConflictError(details)
			}
			return apperrors.NewVersionConflictError("")
		}
		return apperrors.NewInternalError()
	}
	if !deleted {
		return apperrors.NewNotFoundError("todo not found")
	}
	return nil
}

// parseDate validates and parses a YYYY-MM-DD date string.
func parseDate(s string) (time.Time, error) {
	if !dateRegex.MatchString(s) {
		return time.Time{}, errInvalidDateFormat
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errInvalidDateFormat
	}
	return t, nil
}

// formatVersionDetail creates a details string for version conflict errors.
func formatVersionDetail(currentVersion int64) string {
	return "current_version = " + strconv.FormatInt(currentVersion, 10)
}

// toResponse converts a Todo model to a TodoResponse.
func toResponse(t *Todo) *TodoResponse {
	var dueDate *string
	if t.DueDate != nil {
		s := t.DueDate.Format("2006-01-02")
		dueDate = &s
	}

	return &TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     dueDate,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Version:     t.Version,
		UserID:      t.UserID,
	}
}


// ─── backend/internal/errors/errors.go ──────────────────────────────
package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Standard error codes
const (
	ErrorCodeValidation    = "VALIDATION_ERROR"
	ErrorCodeUnauthorized  = "UNAUTHORIZED"
	ErrorCodeForbidden     = "FORBIDDEN"
	ErrorCodeNotFound      = "NOT_FOUND"
	ErrorCodeUsernameTaken = "USERNAME_TAKEN"
	ErrorCodeVersionConflict = "VERSION_CONFLICT"
	ErrorCodeInternal      = "INTERNAL_ERROR"
)

// AppError represents a structured application error.
type AppError struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Details   string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, errorCode, message string) *AppError {
	return &AppError{
		Code:      code,
		ErrorCode: errorCode,
		Message:   message,
	}
}

func NewValidationError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, ErrorCodeValidation, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, ErrorCodeUnauthorized, message)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, ErrorCodeNotFound, message)
}

func NewConflictError(errorCode, message string) *AppError {
	return NewAppError(http.StatusConflict, errorCode, message)
}

func NewVersionConflictError(details string) *AppError {
	err := NewAppError(http.StatusConflict, ErrorCodeVersionConflict, "resource conflict due to version mismatch")
	if details != "" {
		err.Details = details
	}
	return err
}

func NewInternalError() *AppError {
	return NewAppError(http.StatusInternalServerError, ErrorCodeInternal, "internal server error")
}

// NewValidationErrorFromBinding extracts validation errors from ShouldBindJSON errors
// and returns a descriptive AppError with field-level details.
func NewValidationErrorFromBinding(bindErr error) *AppError {
	var ve validator.ValidationErrors
	if ok := AsValidationErrors(bindErr, &ve); ok && len(ve) > 0 {
		var errMsgs []string
		for _, fe := range ve {
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' %s", fe.Field(), fe.Tag()))
		}
		msg := "validation failed: " + strings.Join(errMsgs, "; ")
		return NewAppError(http.StatusBadRequest, ErrorCodeValidation, msg)
	}
	return NewValidationError("invalid request body")
}

// AsValidationErrors checks if the error is a validator.ValidationErrors.
// Exposed as a function to avoid import cycle issues in handlers.
var AsValidationErrors = func(err error, target *validator.ValidationErrors) bool {
	if err == nil {
		return false
	}
	ve, ok := err.(validator.ValidationErrors)
	if ok {
		*target = ve
		return true
	}
	// Also check wrapped errors
	if errors.As(err, &ve) {
		*target = ve
		return true
	}
	return false
}


