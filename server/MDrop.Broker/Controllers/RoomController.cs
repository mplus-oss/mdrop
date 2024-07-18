using MDrop.Broker.Controllers.Middlewares;
using MDrop.Broker.Functions;
using MDrop.Broker.DataStructures;
using System.ComponentModel.DataAnnotations;
using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.ModelBinding;
using JWT.Algorithms;
using JWT.Builder;
using JWT.Exceptions;

namespace MDrop.Broker.Controllers;
[ApiController, Route("/room"), MiddlewareFilter<VerifyPrivateModeTokenPipeline>]
public class RoomController : Controller
{
    [HttpPost("create")]
    public ActionResult CreateRoom(
        [FromQuery, BindRequired, Range(10_000, 59_999)] int port,
        [FromQuery] int durationInHours = 3)
    {
        var response = new RoomReturnJson();

        try
        {
            var expTime = DateTimeOffset.UtcNow.AddHours(durationInHours);
            var portJson = JsonSerializer.Serialize(new RoomReturnJWT { Port = port });
            Console.WriteLine(portJson);
            var jwt = JwtBuilder.Create()
                .WithAlgorithm(new RS2048Algorithm(Constant.Certificate))
                .AddClaim(JWTClaim.Issuer, "MDrop.Broker")
                .AddClaim(JWTClaim.Audience, "MDrop.Client")
                .AddClaim(JWTClaim.ExpirationTime, expTime.ToUnixTimeSeconds())
                .AddClaim(JWTClaim.Value, portJson);

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

    [HttpPost("join")]
    public ActionResult JoinRoom(
        [FromQuery, BindRequired] string token)
    {
        var response = new RoomReturnJson();

        try
        {
            var data = JwtBuilder.Create()
                .WithAlgorithm(new RS2048Algorithm(Constant.Certificate))
                .MustVerifySignature()
                .Decode<IDictionary<string, object>>(token);
            if (data == null) throw new Exception("Value is empty!");
            if (data[JWTClaim.Value] == null) throw new Exception("Value is empty!");

            var portJson = JsonSerializer.Deserialize<RoomReturnJWT>(data[JWTClaim.Value].ToString()!);
            if (portJson == null) throw new Exception("Invalid JSON");

            response.Message = "Ok!";
            response.Port = portJson.Port;
            return Ok(response);
        }
        catch (TokenExpiredException)
        {
            response.Message = "Your token is expired.";
            return Unauthorized(response);
        }
        catch (SignatureVerificationException)
        {
            response.Message = "Your token has invalid signature.";
            return Unauthorized(response);
        }
        catch (Exception e)
        {
            response.Message = $"[{e.GetType()}] {e.Message}";
            return CustomContentResult.JsonCustomResult(500, response);
        }
    }

    private class RoomReturnJson
    {
        [JsonPropertyName("message")]
        public string Message { get; set; } = "";

        [JsonPropertyName("token")]
        [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
        public string? Token { get; set; }

        [JsonPropertyName("port")]
        [JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
        public int? Port { get; set; }
    }

    private class RoomReturnJWT
    {
        [JsonPropertyName("port")]
        public int Port { get; set; }
    }
}
