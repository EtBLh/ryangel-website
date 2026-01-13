package handlers

import (
	"fmt"
	"net/http"
    "sort"
    "strconv"

	"github.com/gin-gonic/gin"

	httpmw "github.com/ryangel/ryangel-backend/internal/http/middleware"
	"github.com/ryangel/ryangel-backend/internal/models"
	"github.com/ryangel/ryangel-backend/internal/repository"
	authsvc "github.com/ryangel/ryangel-backend/internal/services/auth"
)

// CartHandler handles cart-related HTTP requests.
type CartHandler struct {
	Repo           *repository.CartRepository
	EbuyStoreRepo  *repository.EbuyStoreRepository
    DiscountRepo   *repository.DiscountRepository
}

// Register wires the cart routes onto the router.
func (h CartHandler) Register(rg *gin.RouterGroup, authService interface{}) {
	if authSvc, ok := authService.(*authsvc.Service); ok {
		// Endpoints that Support Optional Auth
		rg.GET("/cart", httpmw.OptionalClientAuth(authSvc), h.GetCart)
		rg.POST("/cart/items", httpmw.OptionalClientAuth(authSvc), h.AddItemToCart)
		rg.PATCH("/cart/items/:cart_item_id", httpmw.OptionalClientAuth(authSvc), h.UpdateCartItem)
		rg.DELETE("/cart/items/:cart_item_id", httpmw.OptionalClientAuth(authSvc), h.RemoveCartItem)
		
		// Endpoints that Require Auth
		rg.POST("/cart/checkout", httpmw.ClientAuth(authSvc), h.Checkout)
		rg.POST("/cart/apply-discount", httpmw.ClientAuth(authSvc), h.ApplyDiscount)
		rg.DELETE("/cart/discount", httpmw.ClientAuth(authSvc), h.RemoveDiscount)
	} else {
		// Fallback for no auth service (e.g. testing without auth)
		rg.GET("/cart", h.GetCart)
		rg.POST("/cart/items", h.AddItemToCart)
		rg.PATCH("/cart/items/:cart_item_id", h.UpdateCartItem)
		rg.DELETE("/cart/items/:cart_item_id", h.RemoveCartItem)
	}
}

// getOrCreateCartID gets or creates cart ID.
func (h CartHandler) getOrCreateCartID(c *gin.Context) (string, error) {
	// If authenticated, get or create cart by client_id
	if client, exists := httpmw.ClientFromContext(c); exists {
		cart, err := h.Repo.GetCartByClientID(c.Request.Context(), client.ID)
		if err != nil {
			return "", err
		}

		// Check for anonymous cart to claim if user has no cart
		anonCartID := c.GetHeader("X-Cart-ID")
		if cart == nil && anonCartID != "" && isValidUUID(anonCartID) {
			anonCart, err := h.Repo.GetCartByID(c.Request.Context(), anonCartID)
			if err == nil && anonCart != nil && anonCart.ClientID == nil {
				// Claim the anonymous cart
				if err := h.Repo.UpdateCartOwner(c.Request.Context(), anonCartID, client.ID); err == nil {
					return anonCartID, nil
				}
			}
		}

		if cart != nil {
			return cart.CartID, nil
		}
		// Create new cart for client
		newCart, err := h.Repo.CreateCart(c.Request.Context(), &client.ID)
		if err != nil {
			return "", err
		}
		return newCart.CartID, nil
	}

	// Anonymous: use X-Cart-ID if provided, else create new
	cartIDStr := c.GetHeader("X-Cart-ID")
	if cartIDStr != "" && isValidUUID(cartIDStr) {
		// Check if exists, if not, create
		cart, err := h.Repo.GetCartByID(c.Request.Context(), cartIDStr)
		if err != nil {
			// If cart not found, create new one
			if err.Error() == "cart not found" {
				newCart, err := h.Repo.CreateCart(c.Request.Context(), nil)
				if err != nil {
					return "", err
				}
				return newCart.CartID, nil
			}
			return "", err
		}
		if cart != nil {
			// SECURITY: If cart belongs to a registered client, anonymous user cannot use it.
			if cart.ClientID != nil {
				return "", fmt.Errorf("unauthorized access to client cart")
			}
			return cart.CartID, nil
		}
	}

	// Create new anonymous cart
	newCart, err := h.Repo.CreateCart(c.Request.Context(), nil)
	if err != nil {
		return "", err
	}
	return newCart.CartID, nil
}

// GetCart handles GET /cart.
func (h CartHandler) GetCart(c *gin.Context) {
	cartID, err := h.getOrCreateCartID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := h.Repo.GetCartItems(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart items"})
		return
	}

	// Calculate totals
	var subtotal float64
	for _, item := range items {
		subtotal += item.UnitPrice * float64(item.Quantity)
	}

    // Default shipping fee
    originalShippingFee := 5.0
    finalShippingFee := 5.0
	itemDiscountAmount := 0.0

    // Apply active auto-apply discounts
    if h.DiscountRepo != nil {
        discounts, err := h.DiscountRepo.GetAutoApplyDiscounts(c.Request.Context())
        if err == nil {
            for _, d := range discounts {
                // Logic for BXGY
                if d.DiscountType == "bxgy" {
                    // Check restrictions
                    if d.ProductTypeRestriction != nil {
                         restriction := *d.ProductTypeRestriction
                         // Find applicable items
                         var applicablePrices []float64
                         for _, item := range items {
                             if item.ProductType == restriction {
                                 for k := 0; k < item.Quantity; k++ {
                                     applicablePrices = append(applicablePrices, item.UnitPrice)
                                 }
                             }
                         }
                         
                         if len(applicablePrices) > 0 && d.BuyQuantity != nil && d.GetQuantity != nil {
                             buy := *d.BuyQuantity
                             get := *d.GetQuantity
                             groupSize := buy + get
                             
                             // Sort asc
                             sort.Float64s(applicablePrices)
                             
                             // Calculate how many free
                             numGroups := len(applicablePrices) / groupSize
                             numFree := numGroups * get
                             
                             // Sum the cheapest `numFree` items
                             for i := 0; i < numFree; i++ {
                                 itemDiscountAmount += applicablePrices[i]
                             }
                         }
                    }
                } else if d.DiscountType == "free_shipping" {
					faiachunCount := 0
					for _, item := range items {
						if item.ProductType == "faiachun" {
							faiachunCount += item.Quantity
						}
					}
					if faiachunCount >= 4 {
						finalShippingFee = 0.0
					}
                }
            }
        } else {
             fmt.Printf("Error getting discounts: %v\n", err)
        }
    }
    
    discountedSubtotal := subtotal - itemDiscountAmount
    if discountedSubtotal < 0 {
        discountedSubtotal = 0
    }

    // Total calculation
	total := discountedSubtotal + finalShippingFee
    if total < 0 {
        total = 0
    }

	response := gin.H{
		"items":                   items,
		"subtotal":                subtotal,
        "discounted_subtotal":     discountedSubtotal,
        "shipping_fee":            originalShippingFee,
        "discounted_shipping_fee": finalShippingFee,
		"discount":                itemDiscountAmount,
		"total":                   total,
	}

	c.JSON(http.StatusOK, response)
}

// AddItemToCart handles POST /cart/items.
func (h CartHandler) AddItemToCart(c *gin.Context) {
	var req struct {
		ProductID int64   `json:"product_id" binding:"required"`
		SizeType  *string `json:"size_type"`
		Quantity  int     `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartID, err := h.getOrCreateCartID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert size_type string to SizeType enum
	var sizeType *models.SizeType
	if req.SizeType != nil {
		st := models.SizeType(*req.SizeType)
		sizeType = &st
	}

	err = h.Repo.AddItemToCart(c.Request.Context(), cartID, req.ProductID, sizeType, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart", "cart_id": cartID})
}

// UpdateCartItem handles PATCH /cart/items/:cart_item_id.
func (h CartHandler) UpdateCartItem(c *gin.Context) {
	cartItemIDStr := c.Param("cart_item_id")
	cartItemID, err := strconv.ParseInt(cartItemIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	// Security Check: Verify if user owns the cart containing this item
	cartID, err := h.Repo.GetCartIDByCartItemID(c.Request.Context(), cartItemID)
	if err != nil {
		if err.Error() == "cart item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart info"})
		return
	}

	// Check cart ownership
	cart, err := h.Repo.GetCartByID(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	// Authorization Check
	client, isAuthenticated := httpmw.ClientFromContext(c)
	if cart.ClientID != nil {
		if !isAuthenticated || client.ID != *cart.ClientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to modify this cart"})
			return
		}
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Repo.UpdateCartItem(c.Request.Context(), cartItemID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated"})
}

// RemoveCartItem handles DELETE /cart/items/:cart_item_id.
func (h CartHandler) RemoveCartItem(c *gin.Context) {
	cartItemIDStr := c.Param("cart_item_id")
	cartItemID, err := strconv.ParseInt(cartItemIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	// Security Check: Verify if user owns the cart containing this item
	cartID, err := h.Repo.GetCartIDByCartItemID(c.Request.Context(), cartItemID)
	if err != nil {
		if err.Error() == "cart item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart info"})
		return
	}

	// Check cart ownership
	cart, err := h.Repo.GetCartByID(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	// Authorization Check
	client, isAuthenticated := httpmw.ClientFromContext(c)
	if cart.ClientID != nil {
		if !isAuthenticated || client.ID != *cart.ClientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to modify this cart"})
			return
		}
	}

	err = h.Repo.RemoveCartItem(c.Request.Context(), cartItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item removed"})
}

// Checkout handles POST /cart/checkout.
func (h CartHandler) Checkout(c *gin.Context) {
	// Auth required, so client exists
	client, _ := httpmw.ClientFromContext(c)
	cart, err := h.Repo.GetCartByClientID(c.Request.Context(), client.ID)
	if err != nil || cart == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No cart found"})
		return
	}

	var req struct {
		ShippingAddressID *int64 `json:"shipping_address_id"`
		EbuyStoreID       *string `json:"ebuy_store_id"`
		PaymentMethod     string `json:"payment_method" binding:"required"`
		Notes             string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate shipping: one of address or store
	if (req.ShippingAddressID == nil && req.EbuyStoreID == nil) || (req.ShippingAddressID != nil && req.EbuyStoreID != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide either shipping_address_id or ebuy_store_id"})
		return
	}

	// Validate ebuy_store_id if provided
	if req.EbuyStoreID != nil {
		store, err := h.EbuyStoreRepo.GetEbuyStoreByID(c.Request.Context(), *req.EbuyStoreID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate ebuy store"})
			return
		}
		if store == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ebuy_store_id"})
			return
		}
	}

	// Get cart items
	items, err := h.Repo.GetCartItems(c.Request.Context(), cart.CartID)
	if err != nil || len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Calculate totals (simplified, no discounts yet)
	var subtotal float64
	for _, item := range items {
		subtotal += item.UnitPrice * float64(item.Quantity)
	}

	// Create order (need order repository, but for now, placeholder)
	// This would involve creating order, order_items, clearing cart, etc.
	// For now, just clear cart and return success

	err = h.Repo.ClearCart(c.Request.Context(), cart.CartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to checkout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful", "order_id": 123}) // Placeholder
}

// ApplyDiscount handles POST /cart/apply-discount.
func (h CartHandler) ApplyDiscount(c *gin.Context) {
	client, _ := httpmw.ClientFromContext(c)
	cart, err := h.Repo.GetCartByClientID(c.Request.Context(), client.ID)
	if err != nil || cart == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No cart found"})
		return
	}

	var req struct {
		DiscountCode string `json:"discount_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Placeholder: assume discount_id = 1 for any code
	discountID := int64(1)
	err = h.Repo.ApplyDiscountToCart(c.Request.Context(), cart.CartID, discountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply discount"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discount applied"})
}

// RemoveDiscount handles DELETE /cart/discount.
func (h CartHandler) RemoveDiscount(c *gin.Context) {
	client, _ := httpmw.ClientFromContext(c)
	cart, err := h.Repo.GetCartByClientID(c.Request.Context(), client.ID)
	if err != nil || cart == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No cart found"})
		return
	}

	err = h.Repo.RemoveDiscountFromCart(c.Request.Context(), cart.CartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove discount"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discount removed"})
}

func isValidUUID(u string) bool {
	if len(u) != 36 {
		return false
	}
	// simple check for hyphens
	if u[8] != '-' || u[13] != '-' || u[18] != '-' || u[23] != '-' {
		return false
	}
	return true
}
