FROM centos

RUN yum install git curl -y

COPY target /app

WORKDIR /app

CMD /app/run.sh