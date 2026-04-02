package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Release struct{
	ID string `json:"id"`
	Version string `json:"version"`
	Title string `json:"title"`
	Body string `json:"body"`
	Status string `json:"status"`
	PublishedAt string `json:"published_at"`
	Tags string `json:"tags"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"brandingiron.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS releases(id TEXT PRIMARY KEY,version TEXT NOT NULL,title TEXT DEFAULT '',body TEXT DEFAULT '',status TEXT DEFAULT 'draft',published_at TEXT DEFAULT '',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Release)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO releases(id,version,title,body,status,published_at,tags,created_at)VALUES(?,?,?,?,?,?,?,?)`,e.ID,e.Version,e.Title,e.Body,e.Status,e.PublishedAt,e.Tags,e.CreatedAt);return err}
func(d *DB)Get(id string)*Release{var e Release;if d.db.QueryRow(`SELECT id,version,title,body,status,published_at,tags,created_at FROM releases WHERE id=?`,id).Scan(&e.ID,&e.Version,&e.Title,&e.Body,&e.Status,&e.PublishedAt,&e.Tags,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Release{rows,_:=d.db.Query(`SELECT id,version,title,body,status,published_at,tags,created_at FROM releases ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Release;for rows.Next(){var e Release;rows.Scan(&e.ID,&e.Version,&e.Title,&e.Body,&e.Status,&e.PublishedAt,&e.Tags,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM releases WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM releases`).Scan(&n);return n}
