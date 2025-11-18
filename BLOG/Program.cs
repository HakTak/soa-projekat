var builder = WebApplication.CreateBuilder(args);

// Dodajemo servise za kontrolere i Swagger
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Ucitavamo konfiguraciju iz appsettings.json (sekcija "MongoDatabaseSettings")
builder.Services.Configure<BLOG.Database.MongoDatabaseSettings>(
    builder.Configuration.GetSection("MongoDatabaseSettings")
);

// Registrujemo MongoDatabaseSettings kao singleton, da se može koristiti u repository-ju
builder.Services.AddSingleton(resolver =>
    resolver.GetRequiredService<Microsoft.Extensions.Options.IOptions<BLOG.Database.MongoDatabaseSettings>>().Value
);
// Registracija CommentRepository i CommentService
builder.Services.AddScoped<BLOG.Repositories.ICommentRepository, BLOG.Repositories.CommentRepository>();
builder.Services.AddScoped<BLOG.Services.CommentService>();



var app = builder.Build();

// Swagger samo u development
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.MapControllers();

app.Run();
