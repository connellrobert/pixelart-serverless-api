VERSION 0.6

build-all:
    BUILD ./functions/oracle+compile
    BUILD ./functions/result+compile
    BUILD ./functions/scheduler+compile
    BUILD ./functions/image+compile
    BUILD ./functions/status+compile

deploy:
    BUILD +build-all
    BUILD ./infrastructure+apply

destroy:
    BUILD +build-all
    BUILD ./infrastructure+destroy

update-deps:
    BUILD ./functions/lib+update
    BUILD ./functions/oracle+update
    BUILD ./functions/result+update
    BUILD ./functions/scheduler+update
    BUILD ./functions/image+update
    BUILD ./functions/status+update
    