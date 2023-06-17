package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"smart-contract-service/app/usecase"
	"smart-contract-service/models"
	"time"
)

type DatabaseConnection struct {
	client *gorm.DB
}

func NewDatabaseConnection(client *gorm.DB) usecase.DbRepository {
	return &DatabaseConnection{client: client}
}

func (db *DatabaseConnection) GetCustomerData(id string) (data *models.Customer, err error) {
	err = db.client.Model(&models.Customer{}).Where("id = ?", id).Find(&data).Error
	return
}

func (db *DatabaseConnection) GetCustomerByAccount(account string) (data *models.Customer, err error) {
	err = db.client.Model(&models.Customer{}).Where("account = ?", account).Find(&data).Error
	return
}

func (db *DatabaseConnection) GetUserById(id string) (data *models.Partners, err error) {
	err = db.client.Model(&models.Partners{}).Where("id = ?", id).Find(&data).Error
	return
}

func (db *DatabaseConnection) GetUserByUsername(username string) (data *models.Partners, err error) {
	err = db.client.Model(&models.Partners{}).Where("username = ?", username).Find(&data).Error
	return
}

func (db *DatabaseConnection) GetUserByReferenceNo(referenceNo string) (data *models.Partners, err error) {
	err = db.client.Model(&models.Partners{}).Where("reference_no = ?", referenceNo).Find(&data).Error
	return
}

func (db *DatabaseConnection) InsertUser(input *models.Partners) (id string, err error) {
	id = uuid.New().String()
	timeNow := time.Now()
	err = db.client.Create(&models.Partners{
		Id:        id,
		Username:  input.Username,
		Password:  input.Password,
		CreatedAt: &timeNow,
		UpdatedAt: nil,
	}).Error
	return id, err
}

func (db *DatabaseConnection) InsertPayment(input *models.Payment) (id string, err error) {
	id = uuid.New().String()
	timeNow := time.Now()
	err = db.client.Create(&models.Payment{
		Id:             id,
		PartnerId:      input.PartnerId,
		ConsumerId:     input.ConsumerId,
		Amount:         input.Amount,
		Currency:       input.Currency,
		AdditionalInfo: input.AdditionalInfo,
		CreatedAt:      &timeNow,
		UpdatedAt:      nil,
	}).Error
	return id, err
}
