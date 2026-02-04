package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Import driver MySQL (side-effect import)
	"github.com/zsais/go-gin-prometheus"

	"Test2/config"
	httphandler "Test2/internal/delivery/http"
	"Test2/internal/repository/mysql"
	"Test2/internal/usecase"
)

func main() {
	// 1. Load Configuration
	// Dựa trên config.go để lấy tham số môi trường
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Database Connection
	// Sử dụng DSN từ config.GetDSN()
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Kiểm tra kết nối thực tế (Ping)
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established successfully")

	// 3. Dependency Injection (Wiring Layers)

	// Timeout cho context của mỗi request (được define trong UseCase)
	timeoutContext := 2 * time.Second

	// Layer 1: Repository
	// Lưu ý: Cần thêm hàm NewMysqlPostRepository vào package mysql như đã đề cập ở trên
	postRepo := mysql.NewMysqlPostRepository(db)
	cateRepo := mysql.NewMysqlCateRepository(db)

	// Layer 2: UseCase
	// Tiêm Repository và Timeout vào UseCase
	postUseCase := usecase.NewPostUseCase(postRepo, timeoutContext)
	cateUseCase := usecase.NewCateUseCase(cateRepo, timeoutContext)

	// Layer 3: Delivery (HTTP Handler)
	r := gin.Default()

	// Cấu hình để tự động tạo route /metrics
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	// Đăng ký routes và handler
	httphandler.NewPostHandler(r, postUseCase)
	httphandler.NewCateHandler(r, cateUseCase)

	// 4. Run Server
	serverAddr := cfg.AppPort
	log.Printf("Server is running on port %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
