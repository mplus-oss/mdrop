using System.Text.Json;
using System.Text.Json.Serialization;

namespace MDrop.Broker.Controllers.Middlewares;
public class VerifyPrivateModeTokenPipeline
{
    public void Configure(IApplicationBuilder app)
    {
        app.Use(async (context, next) => 
        {
            var response = new VerifyReturnJson();

            var privToken = Constant.PrivateModeToken;
            if (privToken == "") 
            {
                await next();
                return;
            }

            var token = context.Request.Headers.Authorization;
            var tokenArr = token.ToString().Split(' ');
            if (tokenArr.Length == 1)
            {
                context.Response.ContentType = "application/json";
                context.Response.StatusCode = 401;
                response.Message = "Authorization needed.";
                await context.Response.WriteAsync(JsonSerializer.Serialize(response));
                return;
            }

            token = tokenArr[1];
            if (token != privToken)
            {
                context.Response.ContentType = "application/json";
                context.Response.StatusCode = 401;
                response.Message = "Token is invalid.";
                await context.Response.WriteAsync(JsonSerializer.Serialize(response));
                return;
            }

            await next();
        });
    }

    private class VerifyReturnJson
    {
        [JsonPropertyName("message")]
        public string Message { get; set; } = "";
    }
}
