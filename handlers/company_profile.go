package handlers

import (
	"net/http"
	"strconv"

	"inventory_system/database"
	"inventory_system/models"

	"github.com/gin-gonic/gin"
)

// CompanyProfileRequest represents the request body for company profile operations
type CompanyProfileRequest struct {
	CompanyName     string `json:"company_name" binding:"required"`
	CompanyAddress  string `json:"company_address"`
	CompanyPhone    string `json:"company_phone"`
	CompanyEmail    string `json:"company_email"`
	CompanyWebsite  string `json:"company_website"`
	TaxNumber       string `json:"tax_number"`
	BusinessLicense string `json:"business_license"`
	LogoBase64      string `json:"logo_base64"`
	InvoiceFooter   string `json:"invoice_footer"`
	BankAccount     string `json:"bank_account"`
	Currency        string `json:"currency"`
}

// GetCompanyProfile returns the active company profile
func GetCompanyProfile(c *gin.Context) {
	var profile models.CompanyProfile
	result := database.DB.Where("is_active = ?", true).First(&profile)
	
	if result.Error != nil {
		// If no profile exists, return default structure
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusOK, gin.H{
				"company_name":     "INVENTORY SYSTEM",
				"company_address":  "",
				"company_phone":    "",
				"company_email":    "",
				"company_website":  "",
				"tax_number":       "",
				"business_license": "",
				"logo_base64":      "",
				"invoice_footer":   "Thank you for your business!",
				"bank_account":     "",
				"currency":         "IDR",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// CreateOrUpdateCompanyProfile creates a new company profile or updates existing one
func CreateOrUpdateCompanyProfile(c *gin.Context) {
	var req CompanyProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "IDR"
	}

	// Start transaction
	tx := database.DB.Begin()

	// Check if a profile already exists
	var existingProfile models.CompanyProfile
	result := tx.Where("is_active = ?", true).First(&existingProfile)

	if result.Error != nil && result.Error.Error() != "record not found" {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing profile"})
		return
	}

	if result.Error == nil {
		// Update existing profile
		existingProfile.CompanyName = req.CompanyName
		existingProfile.CompanyAddress = req.CompanyAddress
		existingProfile.CompanyPhone = req.CompanyPhone
		existingProfile.CompanyEmail = req.CompanyEmail
		existingProfile.CompanyWebsite = req.CompanyWebsite
		existingProfile.TaxNumber = req.TaxNumber
		existingProfile.BusinessLicense = req.BusinessLicense
		existingProfile.LogoBase64 = req.LogoBase64
		existingProfile.InvoiceFooter = req.InvoiceFooter
		existingProfile.BankAccount = req.BankAccount
		existingProfile.Currency = req.Currency

		if err := tx.Save(&existingProfile).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company profile"})
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{
			"message": "Company profile updated successfully",
			"data":    existingProfile,
		})
	} else {
		// Create new profile
		newProfile := models.CompanyProfile{
			CompanyName:     req.CompanyName,
			CompanyAddress:  req.CompanyAddress,
			CompanyPhone:    req.CompanyPhone,
			CompanyEmail:    req.CompanyEmail,
			CompanyWebsite:  req.CompanyWebsite,
			TaxNumber:       req.TaxNumber,
			BusinessLicense: req.BusinessLicense,
			LogoBase64:      req.LogoBase64,
			InvoiceFooter:   req.InvoiceFooter,
			BankAccount:     req.BankAccount,
			Currency:        req.Currency,
			IsActive:        true,
		}

		if err := tx.Create(&newProfile).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company profile"})
			return
		}

		tx.Commit()
		c.JSON(http.StatusCreated, gin.H{
			"message": "Company profile created successfully",
			"data":    newProfile,
		})
	}
}

// UpdateCompanyLogo updates only the company logo
func UpdateCompanyLogo(c *gin.Context) {
	var req struct {
		LogoBase64 string `json:"logo_base64" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.CompanyProfile
	result := database.DB.Where("is_active = ?", true).First(&profile)
	
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company profile not found"})
		return
	}

	profile.LogoBase64 = req.LogoBase64
	if err := database.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update logo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logo updated successfully",
		"data":    profile,
	})
}

// DeleteCompanyProfile soft deletes the company profile
func DeleteCompanyProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	var profile models.CompanyProfile
	result := database.DB.First(&profile, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company profile not found"})
		return
	}

	if err := database.DB.Delete(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company profile deleted successfully"})
}

// GetCompanyProfileForInvoice returns simplified company data for invoice generation
func GetCompanyProfileForInvoice(c *gin.Context) {
	var profile models.CompanyProfile
	result := database.DB.Where("is_active = ?", true).First(&profile)
	
	if result.Error != nil {
		// Return default values if no profile exists
		c.JSON(http.StatusOK, gin.H{
			"company_name":    "INVENTORY SYSTEM",
			"company_address": "",
			"company_phone":   "",
			"company_email":   "",
			"logo_base64":     "",
			"invoice_footer":  "Thank you for your business!",
			"currency":        "IDR",
		})
		return
	}

	// Return only invoice-relevant data
	c.JSON(http.StatusOK, gin.H{
		"company_name":    profile.CompanyName,
		"company_address": profile.CompanyAddress,
		"company_phone":   profile.CompanyPhone,
		"company_email":   profile.CompanyEmail,
		"logo_base64":     profile.LogoBase64,
		"invoice_footer":  profile.InvoiceFooter,
		"currency":        profile.Currency,
	})
}