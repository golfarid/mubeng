# ./proxy_checker.sh url proxy_user proxy_password file_with_proxies
# Example: ./proxy_checker.sh "https://wbxcatalog-ru.wildberries.ru/nm-2-card/catalog?appType=64&dest=-1255563%2C-1278703%2C-102269%2C-1029256&emp=0&lang=ru&locale=ru&nm=14161490%3B21131553&pricemarginCoeff=1&pricemarginMax=0&pricemarginMin=0&reg=0&regions=83%2C75%2C64%2C4%2C38%2C30%2C33%2C70%2C71%2C22%2C31%2C66%2C68%2C40%2C48%2C1%2C69%2C80&spp=0&stores=117673%2C122258%2C122259%2C125238%2C125239%2C125240%2C6159%2C507%2C3158%2C117501%2C120602%2C120762%2C6158%2C121709%2C124731%2C159402%2C2737%2C130744%2C117986%2C1733%2C686%2C132043%2C161812%2C1193&version=3" 
# "oAWi6CRQ0xJz" "mayak" "Proxy_oAWi6CRQ0xJz_HTTP.txt" >> "checked.txt"


url=$1
proxy_user=$2
proxy_password=$3
files=()

counter=0
for var in "$@"
do
    ((counter++))
    if [ $counter -gt 3 ]; then
        files+=("$var")
    fi
done

for path in "${files[@]}"
do
    while read line; do
      cline=${line//[$'\t\r\n ']}
      checking_proxy="$(printf 'http://%s:%s@%s' $proxy_user $proxy_password $cline)"

      STATUS_CODE=$(curl \
          --output /dev/null \
          --silent \
          --proxy $checking_proxy \
          --write-out "%{http_code}" \
          "$url")

      if (( STATUS_CODE == 200 ))
      then
        echo "$checking_proxy"
      fi
    done < "$path"
done