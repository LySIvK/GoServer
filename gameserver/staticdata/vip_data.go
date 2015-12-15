package staticdata

type Vip_Data struct {
	ID        int
	Level     int
	FreeMoney int
	Discount  float32
	Icon      int
}

func (self *Vip_Data) GetName() string {
	return "vip"
}

func (self *Vip_Data) GetFilePath() string {
	return "csv/vip_data.csv"
}
