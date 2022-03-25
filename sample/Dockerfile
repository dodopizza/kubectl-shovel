# syntax=docker/dockerfile:1
ARG SDK_IMAGE_TAG
ARG RUNTIME_IMAGE_TAG
FROM mcr.microsoft.com/dotnet/sdk:${SDK_IMAGE_TAG} AS build
ARG FRAMEWORK
WORKDIR /app
COPY . .
RUN dotnet publish \
    --framework ${FRAMEWORK} \
    --configuration Release \
    --output ./output

FROM mcr.microsoft.com/dotnet/aspnet:${RUNTIME_IMAGE_TAG}
WORKDIR /app
COPY --from=build /app/output .
ENTRYPOINT ["dotnet", "sample.dll"]