using System.Security.Cryptography.X509Certificates;

namespace MDrop.Broker;
public static class Constant
{
    public static readonly string PrivateModeToken = Environment.GetEnvironmentVariable("PRIVATE_MODE_TOKEN") ?? "";
    public static X509Certificate2 Certificate = X509Certificate2.CreateFromPemFile("cert.pem", "prikey.pem");
}
