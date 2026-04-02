package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Template struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Subject string `json:"subject"`
	Body string `json:"body"`
	Variables string `json:"variables"`
	Status string `json:"status"`
	UsageCount int `json:"usage_count"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"brandingiron.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS templates(id TEXT PRIMARY KEY,name TEXT NOT NULL,type TEXT DEFAULT 'email',subject TEXT DEFAULT '',body TEXT DEFAULT '',variables TEXT DEFAULT '[]',status TEXT DEFAULT 'active',usage_count INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Template)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO templates(id,name,type,subject,body,variables,status,usage_count,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Type,e.Subject,e.Body,e.Variables,e.Status,e.UsageCount,e.CreatedAt);return err}
func(d *DB)Get(id string)*Template{var e Template;if d.db.QueryRow(`SELECT id,name,type,subject,body,variables,status,usage_count,created_at FROM templates WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Type,&e.Subject,&e.Body,&e.Variables,&e.Status,&e.UsageCount,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Template{rows,_:=d.db.Query(`SELECT id,name,type,subject,body,variables,status,usage_count,created_at FROM templates ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Template;for rows.Next(){var e Template;rows.Scan(&e.ID,&e.Name,&e.Type,&e.Subject,&e.Body,&e.Variables,&e.Status,&e.UsageCount,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Template)error{_,err:=d.db.Exec(`UPDATE templates SET name=?,type=?,subject=?,body=?,variables=?,status=?,usage_count=? WHERE id=?`,e.Name,e.Type,e.Subject,e.Body,e.Variables,e.Status,e.UsageCount,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM templates WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM templates`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Template{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ? OR subject LIKE ? OR body LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["type"];ok&&v!=""{where+=" AND type=?";args=append(args,v)}
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,type,subject,body,variables,status,usage_count,created_at FROM templates WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Template;for rows.Next(){var e Template;rows.Scan(&e.ID,&e.Name,&e.Type,&e.Subject,&e.Body,&e.Variables,&e.Status,&e.UsageCount,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM templates GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
