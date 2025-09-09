package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/linkmeAman/saas-billing/internal/billing"
	"github.com/linkmeAman/saas-billing/internal/db"
	"github.com/linkmeAman/saas-billing/internal/middleware"
	"github.com/linkmeAman/saas-billing/internal/orgs"
	"github.com/linkmeAman/saas-billing/internal/types"
	"github.com/linkmeAman/saas-billing/internal/users"
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
			c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
				Code:       "DATABASE_ERROR",
				Message:    "Database connection failed",
				StatusCode: http.StatusInternalServerError,
			}))
			return
		}
		c.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{"status": "healthy"}, nil))
	})

	// Public routes
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				var req RegisterRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, types.NewErrorResponse(&types.ErrorInfo{
						Code:       "INVALID_REQUEST",
						Message:    err.Error(),
						StatusCode: http.StatusBadRequest,
					}))
					return
				}

				if err := userService.Register(req.Email, req.Password); err != nil {
					c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
						Code:       "REGISTRATION_ERROR",
						Message:    "Failed to register user",
						Details:    err.Error(),
						StatusCode: http.StatusInternalServerError,
					}))
					return
				}

				c.JSON(http.StatusCreated, types.NewSuccessResponse(gin.H{"message": "User registered successfully"}, nil))
			})

			auth.POST("/login", func(c *gin.Context) {
				var req LoginRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, types.NewErrorResponse(&types.ErrorInfo{
						Code:       "INVALID_REQUEST",
						Message:    err.Error(),
						StatusCode: http.StatusBadRequest,
					}))
					return
				}

				token, err := userService.Login(req.Email, req.Password)
				if err != nil {
					c.JSON(http.StatusUnauthorized, types.NewErrorResponse(&types.ErrorInfo{
						Code:       "INVALID_CREDENTIALS",
						Message:    "Invalid credentials",
						Details:    err.Error(),
						StatusCode: http.StatusUnauthorized,
					}))
					return
				}

				c.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{"token": token}, nil))
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
						c.JSON(http.StatusBadRequest, types.NewErrorResponse(&types.ErrorInfo{
							Code:       "INVALID_REQUEST",
							Message:    err.Error(),
							StatusCode: http.StatusBadRequest,
						}))
						return
					}

					userID := c.GetString("userID")
					org, err := orgService.Create(req.Name, userID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
							Code:       "ORGANIZATION_CREATE_ERROR",
							Message:    "Failed to create organization",
							Details:    err.Error(),
							StatusCode: http.StatusInternalServerError,
						}))
						return
					}

					c.JSON(http.StatusCreated, types.NewSuccessResponse(org, nil))
				})

				// List user's organizations
				orgs.GET("", func(c *gin.Context) {
					userID := c.GetString("userID")
					orgs, err := orgService.GetUserOrgs(userID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
							Code:       "ORGANIZATION_FETCH_ERROR",
							Message:    "Failed to fetch organizations",
							Details:    err.Error(),
							StatusCode: http.StatusInternalServerError,
						}))
						return
					}

					c.JSON(http.StatusOK, types.NewSuccessResponse(orgs, nil))
				})

				// Organization-specific routes
				org := orgs.Group("/:orgID")
				{
					// Add member to organization (admin only)
					org.POST("/members", middleware.RequireRole(orgService, "owner", "admin"), func(c *gin.Context) {
						var req AddMemberRequest
						if err := c.ShouldBindJSON(&req); err != nil {
							c.JSON(http.StatusBadRequest, types.NewErrorResponse(&types.ErrorInfo{
								Code:       "INVALID_REQUEST",
								Message:    err.Error(),
								StatusCode: http.StatusBadRequest,
							}))
							return
						}

						orgID := c.Param("orgID")
						if err := orgService.AddMember(orgID, req.UserID, req.Role); err != nil {
							c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
								Code:       "MEMBER_ADD_ERROR",
								Message:    "Failed to add member",
								Details:    err.Error(),
								StatusCode: http.StatusInternalServerError,
							}))
							return
						}
						
						c.JSON(http.StatusOK, types.NewSuccessResponse(gin.H{"message": "Member added successfully"}, nil))
					})

					// Billing routes
					billing := org.Group("/billing")
					billing.Use(middleware.RequireRole(orgService, "owner", "admin"))
					{
						// Get available plans
						billing.GET("/plans", func(c *gin.Context) {
							plans, err := billingService.GetPlans()
							if err != nil {
								c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
									Code:       "PLANS_FETCH_ERROR",
									Message:    "Failed to fetch plans",
									Details:    err.Error(),
									StatusCode: http.StatusInternalServerError,
								}))
								return
							}

							c.JSON(http.StatusOK, types.NewSuccessResponse(plans, nil))
						})

						// Subscribe to plan
						billing.POST("/subscribe/:planID", func(c *gin.Context) {
							orgID := c.Param("orgID")
							planID := c.Param("planID")

							sub, err := billingService.CreateSubscription(orgID, planID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
									Code:       "SUBSCRIPTION_CREATE_ERROR",
									Message:    "Failed to create subscription",
									Details:    err.Error(),
									StatusCode: http.StatusInternalServerError,
								}))
								return
							}

							c.JSON(http.StatusCreated, types.NewSuccessResponse(sub, nil))
						})

						// Get current subscription
						billing.GET("/subscription", func(c *gin.Context) {
							orgID := c.Param("orgID")
							sub, err := billingService.GetOrgSubscription(orgID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
									Code:       "SUBSCRIPTION_FETCH_ERROR",
									Message:    "Failed to fetch subscription",
									Details:    err.Error(),
									StatusCode: http.StatusInternalServerError,
								}))
								return
							}

							if sub == nil {
								c.JSON(http.StatusNotFound, types.NewErrorResponse(&types.ErrorInfo{
									Code:       "SUBSCRIPTION_NOT_FOUND",
									Message:    "No active subscription found",
									StatusCode: http.StatusNotFound,
								}))
								return
							}

							c.JSON(http.StatusOK, types.NewSuccessResponse(sub, nil))
						})

						// Get invoices
						billing.GET("/invoices", func(c *gin.Context) {
							orgID := c.Param("orgID")
							invoices, err := billingService.GetInvoices(orgID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, types.NewErrorResponse(&types.ErrorInfo{
									Code:       "INVOICES_FETCH_ERROR",
									Message:    "Failed to fetch invoices",
									Details:    err.Error(),
									StatusCode: http.StatusInternalServerError,
								}))
								return
							}

							c.JSON(http.StatusOK, types.NewSuccessResponse(invoices, nil))
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
