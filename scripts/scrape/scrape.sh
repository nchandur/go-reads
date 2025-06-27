echo "starting scraping process"

for i in $(seq 21 23);
do 
    go run scripts/scrape/main.go https://www.goodreads.com/list/show/1.Best_Books_Ever?page=$i
    mongodump --host localhost --port 9001 --db books --out data/backup
done
