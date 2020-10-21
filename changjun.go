package main
import (
"fmt"
"os"
"database/sql"
_ "github.com/go-sql-driver/mysql"
)
func main() {
// conn, err := sql.Open("mysql", "계정명:계정패스워드@tcp(DB주소:DB포트)/데이터베이스명")
conn, err := sql.Open("mysql", "root:password@tcp(129.254.170.5:30155)/test")
if err != nil {
fmt.Println(err)
os.Exit(1)
}
// DB에 저장할 데이터 삽입
result, err := conn.Exec( "insert into test.time_test(date, score) values('2020-10-20',
'400')")
if err != nil {
fmt.Println(err)
os.Exit(1)
}
// RowsAffected()를 통해 insert한 데이터 갯수를 확인
nRow, err := result.RowsAffected()
fmt.Println("Insert count: ", nRow)
conn.Close()
}
