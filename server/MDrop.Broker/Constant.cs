using System.Security.Cryptography.X509Certificates;

namespace MDrop.Broker;
public static class Constant
{
    public static readonly string TunnelHost = Environment.GetEnvironmentVariable("TUNNEL_HOST") ?? "127.0.0.1";
    public static readonly int TunnelPort = int.Parse(Environment.GetEnvironmentVariable("TUNNEL_PORT") ?? "2222");
    public static readonly string PrivateModeToken = Environment.GetEnvironmentVariable("PRIVATE_MODE_TOKEN") ?? "";
    public static X509Certificate2 Certificate = X509Certificate2.CreateFromEncryptedPemFile(
        "cert.pem",
        Environment.GetEnvironmentVariable("X509_CERTIFICATE_PASSWORD") ?? "",
        "prikey.pem");
}
