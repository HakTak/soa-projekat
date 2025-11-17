# Stage 1: build
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY ["BLOG.csproj", "./"]
RUN dotnet restore "BLOG.csproj"
COPY . .
RUN dotnet publish "BLOG.csproj" -c Release -o /app/publish

# Stage 2: runtime
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app
COPY --from=build /app/publish .

# Postavljamo port koji aplikacija sluša unutar kontejnera
ENV ASPNETCORE_URLS=http://+:5251

EXPOSE 5251
ENTRYPOINT ["dotnet", "BLOG.dll"]
