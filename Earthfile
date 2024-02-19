VERSION 0.6

build-all:
    BUILD ./functions/oracle+compile
    BUILD ./functions/result+compile
    BUILD ./functions/scheduler+compile
    BUILD ./functions/image+compile
    BUILD ./functions/status+compile

test:
    BUILD ./functions/oracle+test
    BUILD ./functions/result+test
    BUILD ./functions/scheduler+test
    # BUILD ./functions/image+test
    # BUILD ./functions/status+test
    BUILD ./functions/lib+test
    LOCALLY
    RUN cat /dev/null > report.json
    RUN cat \
    ./functions/oracle/reports/report.json \
    ./functions/result/reports/report.json \
    ./functions/scheduler/reports/report.json \
    ./functions/lib/reports/report.json \
    > report.json
    # RUN mkdir coverage && cp \
    # ./functions/oracle/reports/coverage.out \
    # ./functions/result/reports/coverage.out \
    # ./functions/scheduler/reports/coverage.out \
    # ./functions/lib/reports/coverage.out \
    #  ./coverage/

sonar:
    FROM sonarsource/sonar-scanner-cli:latest
    BUILD +test
    COPY --dir . /app
    RUN --secret SONAR_TOKEN \
    sonar-scanner \
    -Dsonar.projectKey=pixelart \
    -Dsonar.sources=/app \
    -Dsonar.host.url=http://localhost:9000
    
sonar-local:
    LOCALLY
    BUILD +test
    RUN sonar-scanner 

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
    
debug:
    FROM sonarsource/sonar-scanner-cli:latest
    RUN curl http://localhost:9000