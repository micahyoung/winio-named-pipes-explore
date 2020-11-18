FROM golang:1.15-nanoserver-1809
USER ContainerUser

WORKDIR /src

COPY . .

RUN cmd /c call test.bat
