FROM mcr.microsoft.com/dotnet/sdk:6.0-focal as tools-install

RUN dotnet tool install -g dotnet-gcdump && \
    dotnet tool install -g dotnet-trace && \
    dotnet tool install -g dotnet-dump

FROM mcr.microsoft.com/dotnet/runtime:6.0.2-focal

ARG DOTNET_TOOLS_PATH="/root/.dotnet/tools"
ENV PATH="${PATH}:${DOTNET_TOOLS_PATH}"

WORKDIR /app
COPY --from=tools-install ${DOTNET_TOOLS_PATH} ${DOTNET_TOOLS_PATH}
COPY ./bin/dumper ./

ENTRYPOINT [ "/app/dumper" ]
