package staticdata

type Mall_Data struct {
	ID       int
	Name     string
	Money    int
	Discount float32
	Icon     int
}

func (self *Mall_Data) GetName() string {
	return "mall"
}

func (self *Mall_Data) GetFilePath() string {
	return "csv/mall_data.csv"
}
