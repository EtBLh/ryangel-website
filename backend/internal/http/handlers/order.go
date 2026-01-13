package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"image"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	httpmw "github.com/ryangel/ryangel-backend/internal/http/middleware"
	"github.com/ryangel/ryangel-backend/internal/repository"
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

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
