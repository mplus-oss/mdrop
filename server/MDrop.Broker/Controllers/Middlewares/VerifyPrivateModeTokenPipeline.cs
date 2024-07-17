namespace MDrop.Broker.Controllers.Middlewares;
public class VerifyPrivateModeTokenPipeline
{
    public void Configure(IApplicationBuilder app)
    {
        app.Use(async (context, next) => 
        {
            var privToken = Constant.PrivateModeToken;
            if (privToken == "") 
            {
                await next();
                return;
            }

            var token = context.Request.Headers.Authorization;

            await next();
        });
    }
}
