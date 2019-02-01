package internal

import (
	"bytes"
	"testing"
)

func TestRouteRegistry(t *testing.T) {
	type result struct {
		hostname string
		port     int
		proxy    bool
	}

	cases := []struct {
		Config string
		Tests  map[string]result
	}{
		{
			Config: `
# local
::1                  1.2.3.4:80
0:0:0:0:0:0:0:1      1.2.3.4:80
localhost            1.2.3.4
/127\.0\.0\.\d+/     1.2.3.4:80

# home
home                 1.2.3.4:80
*.home               2.3.4.5:443

# my
${myid}  1.2.3.4:80
*.${myid} 2.3.4.5:443

# peer
92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8 1.2.3.4:80
*.92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8 2.3.4.5:443
/.*\.[a-zA-Z0-9]{25,}/ 3.4.5.6:28080

# web
go.universe.tf 1.2.3.4
*.universe.tf 2.3.4.5
#
google.* 3.4.5.6
/gooo+gle\.com/ 4.5.6.7
foobar.net 6.7.8.9 PROXY
*  9.2.3.4:8080
#/.*/ 9.2.3.4:8080
`,
			Tests: map[string]result{
				"localhost": result{"1.2.3.4", 0, false},
				"127.0.0.1": result{"1.2.3.4", 80, false},

				"home":        result{"1.2.3.4", 80, false},
				"fe.bbb.home": result{"2.3.4.5", 443, false},

				"921sm3fxr9v5wwh08d7nvnks5a37px0tdj8qd8e0cc60acy514r61r":   result{"1.2.3.4", 80, false},
				"*.921sm3fxr9v5wwh08d7nvnks5a37px0tdj8qd8e0cc60acy514r61r": result{"2.3.4.5", 443, false},

				"92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8":     result{"1.2.3.4", 80, false},
				"git.92114bmb5wjn6hfz0qb2jdr1qc2a5j3hqcr7efsfe2gj09yjmj5eg8": result{"2.3.4.5", 443, false},
				"git.920j8197p3mq6fdb1xhpserz22x146pzmnr2hj987ay8nea39fr40h": result{"3.4.5.6", 28080, false},

				"go.universe.tf":      result{"1.2.3.4", 0, false},
				"foo.universe.tf":     result{"2.3.4.5", 0, false},
				"bar.universe.tf":     result{"2.3.4.5", 0, false},
				"foo.bar.universe.tf": result{"2.3.4.5", 0, false},

				"google.com":         result{"3.4.5.6", 0, false},
				"google.fr":          result{"3.4.5.6", 0, false},
				"google.com.br":      result{"3.4.5.6", 0, false},
				"goooooooooogle.com": result{"4.5.6.7", 0, false},
				"foobar.net":         result{"6.7.8.9", 0, true},

				"blah.com":      result{"9.2.3.4", 8080, false},
				"googlecom.br":  result{"9.2.3.4", 8080, false},
				"goooooglexcom": result{"9.2.3.4", 8080, false},
			},
		},
	}

	for _, test := range cases {
		cfg := NewRouteRegistry("921sm3fxr9v5wwh08d7nvnks5a37px0tdj8qd8e0cc60acy514r61r")

		if err := cfg.Read(bytes.NewBufferString(test.Config)); err != nil {
			t.Fatalf("Failed to read config (%s):\n%q", err, test.Config)
		}

		for hostname, expected := range test.Tests {
			be, proxy := cfg.Match(hostname)
			if expected.hostname != be[0].Hostname || expected.port != be[0].Port {
				t.Errorf("cfg.Match(%q) is %v and %v, want %v", hostname, *be[0], proxy, expected)
			}
		}
	}
}
