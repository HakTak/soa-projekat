
using BLOG.GrpcServices;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using Microsoft.AspNetCore.Builder;
var builder = WebApplication.CreateBuilder(args);

// Dodajemo servise za kontrolere i Swagger
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

builder.Services.AddGrpc();

// --- POCETAK INTEGRACIJE JAEGERA ---

// Preuzimamo ime servisa iz environment varijable (definisane u docker-compose)
// Ako nije definisano, koristimo default "blog-service"
var serviceName = Environment.GetEnvironmentVariable("OTEL_SERVICE_NAME") ?? "blog-service";

builder.Services.AddOpenTelemetry()
    .WithTracing(tracerProviderBuilder =>
    {
        tracerProviderBuilder
            .AddSource(serviceName)
            .SetResourceBuilder(
                ResourceBuilder.CreateDefault()
                    .AddService(serviceName: serviceName, serviceVersion: "1.0.0"))
            
            // Prati sve dolazne HTTP i gRPC zahteve ka tvom servisu
            .AddAspNetCoreInstrumentation()
            
            // Prati sve odlazne HTTP zahteve (ako tvoj blog servis zove nekog drugog)
            .AddHttpClientInstrumentation()
            
            // Šalje podatke Jaegeru koristeći OTLP protokol
            // Čita adresu automatski iz "OTEL_EXPORTER_OTLP_ENDPOINT" (iz docker-compose)
            .AddOtlpExporter();
    });

// --- KRAJ INTEGRACIJE JAEGERA ---


// Ucitavamo konfiguraciju iz appsettings.json (sekcija "MongoDatabaseSettings")
builder.Services.Configure<BLOG.Database.MongoDatabaseSettings>(
    builder.Configuration.GetSection("MongoDatabaseSettings")
);

// Registrujemo MongoDatabaseSettings kao singleton
builder.Services.AddSingleton(resolver =>
    resolver.GetRequiredService<Microsoft.Extensions.Options.IOptions<BLOG.Database.MongoDatabaseSettings>>().Value
);

// Registracija Repository i Service slojeva
builder.Services.AddScoped<BLOG.Repositories.ICommentRepository, BLOG.Repositories.CommentRepository>();
builder.Services.AddScoped<BLOG.Services.CommentService>();
builder.Services.AddScoped<BLOG.Repositories.IBlogRepository, BLOG.Repositories.BlogRepository>();
builder.Services.AddScoped<BLOG.Services.BlogService>();


var app = builder.Build();

// Swagger samo u development
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}
app.UseStaticFiles();
app.UseHttpsRedirection();

app.MapControllers();

app.MapGrpcService<BlogGrpcService>();

app.Run();