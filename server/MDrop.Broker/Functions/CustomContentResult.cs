using System.Text.Json;
using Microsoft.AspNetCore.Mvc;

namespace MDrop.Broker.Functions;
public static class CustomContentResult
{
    public static async Task JsonCustomContextResult<T>(HttpContext context, int statusCode, T content) =>
        await CustomContentResult.GenerateCustomContextResult(
            context,
            statusCode,
            "application/json",
            JsonSerializer.Serialize<T>(content));

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

    public static async Task GenerateCustomContextResult(
        HttpContext context,
        int statusCode,
        string contentType,
        string content)
    {
        context.Response.ContentType = "application/json";
        context.Response.StatusCode = 401;
        await context.Response.WriteAsync(content);
    }
}

