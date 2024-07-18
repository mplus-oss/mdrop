using MDrop.Broker.Functions;
using System.Text.Json.Serialization;

namespace MDrop.Broker.Controllers.Middlewares;
public class VerifyPrivateModeTokenPipeline
{
    public void Configure(IApplicationBuilder app)
    {
        app.Use(async (context, next) => 
        {
            var response = new VerifyReturnJson();

            if (Constant.PrivateModeToken == "") 
            {
                await next();
                return;
            }

            var token = context.Request.Headers.Authorization;
            var tokenArr = token.ToString().Split(' ');
            if (tokenArr.Length == 1)
            {
                response.Message = "Authorization needed.";
                await CustomContentResult.JsonCustomContextResult(context, 401, response);
                return;
            }

            if (tokenArr[0].ToLower() != "bearer")
            {
                response.Message = "Authorization invalid.";
                await CustomContentResult.JsonCustomContextResult(context, 401, response);
                return;
            }
            if (tokenArr[1] != Constant.PrivateModeToken)
            {
                response.Message = "Token is invalid.";
                await CustomContentResult.JsonCustomContextResult(context, 401, response);
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
