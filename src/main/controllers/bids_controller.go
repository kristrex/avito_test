package controllers

import (
	"net/http"
	"strconv"
	"time"

	"zadanie-6105/src/main/models"
	"zadanie-6105/src/main/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBids создает новое предложение
func CreateBids(c *gin.Context) {
	// Декодирование тела запроса
	var req models.Bids
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}
	//Проверка типа(при создании он может быть только organization)
	if req.AuthorType != "Organization" {
		c.JSON(http.StatusForbidden, gin.H{"reason": "Invalid author type"})
		return
	}
	//Проверяем существует ли тендер
	exists_tender, _ := utils.TenderExists(models.Db, req.TenderID)
	if req.TenderID == "" || !exists_tender {
		c.JSON(http.StatusNotFound, gin.H{"reason": "Tender not found"})
		return
	}
	//Проверка существования пользователя
	exists, err_usr := utils.UserIdExists(models.Db, req.AuthorID)
	if err_usr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": err_usr.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "User not found"})
		return
	} else {
		var org_resp models.OrganizationResponsible
		result_org := models.Db.Table("organization_responsible").Where("user_id = ?", req.AuthorID).First(&org_resp)
		if result_org.Error != nil {
			c.JSON(http.StatusForbidden, gin.H{"reason": "User acess denied"})
			return
		}
		result_tndr := models.Db.Table("tender").Where("organizationid = ?", org_resp.OrganizationID)
		if result_tndr.Error != nil {
			c.JSON(http.StatusForbidden, gin.H{"reason": "User acess denied"})
			return
		}
	}
	//Подготавливаем данные для вставки в БД
	create := models.Bids{
		Name:        req.Name,
		Description: req.Description,
		TenderID:    req.TenderID,
		AuthorType:  req.AuthorType,
		AuthorID:    req.AuthorID,
	}
	create.ID = uuid.New().String()
	create.Version = 1
	create.Status = "Created"
	// Вставка данных в базу данных
	tx := models.Db.Table("bid").Create(&create)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Data insertion error!"})
		return
	}

	// Подготовка ответа
	var savedBids models.Bids
	if err := models.Db.Table("bid").First(&savedBids, "id = ?", create.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	// Подготовка ответа
	resp := models.BidsResponse{
		ID:         savedBids.ID,
		Name:       savedBids.Name,
		Status:     savedBids.Status,
		AuthorType: savedBids.AuthorType,
		AuthorID:   savedBids.AuthorID,
		Version:    savedBids.Version,
		CreatedAt:  savedBids.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}

// MyBids показывает предложения пользователя
func MyBids(c *gin.Context) {
	//Получаем параметр - имя пользователя
	usr_name := c.Query("username")
	//Проверяем существует ли пользователь
	exists, err_usr := utils.UserExists(models.Db, usr_name)
	if err_usr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": err_usr.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "User not found"})
		return
	}
	// Получаем параметр ограничения количества выводимых данных
	limitStr := c.DefaultQuery("limit", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid limit parameter"})
		return
	}
	//Получаем параметр смещения
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid offset parameter"})
		return
	}

	var bid []models.Bids
	query := models.Db.Table("bid").Find(&bid)
	// Формируем запрос для получения тендеров у пользователя
	if usr_name != "" {
		var employee models.Employee
		models.Db.Table("employee").Where("username = ?", usr_name).First(&employee)
		query = query.Where("authorid = ?", employee.ID)
	}
	//Реализуем сортировку по алфавитному порядку(вне зависимости от регистра)
	query = query.Order("LOWER(name) ASC")
	// Ограничиваем количество выводимых данных, если параметр limit указан и больше 0
	if limit > 0 {
		query = query.Limit(limit)
	}
	//Реализуем смещение
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Database error"})
		return
	}

	// Проверяем, есть ли данные для возврата
	if len(bid) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	// Возвращаем тендеры без organizationId
	var bidResponses []models.BidsResponse
	for _, bds := range bid {
		bidResponse := models.BidsResponse{
			ID:         bds.ID,
			Name:       bds.Name,
			Status:     bds.Status,
			AuthorType: bds.AuthorType,
			AuthorID:   bds.AuthorID,
			Version:    bds.Version,
			CreatedAt:  bds.CreatedAt.Format(time.RFC3339),
		}
		bidResponses = append(bidResponses, bidResponse)
	}
	c.JSON(http.StatusOK, bidResponses)
}

// BidsList публикует предложение
func BidsList(c *gin.Context) {
	param := c.Param("tenderId")
	//Проверяем существует ли tender c данным id
	exists, _ := utils.TenderExists(models.Db, param)
	if param == "" || !exists {
		c.JSON(http.StatusNotFound, gin.H{"reason": "Tender not found"})
		return
	}
	// Проверяем существует ли пользователь с таким именем
	usr_name := c.Query("username")
	exists_usr, _ := utils.UserExists(models.Db, usr_name)
	if !exists_usr {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "User not found"})
		return
	}
	// Проверяем права
	var employee models.Employee
	models.Db.Table("employee").Where("username = ?", usr_name).First(&employee)
	var tender_resp models.Tender
	models.Db.Table("tender").Where("id = ?", param).First(&tender_resp)
	var org_resp models.OrganizationResponsible
	tx := models.Db.Table("organization_responsible").Where("user_id = ? AND organization_id = ?", employee.ID, tender_resp.OrganizationID).First(&org_resp)
	if tx.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{"reason": "User access denied"})
		return
	}
	//Получаем параметр: ограничение выводимых данных
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid limit parameter"})
		return
	}
	//Получаем параметр: смещение
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid offset parameter"})
		return
	}

	var bid []models.Bids
	query := models.Db.Table("bid").Find(&bid)
	// Формируем запрос для получения тендеров у пользователя
	if usr_name != "" {
		var employee models.Employee
		models.Db.Table("employee").Where("username = ?", usr_name).First(&employee)
		var bids_resp models.BidsResponse
		models.Db.Table("bid").Where("authorid = ? AND tenderid = ?", employee.ID, param).First(&bids_resp)
		query = query.Where("id = ?", bids_resp.ID)
	}
	//Реализуем сортировку по алфавитному порядку(вне зависимости от регистра)
	query = query.Order("LOWER(name) ASC")

	// Ограничиваем количество выводимых данных, если параметр limit указан и больше 0
	if limit > 0 {
		query = query.Limit(limit)
	}
	// Реализуем смещение
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Database error"})
		return
	}

	// Проверяем, есть ли данные для возврата
	if len(bid) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	// Возвращаем тендеры без organizationId
	var bidResponses []models.BidsResponse
	for _, bds := range bid {
		bidResponse := models.BidsResponse{
			ID:         bds.ID,
			Name:       bds.Name,
			Status:     bds.Status,
			AuthorType: bds.AuthorType,
			AuthorID:   bds.AuthorID,
			Version:    bds.Version,
			CreatedAt:  bds.CreatedAt.Format(time.RFC3339),
		}
		bidResponses = append(bidResponses, bidResponse)
	}
	c.JSON(http.StatusOK, bidResponses)
}
