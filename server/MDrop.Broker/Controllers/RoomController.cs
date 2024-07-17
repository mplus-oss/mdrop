using MDrop.Broker.Functions;
using MDrop.Broker.DataStructures;
using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Mvc;
using JWT.Algorithms;
using JWT.Builder;

namespace MDrop.Broker.Controllers;
[ApiController]
[Route("/room")]
public class RoomController : Controller
{
    [HttpPost("create")]
    public ActionResult CreateRoom(
        [FromQuery] int port,
        [FromQuery] int durationInHours = 3)
    {
        var response = new RoomReturnJson();

        if (port < 1024 && port > 59_999)
        {
            response.Message = "Port must be on range 1024-59999";
            return BadRequest(response);
        }

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
            return CustomContentResult.JsonCustomResult(432, response);
        }
    }

    public class RoomReturnJson
    {
        [JsonPropertyName("message")]
        public string Message { get; set; } = "";

        [JsonPropertyName("token")]
        public string? Token { get; set; }
    }
}
