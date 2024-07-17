using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.ModelBinding;

namespace MDrop.Broker.Controllers;
[ApiController, Route("/verify")]
public class VerifyController : Controller
{
    public ActionResult<VerifyReturnJson> VerifyClientOnPrivateMode(
        [FromQuery, BindRequired] string token)
    {
        var response = new VerifyReturnJson();
        if (Constant.PrivateModeToken == "")
        {
            response.Message = "This broker is public mode. Refusing.";
            return BadRequest(response);
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
    }
}
