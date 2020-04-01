package jwt_token

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestKey_Read(t *testing.T) {
	can_read_from_string_reader(t)
	can_read_from_file(t)
}

func can_read_from_string_reader(t *testing.T) {
	test := func(t *testing.T, stringToRead string) {
		key := &Key{buf: nil, reader: bytes.NewReader([]byte(stringToRead))}

		buf, err := key.Read()

		if err != nil {
			t.Fatal(err)
		}
		if string(buf) != stringToRead {
			t.Fatalf("read key failed, want:%s got:%s", stringToRead, string(buf))
		}
	}

	datas := map[string]string{
		"empty string": "",
		"short buf":    "a shourt string",
		"long buf":     certPem,
	}
	for name, data := range datas {
		t.Run(name, func(t *testing.T) {
			test(t, data)
		})
	}
}

func can_read_from_file(t *testing.T) {
	t.Run("can_read_from_file", func(t *testing.T) {
		file, err := os.Open("../../res/certs/cert.pem")
		if err != nil {
			t.Fatal(err)
		}
		reader := bufio.NewReader(file)
		key := &Key{buf: nil, reader: reader}

		buf, err := key.Read()

		if err != nil {
			t.Fatal(err)
		}
		if string(buf) != certPem {
			t.Fatalf("read from file failed, want:%s got:%s", certPem, string(buf))
		}
	})
}

const certPem string = `-----BEGIN RSA PRIVATE KEY-----
MIIJKQIBAAKCAgEAzRImlV2UPPfQdk/PdP4XR7G67UFs7BZdq+FRZ2HtFQ9S+Qan
93rRn5YzuZgz0H1y/ee4KRtXkaJTTNntl5IXLQqSC0N5d7dZDjMhBkZPZustWRVD
nCmDCjyt32Ik1Ap6RvCnhVnztO7e5qNmzX+zeF+0++07LPVsHnZT4wc4ak3teVao
7FLNldtxb0t3S7nr7XS2Bgvu4tbQg8CLOHCEK+RH9YJLDLn8d+h+bffloR0QGskY
ByRzLATf6TGh47Lo49mh3Rj2WFuYlVGYDWMsE2erdeWtCa1zF0pYVV6BmhXzR3jB
k897+t5yUfn/OlPWd1KHF0LPzbZZW51ivF548vITl3JMJLL3aNoGFkULhRB2Rlr4
lHlDZ79M5HXBxcZ+bAHLMSlQysDvI8B6VpwwHRxNF5CNdD67ofAMRAHwBtgB6dh5
Bl+NdERq+Ma57DNy1ojnwtNJlpq0ToTuNwZEjAD5ak3zmyvUJnC0QmTuN4fU4f6A
ke3IbZVb3+mJtRFKcuRyn/gkdt4jk7TPfJ6VJ/s4/8SIZQBiRRWWp4vMiH/OO0bw
LNheujrIj+Vn9JvpJgpYEpgNqn1K2CWD3Evsa6/wctkgMyzwqfmnaMN0GHcqu3m/
3bu8JF9F26MTCjv0Owdny9MUSBJMgLzbWkCgGzyCiejYj9mz+W18fv9zDrcCAwEA
AQKCAgEAxq83QrAP16IiDv70eN1VoNjJyuUAqxxgVQ9Q12FBWo4mTa/tPRF516yn
IMIeMXnZ17aR7wHWrfsfye27DIc2fSUoqlENdrflSKSw8mtcstJYdV7synhNxbVU
oIFuPSKJpgGzzLeCL2LgA4V9LMz1DUNtDkiidMSzpC/wxp1QQ26NDmqv74eqN+8M
0E/FnVDdDItvcBxPrueBzqTWPjCFXiEmtu8t13665VIUbGcpzBdkaw86gHEIRTXT
2xODTiEhEuDgul5dDu5vvvv23cGgCoz+ypNkO2QZhg851jxiRO/PbQfKRXNZI44y
JMqaReDEgYcr65WH0D9EBEdQSkEd3EeXGBfd73TuTmYOwbeXAcGQYj6U9F/mFf7R
/f4Id660u8BXZ87PG5vuiN9GlSJGU3bhtJp3HAhHCfXiVAzh8dCXpKUqabdAQVIy
T6pWVbYy8YDY/0UXoPOSl1HuODLUZvG2yahWjnrurD+kpvhKmkmyU5anN3mNvvBC
9fLa2QoaXyiojw+B5szgxNo+5kUUhl0zO8UMUnbtNi7KI6+dNYC4k+vOix7RluDH
gO25D5HQl/mCI+nSEF8a68QWap6ZIW5sbfI0CRNxwJ8PWWOSPNNUyOBXsxWn802q
4H+Nsl0pvhxY5HoXFkCQ0UPEhvNWeXkXd3Ckl31ir8NPzNmZrAECggEBAPIKqB2L
Lt0TQEg/iqIGHeCT9//k/etJZYhW3kBal92cYJJEpCNpD6pyJoetCIDRn9jf/r6L
+svGSRt9N6b3NArM08Yo9HJ7K/DkV/UtiAJiGQp6SYlnwW69UqXMsxOXWbsv4VRB
6PDBKNmsfmuWBaOJH9zfJzGNZFSOYveRUgpjbSxfrV0kh00UoQaseYCryL/Sf8xS
hcyw9g8Ho1ecxTXCmcOqfmzDKQo4bu8nFLY6wC4qmGQRjGZLJiv6rO+8VGTlpsTS
E0S2FJPNrq6gUONQx68U2rZYwzn4+tXnN5Z5r0gKT32gxd+VKEEIPYBLj/MjBrqc
XxmeYrVBanmglPcCggEBANjlrrsBwuUXhpJ/ItFmQaQeaYSO8p7FYOzG4bquytxq
RpGreylbUcQq2eFdE15xhz26hlaOAtqDezm45/nc/lWSR3ccPgr2gnAnHdk7oI2I
87rUTI9cmWym6vPIjODKDnpvVtsuVJdliCjxnNR4yxS5C7sl9Md/my2i2wG0/hWQ
cpBnGXilNk34yFGNewVrZ/+H0HlzrCm6OPKf19C2LUjDVLsHiD4Dh/Sb/ciz5z1Q
kFSCQFLLqcb2nsmhQG88HyH1jnKLGo4ZigJIUvQbagS73+cMmT0UIKW1+BQ76Owh
cnfUb2NOWcDEPWdy88HKxJJLrI/iTEmKx25nM/b5pEECggEAKxEEWsViIEoFnRVZ
SH0IIeaSMQEAwTW9ECZMw2ybKv5hHIWEIxzVgcFv46JBKhKie4dXn3XuuQVeCrsc
BORlaSqK3+53mEscRW+Lyv8//RSRWhDqNr20aEzdgMzMbEb421qooEJd/UCRUTHl
CKWX+UIz3iwCoEmFOZpgN6auz1Rjn4qioTkXrfpmsHZN5Dwsqxz8SlHApuwxy/jS
8ordeDRZAby7ZATRr5TdAEaW49nOSiigFuYccjMa5qZi0QFUjuNh6hFrBkXToXzu
gPnbiqbb8OYoCFwA2LbZgufyNx8NtibHgBX0P1Ud5Xxe3Q4U3fE0iE28iiVpcNDj
7iJJfwKCAQEArZ10dvaa9dwV+S/RRAJxKpi5Z8Uwygw+YGl2CIOfmD1tjW7RIDKb
ycVvMCjbty5yzeN/Yss714OFYJf1ABl4cDCuCdbOhuH5WSLGjrte8cwdJICJY/wA
R2t5CHiQ4+J5ImH7CWkVhzZbfkKggabLECRrEnv4arRnF2mTmtMwyzwbqCEOz3aX
eGRanIT+Y0EtNuqU4pLHzYLl0LhH/SXGK8dHDIqj8NfMvO0cgMoYoAjb7vlv2ZZy
qPOB+O2dcSyT0xAG3QMh13rz8I3J8OH8xBtKx1xbUPvKgjqdgDzQXisLwIWIP9pY
l6e9axAltArFvEDcuTOwUvHGX0Y230vGwQKCAQAsMC4kybdGClQUdOxszO6pFahK
DMVGosNhWE5VBgqdVqxi+htmzJsTQ3IQ7YL8TD4sCehNF/CIT61hq/02VneKXiVU
x33XEOmEEA+akp2YsHjQ77ah7/yLw7GsZLseNyluBbYm/CmScQNgo2W46cVny/7E
5xJqif/+g/oQfYscl2yEmZY+B9vc+7X8xbwebqU2b4X5NkRVRLAICR0o7y55+PRC
OCfpgL8K/mJflyhx650jlCnDFP7ttk54nFhcyt4tPu6fsi5Jv0fhGu5L0Zy7z9S1
EbPrRXXR8a3HIVifFOpgdPwfYAcxGlJo8TTBhhcvfT+TBnXqLUlhGIuL1AhP
-----END RSA PRIVATE KEY-----
`
