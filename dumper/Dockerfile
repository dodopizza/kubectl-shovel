FROM mcr.microsoft.com/dotnet/sdk:9.0-noble AS tools-install

RUN dotnet tool install -g dotnet-gcdump --version "9.*" && \
    dotnet tool install -g dotnet-trace --version "9.*" && \
    dotnet tool install -g dotnet-dump --version "9.*"

FROM mcr.microsoft.com/dotnet/runtime:9.0-noble

ARG DOTNET_TOOLS_PATH="/root/.dotnet/tools"
ARG DOTNET_RUNTIME_PATH="/usr/share/dotnet/shared/Microsoft.NETCore.App/${DOTNET_VERSION}"
ENV PATH="${PATH}:${DOTNET_TOOLS_PATH}:${DOTNET_RUNTIME_PATH}"

WORKDIR /app
COPY --from=tools-install ${DOTNET_TOOLS_PATH} ${DOTNET_TOOLS_PATH}
COPY ./bin/dumper ./

ENTRYPOINT [ "/app/dumper" ]
