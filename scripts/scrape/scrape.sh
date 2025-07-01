echo "starting scraping process"

for i in $(seq 11 11);
do 
    go run scripts/scrape/main.go https://www.goodreads.com/list/show/264.Books_That_Everyone_Should_Read_At_Least_Once?page=$i
    mongodump --host localhost --port 9001 --db books --out data/backup_$(date +%Y%m%d%H%M%S)
done