FROM centos

RUN yum install git curl -y

COPY target /app

WORKDIR /app

RUN touch /etc/localtime /etc/timezone && \
chown -R 1001 /app /etc/localtime /etc/timezone  && \
chgrp -R 0 /app /etc/localtime /etc/timezone  && \
chmod -R g=u /app /etc/localtime /etc/timezone

USER 1001:0

EXPOSE 9000

CMD /app/run.sh