package pkg

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/ttmars/goproxy"
)

var CaCert = []byte(`-----BEGIN CERTIFICATE-----
MIIC+jCCAeKgAwIBAgIRAODvM6qZEZDFdM1UkL5Be+owDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0yMjExMDgxMjM0MzVaFw0yMzExMDgxMjM0
MzVaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQCtSttU1i5aCN3cVfHQdxNhWV1ZiO0gAz/BQsi/OkpAgBLJU8v/zbpJ
7Ajlg7OzDcKWmK2Zjlo6492YlCtt9cQMrI05GusjjK/+NHenmuaAztUw6rnmxyEJ
Gq11eLNo4GC/O4q0lAmCVAhN4YxAsXG/m+EkCO7T0TpXdYUj/9IT1XQ9MJeWXt/i
fNP9YwZU9+ZEBa9KqNJV7ASmJqX0PgUI2LS3wXfu99YKzjg86GPoAv9nSvQcixPF
yYDPrxiCWGG936cbrga4VChn1ZQAL5GLpZkdAg1IOc2ufZVFuZuSa6KHUPL3a+V6
5d09vBQlCrU3kQ6L8hlFjx8jX9TAcNF5AgMBAAGjSzBJMA4GA1UdDwEB/wQEAwIF
oDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMBQGA1UdEQQNMAuC
CWxvY2FsaG9zdDANBgkqhkiG9w0BAQsFAAOCAQEATLg4+taOXL9Vn1P3bsNeZBru
mz01aagE/dQoOh2SLVUcouIQoB/SoSMPGtWaDmM53NgLypwAQWLbBkftAhcR6KAq
pq1ltTi6aIzi2PddgsDbIfPVWw+ip2REWgvGDrB9sqhlfs+b4grCyhSOSc8tfmLp
BY1rJMO09jJYOPRrEQ8wsmlM0pb8s3p4IBaUnEdKTNxJv23DhpoqcNDcAWDk/hN8
dGPPZjiMn/sf6FVXuIM2dU1fnF4I3d5E4a9xzCsi6jLH7WZAJ32uw1xF42lC4hfd
Y8qbD6GLl/Ngh2YgXn0g9DTbUpNa+EIixdLApremwxD1H1PYOYBoqh08PbjR+w==
-----END CERTIFICATE-----`)

var CaKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCtSttU1i5aCN3c
VfHQdxNhWV1ZiO0gAz/BQsi/OkpAgBLJU8v/zbpJ7Ajlg7OzDcKWmK2Zjlo6492Y
lCtt9cQMrI05GusjjK/+NHenmuaAztUw6rnmxyEJGq11eLNo4GC/O4q0lAmCVAhN
4YxAsXG/m+EkCO7T0TpXdYUj/9IT1XQ9MJeWXt/ifNP9YwZU9+ZEBa9KqNJV7ASm
JqX0PgUI2LS3wXfu99YKzjg86GPoAv9nSvQcixPFyYDPrxiCWGG936cbrga4VChn
1ZQAL5GLpZkdAg1IOc2ufZVFuZuSa6KHUPL3a+V65d09vBQlCrU3kQ6L8hlFjx8j
X9TAcNF5AgMBAAECggEANQAMDOpkysyjblwq1SNWHhQC5Ptn6r6TpTwTwcjGJOwG
0uR6JAZ7z1gNcITTVRQES6LulWRgXFqMz7mhfsQH7ghoOOrut7Szrv/FCNHrZcHc
mlVv/hExHWO3YZJE7PKTJGnFhm0wa1fgIlG1X6PlskCunyLMSKRZP56F1fjL+5xw
UixYiqkYhPMwVCSP4+b/rjOqrlgosVRQ0nBHmgbJkXvg/J7dPHSBodhlSKfOx7D8
EPLnTdM/QDF/2FI/wuTwHXhVsgyI3I0Ht3pc0/EQ2Rjn1nGKYk8Nv7d1M0/otVOa
1rhVVF1NLyo/PlPlf9AIJUQYkG8Bp0jbIQexVOU0sQKBgQDV7a90Afm7Wb0yCFg4
ViWV13Bdp7UXiYlLSnyMFmPGp4R0DHNGoz8ewHOeDnTaH0/mYuJmsazsMEDIuHcF
1SUZbXRqlqCig9eUEFUE4D/17zmSDJVqjJR/niFNbZ5UVCofJBuZ9L788k3va6Qz
VScrZKyOx2D4EaDEsOTmKzp84wKBgQDPX1VFaZARGEHjxD3+eLsTGot/sR4/Ms2s
jQ1WlI/3q+YGbEPmdg0I8xIElWEAKQ3IBMd4Txg6PP1aA+GUoPtAhJ/rG1LVWad9
WIxoUeAeQ/eJMU6i/YfGi1AcLhKk7A1MsbCE79dxzM19N1wiKCYOErKqUSpdqUxf
8MioOkuC8wKBgEx7K0zoH+YxEQjAHvoVIl7NpOh2urFthF4chSZ4Ire00A/FG7lX
R4uw9iS9ulz48NHG7HYWc2IFZkPcXwEA0MCkdwhcTZWMWRggNqUFnxhHrGdghFKR
a82sNO+/julLJbv4Zr3F2DoKTn6YFx4bBWPoHCD3et11P+rR6yO2tLRNAoGAWPb5
SjIjkHHrsp96STXabDOzLTD7XPmaqzBITKCnswWYRaEk8DYtGW3OiRDc8IisVOdX
/BFSv4ly169ak70MjX1YbjDmtIkmBex7MDYQBGv8QmtY5SwHl/IkiDJr5T0v53tD
04Rh1Xybm+CoMi8vRFJZPCBeIJiFH5PZQfLgemkCgYEAinsbP5u2yCxFyk8MUfxs
Up2ItRdXMI5VfsV6Vi9oAnvvdm1ayPtqJkWu5SdrQqSKmU3HYm+zRhJnkS6Q5/QO
omSo6JyabsLIxn10jlV3U/LthCKEKLmdeb5noI6fsWvly8uOJwrI+uuRvAvwuIvQ
24bk3gbMIL5iSo6OGY65yDg=
-----END PRIVATE KEY-----`)

func SetCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}