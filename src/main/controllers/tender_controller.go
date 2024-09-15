package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"zadanie-6105/src/main/models"
	"zadanie-6105/src/main/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TenderHandler(c *gin.Context) {
	// Получаем параметр фильтрации по типу услуг
	serviceType := c.Query("service_type")
	if serviceType != "Delivery" && serviceType != "Construction" && serviceType != "Manufacture" && serviceType != "" {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid service_type parameter"})
		return
	}
	// Получаем параметр ограничения количества выводимых данных
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid limit parameter"})
		return
	}
	// Получаем параметр смещения
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid offset parameter"})
		return
	}

	// Создаем запрос к базе
	var tender []models.Tender
	//query := models.Db.Model(&models.Tender{})
	query := models.Db.Table("tender").Find(&tender)
	// Фильтруем тендеры по типу услуг, если параметр указан
	if serviceType != "" {
		query = query.Where("servicetype = ?", serviceType)
	}
	query = query.Order("LOWER(name) ASC")
	// Ограничиваем количество выводимых данных, если параметр limit указан и больше 0
	if limit > 0 {
		query = query.Limit(limit)
	}
	// Реализуем смещение
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Database error"})
		return
	}

	// Проверяем, есть ли данные для возврата
	if len(tender) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	// Возвращаем тендеры без organizationId
	var tenderResponses []models.TenderResponse
	for _, tndr := range tender {
		tenderResponse := models.TenderResponse{
			ID:          tndr.ID,
			Name:        tndr.Name,
			Description: tndr.Description,
			ServiceType: tndr.ServiceType,
			Status:      tndr.Status,
			Version:     tndr.Version,
			CreatedAt:   tndr.CreatedAt.Format(time.RFC3339),
		}
		tenderResponses = append(tenderResponses, tenderResponse)
	}

	// Возвращаем отфильтрованные тендеры в формате JSON
	c.JSON(http.StatusOK, tenderResponses)
}

// CreateTender создает новый тендер
func CreateTender(c *gin.Context) {
	// Декодирование тела запроса
	var req models.CreationTender
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}
	//Проверка типа тендера
	if req.ServiceType != "Delivery" && req.ServiceType != "Construction" && req.ServiceType != "Manufacture" {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid serviceType parameter"})
		return
	}
	//Проверка существования пользователя
	exists, err_usr := utils.UserExists(models.Db, req.CreatorUsername)
	if err_usr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": err_usr.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "User not found"})
		return
	} else {
		var employee models.Employee
		result := models.Db.Table("employee").Where("username = ?", req.CreatorUsername).First(&employee)
		fmt.Println(employee.ID, employee.Username, result.Error)
		var org_resp models.OrganizationResponsible
		result_org := models.Db.Table("organization_responsible").Where("user_id = ? AND organization_id = ?", employee.ID, req.OrganizationID).First(&org_resp)
		if result_org.Error != nil {
			c.JSON(http.StatusForbidden, gin.H{"reason": "User acess denied"})
			return
		}
	}
	//Подготавливаем данные для вставки
	create := models.Tender{
		Name:           req.Name,
		Description:    req.Description,
		ServiceType:    req.ServiceType,
		OrganizationID: req.OrganizationID,
	}
	create.ID = uuid.New().String()
	create.Version = 1
	create.Status = "Created"
	// Вставка данных в базу данных
	tx := models.Db.Table("tender").Create(&create)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Data insertion error!"})
		return
	}

	// Подготовка ответа
	var savedTender models.Tender
	if err := models.Db.Table("tender").First(&savedTender, "name = ?", req.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	resp := models.TenderResponse{
		ID:          savedTender.ID,
		Name:        savedTender.Name,
		Description: savedTender.Description,
		Status:      savedTender.Status,
		ServiceType: savedTender.ServiceType,
		Version:     savedTender.Version,
		CreatedAt:   savedTender.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}

// MyTender показывает тендеры пользователя
func MyTender(c *gin.Context) {
	//Получаем параметр: имя пользователя
	usr_name := c.Query("username")
	//Проверка имени пользователя(идентификация)
	exists, err_usr := utils.UserExists(models.Db, usr_name)
	if err_usr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": err_usr.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "User not found"})
		return
	}
	// Получаем и проверяем параметр ограничения количества выводимых данных
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid limit parameter"})
		return
	}
	//Получаем и проверяем параметр: смещение
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid offset parameter"})
		return
	}

	var tender []models.Tender
	query := models.Db.Table("tender").Find(&tender)
	// Формируем запрос для получения тендеров у пользователя
	if usr_name != "" {
		var employee models.Employee
		models.Db.Table("employee").Where("username = ?", usr_name).First(&employee)
		var org_resp models.OrganizationResponsible
		models.Db.Table("organization_responsible").Where("user_id = ?", employee.ID).First(&org_resp)
		query = query.Where("organizationid = ?", org_resp.OrganizationID)
	}
	//Реализуем сортировку по алфавитному порядку( вне зависимости от регистра)
	query = query.Order("LOWER(name) ASC")
	// Ограничиваем количество выводимых данных, если параметр limit указан и больше 0
	if limit > 0 {
		query = query.Limit(limit)
	}
	// Реализуем смещенеие
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Database error"})
		return
	}

	// Проверяем, есть ли данные для возврата
	if len(tender) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	// Возвращаем тендеры без organizationId
	var tenderResponses []models.TenderResponse
	for _, tndr := range tender {
		tenderResponse := models.TenderResponse{
			ID:          tndr.ID,
			Name:        tndr.Name,
			Description: tndr.Description,
			ServiceType: tndr.ServiceType,
			Status:      tndr.Status,
			Version:     tndr.Version,
			CreatedAt:   tndr.CreatedAt.Format(time.RFC3339),
		}
		tenderResponses = append(tenderResponses, tenderResponse)
	}
	c.JSON(http.StatusOK, tenderResponses)
}

// Получение статуса тендера
func GetStatus(c *gin.Context) {
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
	//Проверяем есть ли права у пользователя
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

	c.JSON(http.StatusOK, tender_resp.Status)
}

// PutStatus меняем значение статуса
func PutStatus(c *gin.Context) {
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
	statusStr := c.Query("status")
	if statusStr != "Created" && statusStr != "Published" && statusStr != "Closed" {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid status parameter"})
		return
	}
	//Проверяем есть ли права у пользователя
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
	// Реализуем смену статуса
	query := `UPDATE tender SET status = $1 WHERE id = $2`
	models.Db.Exec(query, statusStr, tender_resp.ID)
	// Подготавливаем ответ
	models.Db.Table("tender").Where("id = ?", param).First(&tender_resp)
	tenderResponse := models.TenderResponse{
		ID:          tender_resp.ID,
		Name:        tender_resp.Name,
		Description: tender_resp.Description,
		ServiceType: tender_resp.ServiceType,
		Status:      tender_resp.Status,
		Version:     tender_resp.Version,
		CreatedAt:   tender_resp.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, tenderResponse)
}
