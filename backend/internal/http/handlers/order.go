package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"image"
	"image/jpeg"
	_ "image/png" // Register PNG decoder
	"os"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	httpmw "github.com/ryangel/ryangel-backend/internal/http/middleware"
	"github.com/ryangel/ryangel-backend/internal/repository"
	"github.com/ryangel/ryangel-backend/internal/models"
	authsvc "github.com/ryangel/ryangel-backend/internal/services/auth"
)

type OrderHandler struct {
    Orders *repository.OrderRepository
}

func (h OrderHandler) Register(rg *gin.RouterGroup, authSvc *authsvc.Service) {
	orders := rg.Group("/orders")
	if authSvc != nil {
		orders.Use(httpmw.ClientAuth(authSvc))
	}
	orders.GET("", h.listOrders)
	orders.POST("", h.createOrder)
}

func (h OrderHandler) RegisterAdmin(rg *gin.RouterGroup, authSvc *authsvc.Service) {
	admin := rg.Group("/admin/orders")
	if authSvc != nil {
		admin.Use(httpmw.AdminAuth(authSvc))
	}
	admin.GET("/stats", h.getDashboardStats)
	admin.GET("", h.adminListOrders)
	admin.GET("/:id/items", h.adminGetOrderItems)
	admin.PATCH("/:id/status", h.adminUpdateStatus)
}

func (h OrderHandler) adminListOrders(c *gin.Context) {
    page := 1
    limit := 50
    
    if p := c.Query("page"); p != "" {
        fmt.Sscanf(p, "%d", &page)
        if page < 1 {
            page = 1
        }
    }

    offset := (page - 1) * limit
    orders, err := h.Orders.GetOrders(c.Request.Context(), limit, offset)
    if err != nil {
        writeError(c, http.StatusInternalServerError, "DB_ERROR", "Failed to list orders", nil)
        return
    }
    
    if orders == nil {
        orders = []*models.Order{}
    }
    c.JSON(http.StatusOK, orders)
}

func (h OrderHandler) adminGetOrderItems(c *gin.Context) {
    idStr := c.Param("id")
    var id int64
    _, err := fmt.Sscanf(idStr, "%d", &id)
    if err != nil {
        writeError(c, http.StatusBadRequest, "INVALID_ID", "Invalid order ID", nil)
        return
    }
    
    items, err := h.Orders.GetOrderItems(c.Request.Context(), id)
    if err != nil {
         writeError(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch items", nil)
         return
    }
    c.JSON(http.StatusOK, items)
}

type updateStatusRequest struct {
    Status string `json:"status" binding:"required"`
}

func (h OrderHandler) adminUpdateStatus(c *gin.Context) {
    idStr := c.Param("id")
    var id int64
    _, err := fmt.Sscanf(idStr, "%d", &id)
    if err != nil {
        writeError(c, http.StatusBadRequest, "INVALID_ID", "Invalid order ID", nil)
        return
    }

    var req updateStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        writeValidationError(c, err)
        return
    }

    // Basic validation of status enum could be here

    if err := h.Orders.UpdateStatus(c.Request.Context(), id, models.OrderStatus(req.Status)); err != nil {
         writeError(c, http.StatusInternalServerError, "DB_ERROR", "Failed to update status", nil)
         return
    }
    c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h OrderHandler) getDashboardStats(c *gin.Context) {
	stats, err := h.Orders.GetDashboardStats(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch stats", nil)
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h OrderHandler) createOrder(c *gin.Context) {
	client, ok := httpmw.ClientFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not logged in", nil)
		return
	}

	// Parse multipart form
	ebuyStoreID := c.PostForm("ebuy_store_id")
	name := c.PostForm("name")
	email := c.PostForm("email")
	instagram := c.PostForm("instagram")
	phone := c.PostForm("phone") // verification only

	if ebuyStoreID == "" || name == "" {
		writeError(c, http.StatusBadRequest, "INVALID_DATA", "Missing required fields (store, name)", nil)
		return
	}

	// Handle file upload
	fileHeader, err := c.FormFile("payment_proof")
	var proofPath string
	if err == nil {
		if fileHeader.Size > 5*1024*1024 {
			writeError(c, http.StatusBadRequest, "UPLOAD_ERROR", "File too large (max 5MB)", nil)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			writeError(c, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to open file", nil)
			return
		}
		defer file.Close()

		// Decode image
		img, _, err := image.Decode(file)
		if err != nil {
			writeError(c, http.StatusBadRequest, "UPLOAD_ERROR", "Invalid image format", nil)
			return
		}

		// Resize if too large (e.g. max width 1024px)
		if img.Bounds().Dx() > 1024 {
			img = imaging.Resize(img, 1024, 0, imaging.Lanczos)
		}

		// Save to /var/www/media/uploads/proofs/
		filename := fmt.Sprintf("%d_%d_%s", client.ID, time.Now().Unix(), "proof.jpg") // Force jpg
		saveDir := "/var/www/media/uploads/proofs"
		
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			// Log error via fmt or logger if available
			fmt.Printf("Failed to create directory %s: %v\n", saveDir, err)
			writeError(c, http.StatusInternalServerError, "UPLOAD_ERROR", "Internal server error", nil)
			return
		}

		fullPath := filepath.Join(saveDir, filename)

		// Create file
		out, err := os.Create(fullPath)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to save file", nil)
			return
		}
		defer out.Close()

		// Compress and save as JPEG quality 80
		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 80}); err != nil {
			writeError(c, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to encode image", nil)
			return
		}

		// DB stores the path served by Nginx: /media/uploads/proofs/filename 
		// (assuming nginx maps /media -> /var/www/media)
		proofPath = "/media/uploads/proofs/" + filename
	} else if err != http.ErrMissingFile {
		writeError(c, http.StatusBadRequest, "UPLOAD_ERROR", "Error uploading file", nil)
		return
	}

	order, err := h.Orders.CreateOrder(c.Request.Context(), repository.CreateOrderParams{
		ClientID: client.ID,
		EbuyStoreID: ebuyStoreID,
		Name: name,
		Email: email,
		Instagram: instagram,
		Phone: phone,
		ProofPath: proofPath,
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, "CREATE_ERROR", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": order})
}


func (h OrderHandler) listOrders(c *gin.Context) {
	client, ok := httpmw.ClientFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not logged in", nil)
		return
	}

	orders, err := h.Orders.GetByClientID(c.Request.Context(), client.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch orders", nil)
		return
	}

    var result []models.OrderWithItems
    for _, o := range orders {
        items, err := h.Orders.GetOrderItems(c.Request.Context(), o.OrderID)
        if err != nil {
            // Log error but maybe continue or fail? Failure is safer.
            writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch order items", nil)
            return
        }
        result = append(result, models.OrderWithItems{
            Order: *o,
            Items: items,
        })
    }

	c.JSON(http.StatusOK, gin.H{"orders": result})
}
