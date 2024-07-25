using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.ModelBinding;

namespace MDrop.Broker.Controllers;
[ApiController, Route("/verify")]
public class VerifyController : Controller
{
    [HttpPost]
    public ActionResult<VerifyReturnJson> VerifyClientOnPrivateMode(
        [FromQuery, BindRequired] string token)
    {
        var response = new VerifyReturnJson();
        if (Constant.PrivateModeToken == "")
        {
            response.Message = "This broker is public mode. Refusing.";
            response.IsPublic = true;
            return Ok(response);
        }

        if (Constant.PrivateModeToken != token)
        {
            response.Message = "Wrong token.";
            return Unauthorized(response);
        }

        response.Message = "Authenticated";
        return Ok(response);
    }

    public class VerifyReturnJson
    {
        [JsonPropertyName("message")]
        public string Message { get; set; } = "";

        [JsonPropertyName("isPublic")]
        public bool IsPublic { get; set; } = false;

        [JsonPropertyName("tunnel")]
        public TunnelProperty Tunnel { get; set; } = new TunnelProperty();
        public class TunnelProperty
        {
            [JsonPropertyName("host")]
            public string Host { get; set; } = Constant.TunnelHost;

            [JsonPropertyName("port")]
            public int Port { get; set; } = Constant.TunnelPort;

            [JsonPropertyName("proxy")]
            public string Proxy { get; set; } = Constant.TunnelProxy;
        }
    }
}
