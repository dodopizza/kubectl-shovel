FROM mcr.microsoft.com/dotnet/sdk:6.0-focal as tools-install

RUN dotnet tool install -g dotnet-gcdump && \
    dotnet tool install -g dotnet-trace && \
    dotnet tool install -g dotnet-dump

FROM mcr.microsoft.com/dotnet/runtime:6.0.3-focal

ARG DOTNET_TOOLS_PATH="/root/.dotnet/tools"
ARG DOTNET_RUNTIME_PATH="/usr/share/dotnet/shared/Microsoft.NETCore.App/6.0.3"
ENV PATH="${PATH}:${DOTNET_TOOLS_PATH}:${DOTNET_RUNTIME_PATH}"

WORKDIR /app
COPY --from=tools-install ${DOTNET_TOOLS_PATH} ${DOTNET_TOOLS_PATH}
COPY ./bin/dumper ./

ENTRYPOINT [ "/app/dumper" ]
