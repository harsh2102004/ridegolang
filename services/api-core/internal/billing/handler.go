package billing

import (
	"context"
	"log"
	"net/http"

	pb "org/api-core/proto"

	"github.com/gin-gonic/gin"
)

func CreateCheckoutSession(c *gin.Context) {
	userID := c.GetString("user_id")
	log.Println("USER ID FROM MIDDLEWARE:", userID)

	var body struct {
		Plan string `json:"plan"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	if body.Plan == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "plan is required",
		})
		return
	}

	res, err := Client.CreateCheckoutSession(
		context.Background(),
		&pb.CreateCheckoutSessionRequest{
			UserId: userID,
			Plan:   body.Plan,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkout_url": res.CheckoutUrl,
	})
}



func GetSubscription(c *gin.Context) {
	userID := c.GetString("user_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	plan, status, err := GetUserSubscription(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plan":   plan,
		"status": status,
	})
}


func CreatePortalSession(c *gin.Context) {
	userID := c.GetString("user_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	res, err := Client.CreatePortalSession(
		context.Background(),
		&pb.CreatePortalSessionRequest{
			UserId: userID,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"portal_url": res.PortalUrl,
	})
}