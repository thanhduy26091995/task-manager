package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"task-manager-api/utils"

	"github.com/gin-gonic/gin"
)

func HandleAppleWebhook(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Failed to read request body", err.Error())
		return
	}

	// Get the X-Apple-Signature header
	rawSig := strings.TrimSpace(c.GetHeader("X-Apple-Signature"))
	if rawSig == "" {
		utils.Error(c, http.StatusBadRequest, "Missing X-Apple-Signature header", "The X-Apple-Signature header is required for validation")
		return
	}

	// Validate the signature format
	prefix := "hmacsha256="
	if !strings.HasPrefix(rawSig, prefix) {
		utils.Error(c, http.StatusBadRequest, "Invalid signature format", "The X-Apple-Signature header must start with 'hmacsha256='")
		return
	}

	receivedSig := strings.TrimPrefix(rawSig, prefix)

	// Check if the APPLE_WEBHOOK_SECRET environment variable is set
	secret := os.Getenv("APPLE_WEBHOOK_SECRET")
	if secret == "" {
		utils.Error(c, http.StatusInternalServerError, "Missing APPLE_WEBHOOK_SECRET", "The APPLE_WEBHOOK_SECRET environment variable must be set for signature validation")
		return
	}

	// Compute expected HMAC-SHA256 signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(bodyBytes)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	if !hmac.Equal([]byte(receivedSig), []byte(expectedSig)) {
		utils.Error(c, http.StatusUnauthorized, "Invalid signature", "The signature in the X-Apple-Signature header does not match the computed signature")
		return
	}

	// Parse the JSON payload
	var payload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	fmt.Println("Received Apple webhook payload:", payload)
	utils.Success(c, http.StatusOK, payload)
}
