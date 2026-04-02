package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-brandingiron/internal/server";"github.com/stockyard-dev/stockyard-brandingiron/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8840"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./brandingiron-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("brandingiron: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Brandingiron — Self-hosted release notes manager\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n\n",port,port,dataDir)
log.Printf("brandingiron: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
