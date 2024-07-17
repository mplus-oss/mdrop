namespace MDrop.Broker;
public static class Constant
{
    public static readonly string PrivateModeToken = 
        Environment.GetEnvironmentVariable("PRIVATE_MODE_TOKEN") ?? "";
}
