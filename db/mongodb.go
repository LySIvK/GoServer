package db

import (
	"errors"
	"gopath/gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"loger"
)

var (
	db_Url        string             //! MongoDB连接Url
	db_Connection *mgo.Session = nil //! MongoDB连接Session
)

//! 初始化Mongodb连接
func Init(url string, limit int) {
	db_Url = url

	var err error
	db_Connection, err = mgo.Dial(db_Url)
	if err != nil {
		//! 连接Mongodb失败
		loger.Fatal("connect mongodb failed! error: %v", err.Error())
	}

	//! 设置Mongodb连接池限制 默认值为4096
	db_Connection.SetPoolLimit(limit)
}

//! 使用账号密码初始化Mongodb连接
func InitWithAccount(url string, username string, password string, limit int) {
	db_Url = url

	var err error
	db_Connection, err = mgo.DialWithInfo(&mgo.DialInfo{Addrs: []string{db_Url}, Username: username, Password: password, PoolLimit: limit})
	if err != nil {
		//! 连接Mongodb失败
		loger.Fatal("connect mongodb failed! error: %v", err.Error())
	}
}

//! 获取Mongodb连接的Url
func GetDBUrl() string {
	return db_Url
}

//! 获取Mongodb连接
func GetDBSession() *mgo.Session {
	if db_Connection == nil {
		loger.Fatal("db_Connection is nil !")
	}

	return db_Connection.Clone()
}

//! 插入一条数据
func Insert(dbName string, tableName string, data interface{}) error {
	db_session := GetDBSession()
	defer db_session.Close()

	//! 获取指定数据库表Collection对象
	collection := db_session.DB(dbName).C(tableName)
	err := collection.Insert(data)
	if err != nil {
		if mgo.IsDup(err) == true {
			loger.Warn("Insert error: %v \r\ndbName: %s \r\ntable: %s \r\ndata: %v",
				err.Error(), dbName, tableName, data)
		} else {
			loger.Error("Insert error: %v \r\ndbName: %s \r\ntable: %s \r\ndata: %v",
				err.Error(), dbName, tableName, data)
		}
		return err
	}

	return nil
}

//! 增加一个字段值
func IncFieldValue(dbName string, tableName string, find string, find_value interface{}, filed string, value int) error {
	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Update(bson.M{find: find_value}, bson.M{"$inc": bson.M{filed: value}})
	if err != nil {
		loger.Error("IncFieldValue error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%v \r\n value: %s:%d \r\n",
			err.Error(), dbName, tableName, find, find_value, filed, value)
		return err
	}

	return nil
}

//! 获取集合中元素数量
func Count(dbName string, tableName string) (int, error) {
	db_session := GetDBSession()
	defer db_session.Close()

	return db_session.DB(dbName).C(tableName).Count()
}

//! 查询一条数据
func Find(dbName string, tableName string, find string, find_value interface{}, data interface{}) error {
	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Find(bson.M{find: find_value}).One(data)
	if err != nil {
		loger.Error("Find error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%v \r\n",
			err.Error(), dbName, tableName, find, find_value)
	}

	return err
}

//! 条件查询
//! mongodb.Find_Conditional("TestDB", "TestTable", "money", ">=", 999999999, &player)
func Find_Conditional(dbName string, tableName string, find string, conditional string, find_value interface{}, lst interface{}) error {
	db_session := GetDBSession()
	defer db_session.Close()

	var con string

	switch conditional {
	case ">":
		con = "$gt" //! greater than
	case "<":
		con = "$lt" //! less than
	case ">=":
		con = "$gte" //! greater than equl
	case "<=":
		con = "$lte"
	case "!=":
		con = "$ne"
	default:
		err := errors.New("conditional wrong")
		loger.Error("Find_Conditional error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%s%v \r\n",
			err.Error(), dbName, tableName, find, conditional, find_value)
		return err
	}

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Find(bson.M{find: bson.M{con: find_value}}).All(lst)
	if err != nil {
		loger.Error("Find_Conditional error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%s%v \r\n",
			err.Error(), dbName, tableName, find, conditional, find_value)
	}

	return err
}

//! 范围查询
//! mongodb.Find_Range("TestDB", "TestTable", "level", 10, 100, true, &player)
func Find_Range(dbName string, tableName string, find string, range_begin interface{}, range_end interface{}, isEqul bool, lst interface{}) error {
	var greater, less string

	if isEqul == true {
		greater = "$gte"
		less = "$lte"
	} else {
		greater = "$gt"
		less = "$lt"
	}

	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Find(bson.M{find: bson.M{greater: range_begin, less: range_end}}).All(lst)
	if err != nil {
		loger.Error("Find_Range error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%v --- %v \r\n",
			err.Error(), dbName, tableName, find, range_begin, range_end)
	}

	return err
}

//! 排序查找
//! order 1 -> 正序  -1 -> 倒序
//! Mongodb语句: db.getCollection('TestTable').find().sort({'money' : -1}).limit(1)
func Find_Sort(dbName string, tableName string, find string, order int, number int, lst interface{}) error {
	db_session := GetDBSession()
	defer db_session.Close()

	strSort := ""
	if order == 1 {
		strSort = "+" + find
	} else {
		strSort = "-" + find
	}

	collection := db_session.DB(dbName).C(tableName)
	query := collection.Find(nil).Sort(strSort).Limit(number)

	err := query.All(lst)
	if err != nil {
		loger.Error("Find_Sort error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s \r\norder: %d\r\nlimit: %d\r\n",
			err.Error(), dbName, tableName, find, order, number)
	}

	return err
}

//! 查询所有数据
func FindAll(dbName string, tableName string, lst interface{}) error {
	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Find(nil).Iter().All(lst)
	if err != nil {
		loger.Error("FindAll error: %v \r\ndbName: %s \r\ntable: %s \r\n",
			err.Error(), dbName, tableName)
	}

	return err
}

//! 更新字段值
func UpdateField(dbName string, tableName string, find string, find_value interface{}, update string, update_value interface{}) error {
	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Update(bson.M{find: find_value}, bson.M{"$set": bson.M{update: update_value}})
	if err != nil {
		loger.Error("UpdateField error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%v \r\nupdate: %s:%v\r\n",
			err.Error(), dbName, tableName, find, find_value, update, update_value)
	}

	return err
}

//! 更新整个字段
func Update(dbName string, tableName string, find string, find_value interface{}, data interface{}) error {
	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Update(bson.M{find: find_value}, data)
	if err != nil {
		loger.Error("Update error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%v \r\nupdate: %v\r\n",
			err.Error(), dbName, tableName, find, find_value, data)
	}

	return err
}

//! 删除一条记录
func Remove(dbName string, tableName string, find string, find_value interface{}) error {
	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	err := collection.Remove(bson.M{find: find_value})
	if err != nil {
		loger.Error("Remove error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%v \r\n",
			err.Error(), dbName, tableName, find, find_value)
	}

	return err
}

//! 删除所有记录
func RemoveAll(dbName string, tableName string, find string, find_value interface{}) (*mgo.ChangeInfo, error) {
	db_session := GetDBSession()
	defer db_session.Close()

	collection := db_session.DB(dbName).C(tableName)
	info, err := collection.RemoveAll(bson.M{find: find_value})
	if err != nil {
		loger.Error("RemoveAll error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s:%v \r\n",
			err.Error(), dbName, tableName, find, find_value)
	}

	return info, err
}

//! 增加数组字段
func AddToArray(dbname string, tableName string, find string, find_value interface{}, fieldname string, data interface{}) error {
	return Update(dbname, tableName, find, find_value, bson.M{"$push": bson.M{fieldname: data}})
}

//删掉数组字段
func RemoveFromArray(dbname string, tableName string, find string, find_value interface{}, fieldname string, data interface{}) error {
	return Update(dbname, tableName, find, find_value, bson.M{"$pull": bson.M{fieldname: data}})
}
