package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var DB *sql.DB
var Dsn string
var DBName string
var OutDir string
var InTpl string
var Tables []string
var GoPackage string
var Package string

func init() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	Dsn = viper.GetString("db.dsn")
	DBName = viper.GetString("db.db_name")
	OutDir = viper.GetString("config.out_dir")
	if OutDir == "" {
		OutDir = "message"
	}
	GoPackage = viper.GetString("config.go_package")
	if GoPackage == "" {
		GoPackage = OutDir
	}	
	InTpl = viper.GetString("config.in_tpl")
	if InTpl == "" {
		InTpl = "proto.tpl"
	}
	Package = viper.GetString("config.package")

	Tables = viper.GetStringSlice("db.db_tables")
}

func main() {
	Run(DBName)
}

func Run(dbName string) {
	db, err := Connect("mysql", Dsn)
	if err != nil {
		log.Fatal("db连接失败")
	}
	DB = db

	all := GetTables(dbName)

	if len(Tables) > 0 {
		var b []string
		for _, v := range Tables {
			for _, v2 := range all {
				if v2 == v {
					b = append(b, v2)
					break
				}
			}
		}
		all = b
	}

	if len(all) == 0 {
		return
	}

	var messages []Message
	for _, v := range all {
		var one Message
		one.Name = v
		all := GetStruct(v, dbName)
		for k, v2 := range all {
			all[k].Type = TypeMToP(v2.Type)
		}
		one.MessageDetail = all
		messages = append(messages, one)
	}
	Genarate(OutDir, messages)

}

func Genarate(dir string, all []Message) {
	tmpl, err := template.ParseFiles(InTpl)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stat(dir)
	if err != nil {
		// 创建文件夹
		os.MkdirAll(dir, 0755)
	}
	// 循环，创建具体得文件
	for _, v := range all {
		filepath := fmt.Sprintf("%v/%v.proto", dir, v.Name)
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		v.GoPackage = GoPackage
		v.Package = Package
		err = tmpl.Execute(file, v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetTables(dbName string) []string {
	rows, err := DB.Query(fmt.Sprintf(`SELECT TABLE_NAME FROM information_schema.TABLES t WHERE  t.TABLE_SCHEMA = "%v"`, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var all []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		all = append(all, name)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return all
}

func GetStruct(dbTable, dbName string) []TableFied {
	rows, err := DB.Query("SELECT c.COLUMN_NAME,c.COLUMN_TYPE,c.COLUMN_COMMENT FROM INFORMATION_SCHEMA.Columns c WHERE c.`TABLE_SCHEMA`='" + dbName + "' AND c.TABLE_NAME = '" + dbTable + "' ORDER BY c.ORDINAL_POSITION" )
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var all []TableFied
	var num = 1
	for rows.Next() {
		var one TableFied
		rows.Scan(&one.Name, &one.Type, &one.Comment)
		one.Num = num

		n := strings.Index(one.Type, "(")
		if n > 0 {
			one.Type = one.Type[0:n]
		}
		n = strings.Index(one.Type, " ")
		if n > 0 {
			one.Type = one.Type[0:n]
		}
		all = append(all, one)
		num++
	}
	return all

}

func Connect(driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)

	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(30)
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	return db, err
}

func TypeMToP(m string) string {
	if _, ok := typeArr[m]; ok {
		return typeArr[m]
	}
	return "string"
}
