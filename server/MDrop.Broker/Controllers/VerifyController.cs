using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Mvc;

namespace MDrop.Broker.Controllers;
[ApiController]
[Route("/verify")]
public class VerifyController : Controller
{
    [HttpPost]
    public ActionResult<VerifyReturnJson> VerifyClientOnPrivateMode([FromQuery] string token)
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
