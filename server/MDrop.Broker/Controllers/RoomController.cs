using MDrop.Broker.Functions;
using MDrop.Broker.DataStructures;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.ModelBinding;
using JWT.Algorithms;
using JWT.Builder;

namespace MDrop.Broker.Controllers;
[ApiController]
[Route("/room")]
public class RoomController : Controller
{
    [HttpPost("create")]
    public ActionResult CreateRoom(
        [FromQuery, BindRequired, Range(1_024, 59_999)] int port,
        [FromQuery] int durationInHours = 3)
    {
        var response = new RoomReturnJson();

        try
        {
            var expTime = DateTimeOffset.UtcNow.AddHours(durationInHours);
            var jwt = JwtBuilder.Create()
                .WithAlgorithm(new RS2048Algorithm(Constant.Certificate))
                .AddClaim(JWTClaim.Issuer, "MDrop.Broker")
                .AddClaim(JWTClaim.Audience, "MDrop.Client")
                .AddClaim(JWTClaim.ExpirationTime, expTime.ToUnixTimeSeconds())
                .AddClaim(JWTClaim.BodyValue, "limau");

            response.Message = "Ok!";
            response.Token = jwt.Encode();
            return Ok(response);
        }
        catch (Exception e)
        {
            response.Message = e.Message;
            return CustomContentResult.JsonCustomResult(500, response);
        }
    }

    public class RoomReturnJson
    {
        [JsonPropertyName("message")]
        public string Message { get; set; } = "";

        [JsonPropertyName("token")]
        [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
        public string? Token { get; set; }
    }
}
