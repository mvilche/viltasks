FROM mvilche/viltasks

USER root

RUN yum install git curl golang -y 

WORKDIR /opt/go

ENV GOPATH=/opt/go

RUN go get -v github.com/revel/revel && \
    go get -v github.com/revel/cmd/revel 

COPY . /opt/go/src/viltasks

ENV PATH=$PATH:$GOPATH/bin

RUN revel build /opt/go/src/viltasks /app

WORKDIR /app

RUN mkdir -p /app/database && ln -s src/viltasks/conf conf && touch /etc/localtime /etc/timezone && \
chown -R 1001 /app /etc/localtime /etc/timezone  && \
chgrp -R 0 /app /etc/localtime /etc/timezone  && \
chmod -R g=u /app /etc/localtime /etc/timezone

USER 1001:0

EXPOSE 9000

CMD /app/run.sh