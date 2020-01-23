FROM aliyunfc/runtime-python3.6:1.7.7

EXPOSE 5000

WORKDIR /code
ADD . /code

RUN pip install -r requirements.txt --target=./

ENTRYPOINT [ "" ]

CMD python app.py
