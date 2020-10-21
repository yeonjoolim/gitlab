package mysql

import (
 "time"
//"os"
"database/sql"
"log"
_ "github.com/go-sql-driver/mysql"
)
func main(log string) {


//create("mec_sju")
   
Insert(log)


}
func create(name string) {

   db, err := sql.Open("mysql", "root:password@tcp(129.254.170.5:30155)/")
   if err != nil {
 	log.Fatal(err)
   }
   defer db.Close()

   _,err = db.Exec("CREATE DATABASE "+name)
   if err != nil {
 	log.Fatal(err)
   }

   _,err = db.Exec("USE "+name)
   if err != nil {
 	log.Fatal(err)
   }

   _,err = db.Exec("CREATE TABLE mec_sju.grafana ( time timestamp, log text )")
   if err != nil {
 	log.Fatal(err)
   }
}

func Insert(str string){

   db, err := sql.Open("mysql", "root:password@tcp(129.254.170.5:30155)/mec_sju")
    if err != nil {
 	log.Fatal(err)
   }

   t := time.Now()
   ts := t.Format("2006-01-02 15:04:05")   
   result, err := db.Exec( "insert into mec_sju.grafana(time, log) values('"+ts+"','"+str+"')" )
     if err != nil {
 	log.Fatal(err)
   }
   _, err = result.RowsAffected()
   db.Close()


}

