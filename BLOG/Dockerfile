# Stage 1: Build
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

# 1. Copy csproj
COPY ["BLOG/BLOG.csproj", "BLOG/"]

# 2. Restore (Downloads Nuget packages)
RUN dotnet restore "BLOG/BLOG.csproj"

# 3. Copy Source Code
COPY ["BLOG/", "BLOG/"]

# 4. Copy COMMON (So the ../COMMON path in csproj works)
COPY ["COMMON/", "COMMON/"]

# 5. Build
WORKDIR "/src/BLOG"
RUN dotnet publish "BLOG.csproj" -c Release -o /app/publish /p:UseAppHost=false

# Stage 2: Runtime
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS final
WORKDIR /app
EXPOSE 5251
EXPOSE 9091
COPY --from=build /app/publish .
ENTRYPOINT ["dotnet", "BLOG.dll"]