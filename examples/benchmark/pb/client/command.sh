#!/bin/sh

net="tcp"
host="127.0.0.1"
pbar="false"

./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=false -noresponse=false -clients=1 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=false -noresponse=false -clients=12 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=false -noresponse=false -clients=48 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=false -noresponse=false -clients=192 -bar=$pbar

./client -network=$net -h=$host -total=1000000 -concurrent=true -batch=false -noresponse=false -clients=1 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=true -batch=false -noresponse=false -clients=2 -bar=$pbar

./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=true -noresponse=false -clients=1 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=true -noresponse=false -clients=2 -bar=$pbar

./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=false -noresponse=true -clients=1 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=false -noresponse=true -clients=2 -bar=$pbar

./client -network=$net -h=$host -total=1000000 -concurrent=true -batch=true -noresponse=false -clients=1 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=true -batch=true -noresponse=false -clients=2 -bar=$pbar

./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=true -noresponse=true -clients=1 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=false -batch=true -noresponse=true -clients=2 -bar=$pbar

./client -network=$net -h=$host -total=1000000 -concurrent=true -batch=false -noresponse=true -clients=1 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=true -batch=false -noresponse=true -clients=2 -bar=$pbar

./client -network=$net -h=$host -total=1000000 -concurrent=true -batch=true -noresponse=true -clients=1 -bar=$pbar
./client -network=$net -h=$host -total=1000000 -concurrent=true -batch=true -noresponse=true -clients=2 -bar=$pbar

