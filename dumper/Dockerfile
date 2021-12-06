ARG RUNTIME_VERSION=3.1.20

FROM mcr.microsoft.com/dotnet/core/sdk:3.1-bionic as tools-install

RUN dotnet tool install -g dotnet-gcdump && \
    dotnet tool install -g dotnet-trace && \
    dotnet tool install -g dotnet-dump

FROM mcr.microsoft.com/dotnet/core/runtime:${RUNTIME_VERSION}-bionic

ARG RUNTIME_VERSION=3.1.20
ARG DOTNET_TOOLS_PATH="/root/.dotnet/tools"
ARG DOTNET_INTERNAL_TOOLS_PATH="/usr/share/dotnet/shared/Microsoft.NETCore.App/${RUNTIME_VERSION}"
ENV PATH="${PATH}:${DOTNET_TOOLS_PATH}:${DOTNET_INTERNAL_TOOLS_PATH}"
WORKDIR /app

COPY --from=tools-install ${DOTNET_TOOLS_PATH} ${DOTNET_TOOLS_PATH}

COPY ./bin/dumper ./

ENTRYPOINT [ "/app/dumper" ]
