namespace MDrop.Broker;
public class Program
{
    static void Main(string[] args) => new Program().Start(args);
    public void Start(string[] args)
    {
        var builder = WebApplication.CreateBuilder(args);

        // Add services to the container
        var app = ConfigureServiceBuilder(ref builder);

        // Configure the HTTP request pipeline
        if (!app.Environment.IsDevelopment()) ConfigureDevelopmentRequestPipeline(ref app);

        // Set middleware
        ConfigureMiddleware(ref app);

        // Run app
        app.Run();
    }

    /**
     * Configure service for app builder.
     * This method can be used for insert singleton or service that can be injected on controller.
     */
    public WebApplication ConfigureServiceBuilder(ref WebApplicationBuilder builder)
    {
        // Insert singleton if needed
        return builder.Build();
    }

    /**
     * Configure middleware for webapp.
     */
    public static void ConfigureMiddleware(ref WebApplication app)
    {
        app.UseRouting();
        //app.MapControllers();

        // Health Check
        app.MapGet("/", async context =>
        {
            await context.Response.WriteAsJsonAsync(new { Message = "Hello world!" });
        });
    }

    /**
     * Configure exception request pipeline.
     */
    public void ConfigureDevelopmentRequestPipeline(ref WebApplication app)
    {
        app.UseExceptionHandler("/error");
    }
}
