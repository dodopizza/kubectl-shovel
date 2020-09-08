FROM mcr.microsoft.com/dotnet/core/sdk:3.1-bionic as tools-install

RUN dotnet tool install -g dotnet-gcdump


FROM mcr.microsoft.com/dotnet/core/runtime:3.1

ARG DOTNET_TOOLS_PATH="/root/.dotnet/tools"
ENV PATH="${PATH}:${DOTNET_TOOLS_PATH}"
WORKDIR /app

COPY --from=tools-install ${DOTNET_TOOLS_PATH} ${DOTNET_TOOLS_PATH}

COPY ./dumper ./

ENTRYPOINT [ "/app/dumper" ]
