using System.Text.Json;
using Microsoft.AspNetCore.Mvc;

namespace MDrop.Broker.Functions;
public static class CustomContentResult
{
    public static ContentResult JsonCustomResult<T>(int statusCode, T content) =>
        CustomContentResult.GenerateCustomResult(
            statusCode,
            "application/json",
            JsonSerializer.Serialize<T>(content));

    public static ContentResult GenerateCustomResult(
        int statusCode,
        string contentType,
        string content)
    {
        return new ContentResult()
        {
            StatusCode = statusCode,
            ContentType = contentType,
            Content = content
        };
    }
}
