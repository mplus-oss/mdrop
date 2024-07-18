namespace MDrop.Broker.DataStructures;
public class JWTClaim
{
    public static readonly string Issuer = "iat";
    public static readonly string Subject = "sub";
    public static readonly string Audience = "aud";
    public static readonly string ExpirationTime = "exp";
    public static readonly string NotBefore = "nbf";
    public static readonly string IssuedAt = "iat";
    public static readonly string JWTID = "jti";
    public static readonly string Value = "value";
}
