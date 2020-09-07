FROM mcr.microsoft.com/dotnet/core/sdk:3.1-bionic as dump-install

WORKDIR /app
ENV PATH="$PATH:/root/.dotnet/tools"

RUN dotnet tool install -g dotnet-gcdump

COPY ./dumper ./

