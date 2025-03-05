package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"

	"gorm.io/gorm"
)

type AcademicNavigatorRepository interface {
	AddAcademicNavigator(academicNavigator *model.AcademicNavigator) error
	GetAcademicNavigator(academicNavigator *model.AcademicNavigator, page int64, pageSize int64) ([]*model.AcademicNavigator, int64, error)
	GetAcademicNavigatorByID(ID int64) (*model.AcademicNavigator, error)
	DeleteAcademicNavigator(ID int64) error
}

type AcademicNavigatorRepositoryImpl struct {
	DB            *gorm.DB
	Elasticsearch storage.ElasticsearchImpl
}

// NewAcademicNavigatorRepositoryImpl func NewAcademicNavigatorRepositoryImpl(db *gorm.DB, cfg *elasticsearch.Config) *AcademicNavigatorRepositoryImpl {
func NewAcademicNavigatorRepositoryImpl(db *gorm.DB) *AcademicNavigatorRepositoryImpl {
	//elastic, err := storage.NewElasticSearchClient(cfg)
	//if err != nil {
	//	return nil
	//}
	return &AcademicNavigatorRepositoryImpl{
		DB: db,
		//Elasticsearch: elastic,
	}
}

func (r *AcademicNavigatorRepositoryImpl) AddAcademicNavigator(academicNavigator *model.AcademicNavigator) error {
	return r.DB.Create(academicNavigator).Error
}

func (r *AcademicNavigatorRepositoryImpl) GetAcademicNavigator(academicNavigator *model.AcademicNavigator, page int64, pageSize int64) ([]*model.AcademicNavigator, int64, error) {
	var total int64
	var academicNavigators []*model.AcademicNavigator
	// 先从Elasticsearch中查询到符合条件的ID
	//ids, err := r.Elasticsearch.SearchByFields("academic_navigators", map[string]interface{}{
	//	"user_id":   academicNavigator.UserId,
	//	"content":   academicNavigator.Content,
	//	"education": academicNavigator.Education,
	//	"major":     academicNavigator.Major,
	//	"school":    academicNavigator.School,
	//})
	//if err != nil {
	//	return nil, 0, err
	//}
	// 再根据ID从数据库中查询到对应的数据
	//err = r.DB.Model(&model.AcademicNavigator{}).Where("id IN (?)", ids).Find(&academicNavigators).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Count(&total).Error
	//if err != nil {
	//	return nil, 0, err
	//}
	//return academicNavigators, total, nil

	err := r.DB.Model(&model.AcademicNavigator{}).Where("user_id = ? AND content = ? AND education = ? AND major = ? AND school = ?", academicNavigator.UserId, academicNavigator.Content, academicNavigator.Education, academicNavigator.Major, academicNavigator.School).Find(&academicNavigators).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Count(&total).Error
	return academicNavigators, total, err
}

func (r *AcademicNavigatorRepositoryImpl) GetAcademicNavigatorByID(ID int64) (*model.AcademicNavigator, error) {
	var academicNavigator model.AcademicNavigator
	err := r.DB.First(&academicNavigator, ID).Error
	return &academicNavigator, err
}

func (r *AcademicNavigatorRepositoryImpl) DeleteAcademicNavigator(ID int64) error {
	return r.DB.Delete(&model.AcademicNavigator{}, ID).Error
}
