FROM centos:centos6

RUN rpm --import http://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-6
RUN yum install -y http://dl.fedoraproject.org/pub/epel/6/x86_64/epel-release-6-8.noarch.rpm
RUN yum install -y golang git gcc

RUN mkdir /src
WORKDIR /src
RUN git clone https://github.com/maebashi/docker-metricsd.git
WORKDIR /src/docker-metricsd
ENV GOPATH /
RUN go get -d; go build docker-metricsd.go
RUN cp docker-metricsd /
ADD installer /installer
CMD /installer

# sudo docker run -v /usr/local/bin:/target maebashi/docker-metricsd
# sudo /usr/local/bin/docker-metricsd
