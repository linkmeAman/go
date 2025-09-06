package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourusername/saas-billing/internal/billing"
	"github.com/yourusername/saas-billing/internal/db"
	"github.com/yourusername/saas-billing/internal/middleware"
	"github.com/yourusername/saas-billing/internal/orgs"
	"github.com/yourusername/saas-billing/internal/types"
	"github.com/yourusername/saas-billing/internal/users"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateOrgRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=admin member"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database, err := db.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize services
	userService := users.NewUserService(database)
	orgService := orgs.NewOrganizationService(database)
	billingService := billing.NewBillingService(database)

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		if err := database.Ping(); err != nil {
			c.JSON(500, types.ApiResponse{Success: false, Error: "Database connection failed"})
			return
		}
		c.JSON(200, types.ApiResponse{Success: true, Data: gin.H{"status": "healthy"}})
	})

	// Public routes
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				var req RegisterRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, types.ApiResponse{Success: false, Error: err.Error()})
					return
				}

				if err := userService.Register(req.Email, req.Password); err != nil {
					c.JSON(http.StatusInternalServerError, types.ApiResponse{Success: false, Error: "Failed to register user"})
					return
				}

				c.JSON(http.StatusCreated, types.ApiResponse{Success: true, Data: gin.H{"message": "User registered successfully"}})
			})

			auth.POST("/login", func(c *gin.Context) {
				var req LoginRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, types.ApiResponse{Success: false, Error: err.Error()})
					return
				}

				token, err := userService.Login(req.Email, req.Password)
				if err != nil {
					c.JSON(http.StatusUnauthorized, types.ApiResponse{Success: false, Error: "Invalid credentials"})
					return
				}

				c.JSON(http.StatusOK, types.ApiResponse{Success: true, Data: gin.H{"token": token}})
			})
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired())
		{
			orgs := protected.Group("/organizations")
			{
				// Create organization
				orgs.POST("", func(c *gin.Context) {
					var req CreateOrgRequest
					if err := c.ShouldBindJSON(&req); err != nil {
						c.JSON(http.StatusBadRequest, types.ApiResponse{Success: false, Error: err.Error()})
						return
					}

					userID := c.GetString("userID")
					org, err := orgService.Create(req.Name, userID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, types.ApiResponse{Success: false, Error: "Failed to create organization"})
						return
					}

					c.JSON(http.StatusCreated, types.ApiResponse{Success: true, Data: org})
				})

				// List user's organizations
				orgs.GET("", func(c *gin.Context) {
					userID := c.GetString("userID")
					orgs, err := orgService.GetUserOrgs(userID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, types.ApiResponse{Success: false, Error: "Failed to fetch organizations"})
						return
					}

					c.JSON(http.StatusOK, types.ApiResponse{Success: true, Data: orgs})
				})

				// Organization-specific routes
				org := orgs.Group("/:orgID")
				{
					// Add member to organization (admin only)
					org.POST("/members", middleware.RequireRole(orgService, "owner", "admin"), func(c *gin.Context) {
						var req AddMemberRequest
						if err := c.ShouldBindJSON(&req); err != nil {
							c.JSON(http.StatusBadRequest, types.ApiResponse{Success: false, Error: err.Error()})
							return
						}

						orgID := c.Param("orgID")
						if err := orgService.AddMember(orgID, req.UserID, req.Role); err != nil {
							c.JSON(http.StatusInternalServerError, types.ApiResponse{Success: false, Error: "Failed to add member"})
							return
						}

						c.JSON(http.StatusOK, types.ApiResponse{Success: true, Data: gin.H{"message": "Member added successfully"}})
					})

					// Billing routes
					billing := org.Group("/billing")
					billing.Use(middleware.RequireRole(orgService, "owner", "admin"))
					{
						// Get available plans
						billing.GET("/plans", func(c *gin.Context) {
							plans, err := billingService.GetPlans()
							if err != nil {
								c.JSON(http.StatusInternalServerError, types.ApiResponse{Success: false, Error: "Failed to fetch plans"})
								return
							}

							c.JSON(http.StatusOK, types.ApiResponse{Success: true, Data: plans})
						})

						// Subscribe to plan
						billing.POST("/subscribe/:planID", func(c *gin.Context) {
							orgID := c.Param("orgID")
							planID := c.Param("planID")

							sub, err := billingService.CreateSubscription(orgID, planID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, types.ApiResponse{Success: false, Error: "Failed to create subscription"})
								return
							}

							c.JSON(http.StatusCreated, types.ApiResponse{Success: true, Data: sub})
						})

						// Get current subscription
						billing.GET("/subscription", func(c *gin.Context) {
							orgID := c.Param("orgID")
							sub, err := billingService.GetOrgSubscription(orgID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, types.ApiResponse{Success: false, Error: "Failed to fetch subscription"})
								return
							}

							if sub == nil {
								c.JSON(http.StatusNotFound, types.ApiResponse{Success: false, Error: "No active subscription found"})
								return
							}

							c.JSON(http.StatusOK, types.ApiResponse{Success: true, Data: sub})
						})

						// Get invoices
						billing.GET("/invoices", func(c *gin.Context) {
							orgID := c.Param("orgID")
							invoices, err := billingService.GetInvoices(orgID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, types.ApiResponse{Success: false, Error: "Failed to fetch invoices"})
								return
							}

							c.JSON(http.StatusOK, types.ApiResponse{Success: true, Data: invoices})
						})
					}
				}
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
