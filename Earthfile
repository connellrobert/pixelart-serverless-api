VERSION 0.6

build-all:
    BUILD ./functions/oracle+compile
    BUILD ./functions/poll+compile
    BUILD ./functions/result+compile
    BUILD ./functions/scheduler+compile

deploy:
    BUILD +build-all
    BUILD ./infrastructure+plan

destroy:
    BUILD +build-all
    BUILD ./infrastructure+destroy