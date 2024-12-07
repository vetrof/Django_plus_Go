FROM python:3.10-alpine

WORKDIR /code
COPY ./requirements.txt .
RUN pip install --upgrade pip
RUN pip install -r requirements.txt 
RUN pip install psycopg2-binary==2.9.10


COPY . .

