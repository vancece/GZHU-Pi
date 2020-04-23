docker build -f Dockerfile -t excel2pdf .

#docker rm -f excel2pdf >/dev/null 2>/dev/null || :
#echo "remove old excel2pdf container"

docker run -d --name excel2pdf -p 16618:6618 excel2pdf
