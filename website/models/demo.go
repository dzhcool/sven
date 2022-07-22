package model

type Demo struct {
	Id    int    `json:"id" xorm:"id"`
	Title string `json:"title" xorm:"title"`
}

func (Demo) TableName() string {
	return "t_demo"
}

type demoModel struct {
	model
}

func NewDemoModel() *demoModel {
	return new(demoModel)
}

func (p *demoModel) Get(id int) (*Demo, error) {
	var row Demo
	if err := db.Where("id = ? ", id).First(&row).Error; err != nil {
		return nil, err
	}

	return &row, nil
}
